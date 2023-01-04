package storage

import (
	"encoding/gob"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"os"
	"sync"
)

type fileDB struct {
	/*
		Storage implementation in file.
	*/
	PathToFile string
	mutex      *sync.RWMutex
}

func NewFileDB(pathFile string) DatabaseRepository {
	/*
		Initialize storage in file.
		param pathFile: path to file (with full name)
	*/
	return &fileDB{
		PathToFile: pathFile,
		mutex:      new(sync.RWMutex),
	}
}

// TODO: НЕ РАЮБОТАЕТ ! сделать на подобие с мапой!
func (f *fileDB) save(userID uuid.UUID, k, v string, empty bool) error {
	/*
		Create/overwrite and write to a file.gob (gob-format).
	*/
	defer f.mutex.Unlock()
	f.mutex.Lock()

	data, err := f.get()
	if err != nil {
		data = make(map[uuid.UUID]map[string]string)
	}
	if empty {
		data[userID] = make(map[string]string)
	}

	_, ok := data[userID][k]
	if ok {
		return errHTTPConflict
	}
	data[userID][k] = v

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
	err = encoder.Encode(data)
	if err != nil {
		log.Printf("Cannot save to -> %v", f.PathToFile)
		return err
	}
	return nil
}

func (f *fileDB) get() (map[uuid.UUID]map[string]string, error) {
	/*
		Open file and take map-object from there.
	*/
	var data map[uuid.UUID]map[string]string

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
		result = v[k]
		if result != "" {
			break
		}
	}
	if result == "" {
		return "", errDBUnknownID
	}
	return result, nil
}

func (f *fileDB) SaveShortURL(userID uuid.UUID, k, v string) error {
	/*
		Calling a func to save info to a file.
	*/
	err := f.save(userID, k, v, false)
	return err
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

func (f *fileDB) NewUser(userID uuid.UUID) {
	_ = f.save(userID, "", "", true)
}

func (f *fileDB) GetURLArrayByUser(userID uuid.UUID) (map[string]string, error) {
	defer f.mutex.RUnlock()

	fileMap, _ := f.get()

	f.mutex.RLock()
	urlArray, ok := fileMap[userID]
	if !ok {
		return nil, errDBEmpty
	}
	return urlArray, nil
}
