package storage

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sync"
)

type fileDB struct {
	/*
		Storage implementation in file.
	*/
	PathToFile string
	mutex      *sync.Mutex
}

func NewFileDB(pathFile string) DBRepo {
	/*
		Initialize storage in file.
		param pathFile: path to file (with full name)
	*/
	return &fileDB{
		PathToFile: pathFile,
		mutex:      &sync.Mutex{},
	}
}

func (f *fileDB) save(k, v string) error {
	/*
		Create/overwrite and write to a file.gob (gob-format).
	*/
	defer f.mutex.Unlock()
	f.mutex.Lock()

	DATA, err := f.get()
	if err != nil {
		DATA = make(map[string]string)
	}
	DATA[k] = v

	saveTo, errOpen := os.OpenFile(f.PathToFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0700)

	if errOpen != nil {
		fmt.Println("Cannot create ->", f.PathToFile)
		return errOpen
	}

	defer saveTo.Close()

	encoder := gob.NewEncoder(saveTo)
	err = encoder.Encode(DATA)
	if err != nil {
		log.Printf("Cannot save to -> %v", f.PathToFile)
		return err
	}
	return nil
}

func (f *fileDB) get() (map[string]string, error) {
	/*
		Open file and take map-object from there.
	*/
	var DATA map[string]string

	loadFrom, err := os.OpenFile(f.PathToFile, os.O_RDONLY, 0700)

	if err != nil {
		if os.IsPermission(err) {
			log.Println("Error: Read permission denied.")
		}
		log.Printf("Empty key/value store! ERROR : %v", err)
		return nil, err
	}
	decoder := gob.NewDecoder(loadFrom)
	if err := decoder.Decode(&DATA); err != nil {
		return DATA, err
	}
	return DATA, nil
}

func (f *fileDB) GetURL(k string) (string, error) {
	/*
		get value by key from file.
	*/
	if err := f.Ping(); err != nil {
		log.Fatal(err)
	}
	defer f.mutex.Unlock()
	fileMap, _ := f.get()
	f.mutex.Lock()
	originalURL, ok := fileMap[k]
	if !ok {
		return "", errDBUnknownID
	}
	return originalURL, nil
}

func (f *fileDB) GetAllURL() ([]map[string]string, error) {
	/*
		Take all saved urls.
	*/
	if err := f.Ping(); err != nil {
		log.Fatal(err)
	}
	defer f.mutex.Unlock()
	f.mutex.Lock()

	var allURL []map[string]string

	allMap, _ := f.get()
	if len(allMap) == 0 {
		return nil, errDBEmpty
	}
	for k, v := range allMap {
		temp := make(map[string]string)
		temp[k] = v
		allURL = append(allURL, temp)
	}

	return allURL, nil
}

func (f *fileDB) SaveURL(k, v string) error {
	/*
		Calling a func to save info to a file.
	*/
	err := f.save(k, v)
	return err
}

func (f *fileDB) Ping() error {
	/*
		return error if file or filename does not exist
	*/
	if f.PathToFile == "" {
		log.Println("--- missing filename ---")
		return errDBFileDoesNotExist
	}
	if _, err := os.Stat(f.PathToFile); err != nil {
		if os.IsNotExist(err) {
			log.Println("--- file does not exist ---")
			return errDBFileDoesNotExist
		}
	}
	return nil
}
