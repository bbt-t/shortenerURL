package storage

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/bbt-t/shortenerURL/internal/entity"

	"github.com/gofrs/uuid"
)

type fileDB struct {
	/*
		Storage implementation in file.
	*/
	PathToFile string
	mutex      *sync.RWMutex
}

func newFileDB(pathFile string) *fileDB {
	/*
		Initialize storage in file.
		param pathFile: path to file (with full name)
	*/
	return &fileDB{
		PathToFile: pathFile,
		mutex:      new(sync.RWMutex),
	}
}

func (f *fileDB) saveToFile(data map[uuid.UUID][]entity.DBMapFilling) error {
	/*
		Create/overwrite and write to a file.gob (gob-format).
	*/
	saveTo, errOpen := os.OpenFile(f.PathToFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0700,
	)

	if errOpen != nil {
		fmt.Println("Cannot read file ->", f.PathToFile)
		return errOpen
	}

	defer saveTo.Close()

	encoder := gob.NewEncoder(saveTo)
	err := encoder.Encode(data)
	if err != nil {
		log.Printf("Cannot save to -> %+v", f.PathToFile)
		return err
	}
	return nil
}

func (f *fileDB) get() (map[uuid.UUID][]entity.DBMapFilling, error) {
	/*
		Open file and take map-object from there.
	*/
	data := make(map[uuid.UUID][]entity.DBMapFilling)

	loadFrom, err := os.OpenFile(f.PathToFile,
		os.O_RDONLY,
		0700,
	)

	if err != nil {
		if os.IsPermission(err) {
			log.Println("Error: Read permission denied.")
		}
		return nil, err
	}
	decoder := gob.NewDecoder(loadFrom)
	if err := decoder.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func (f *fileDB) NewUser(uid uuid.UUID) {
	/*
		Create new user in DB.
	*/

	data, _ := f.get()
	data[uid] = []entity.DBMapFilling{}
	if errSave := f.saveToFile(data); errSave != nil {
		log.Println(errSave)
	}
}

func (f *fileDB) GetOriginalURL(k string) (string, error) {
	/*
		get value by key from file.
	*/
	var result string

	defer f.mutex.RUnlock()

	if err := f.PingDB(); err != nil {
		log.Fatal(err)
	}
	fileMap, _ := f.get()
	f.mutex.RLock()
	for _, v := range fileMap {
		for _, val := range v {
			if k == val.ShortURL && val.Deleted {
				return "", errDBUnknownID
			}
			if k == val.ShortURL {
				result = val.OriginalURL
			}
			if result != "" {
				break
			}
		}
	}
	if result == "" {
		return "", errDBUnknownID
	}
	return result, nil
}

func (f *fileDB) GetURLArrayByUser(uid uuid.UUID, baseURL string) ([]map[string]string, error) {
	defer f.mutex.RUnlock()
	fileMap, _ := f.get()
	f.mutex.RLock()

	allURL, ok := fileMap[uid]
	if !ok || len(allURL) == 0 {
		return nil, errDBEmpty
	}

	convInfo := make(map[string]string)
	for _, item := range allURL {
		convInfo[item.ShortURL] = item.OriginalURL
	}
	result := convertToArrayMap(convInfo, baseURL)

	return result, nil
}

func (f *fileDB) SaveShortURL(uid uuid.UUID, k, v string) error {
	/*
		Calling a func to save info to a file.
	*/

	f.mutex.Lock()

	data, err := f.get()
	if err != nil {
		data = make(map[uuid.UUID][]entity.DBMapFilling) //map[uuid.UUID][]entity.DBMapFilling
	}

	for _, v := range data[uid] {
		if v.ShortURL == k {
			return errHTTPConflict
		}
	}
	f.mutex.Unlock()

	f.mutex.Lock()
	data[uid] = append(data[uid], entity.DBMapFilling{
		OriginalURL: v,
		ShortURL:    k,
		Deleted:     false,
	})
	f.mutex.Unlock()

	if errSave := f.saveToFile(data); errSave != nil {
		log.Println(errSave)
		return errSave
	}

	return nil
}

func (f *fileDB) PingDB() error {
	/*
		return error if file or filename does not exist
	*/
	_, err := os.Stat(f.PathToFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("--- file does not exist ---\n:: create new file ::")
			_, err = os.OpenFile(f.PathToFile,
				os.O_RDWR|os.O_TRUNC|os.O_CREATE,
				0700,
			)
		}
	}
	if err == nil {
		log.Println("FILE IS READY!")
	}
	return err
}

func (f *fileDB) DelURLArray(ctx context.Context, uid uuid.UUID, inpURLs []string) error {
	fileMap, _ := f.get()
	for i, item := range fileMap[uid] {
		for _, v := range inpURLs {
			if item.ShortURL == v {
				fileMap[uid][i].Deleted = true
			}
		}
	}
	if errSave := f.saveToFile(fileMap); errSave != nil {
		log.Println(errSave)
		return errSave
	}
	ctx.Done()

	return nil
}

func (f *fileDB) SaveURLArray(ctx context.Context, uid uuid.UUID, urlBatch []entity.URLBatchInp) error {
	fileMap, _ := f.get()

	for i, item := range urlBatch {
		temp := strings.Split(item.ShortURL, "/")
		urlBatch[i].ShortURL = temp[len(temp)-1]
	}

	for _, v := range fileMap[uid] {
		for _, item := range urlBatch {
			if v.OriginalURL != item.OriginalURL {
				fileMap[uid] = append(fileMap[uid], entity.DBMapFilling{
					OriginalURL: item.OriginalURL,
					ShortURL:    item.ShortURL,
					Deleted:     false,
				})
			}
		}
	}
	if errSave := f.saveToFile(fileMap); errSave != nil {
		log.Println(errSave)
		return errSave
	}
	ctx.Done()

	return nil
}
