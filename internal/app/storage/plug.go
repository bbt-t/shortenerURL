package storage

import (
	"errors"
	"log"
	"sync"
)

type mapDBPlug struct {
	/*
		Simple DB stub.
	*/
	mapURL map[string]string
	mutex  sync.Mutex
}

func NewMapDBPlug() DBRepo {
	/*
		return: object with an empty map to write data.
	*/
	return &mapDBPlug{mapURL: map[string]string{}}
}

func (m *mapDBPlug) GetURL(k string) (string, error) {
	/*
		Get info from the map by key.
	*/
	defer m.mutex.Unlock()

	m.mutex.Lock()
	result, ok := m.mapURL[k]

	if !ok {
		return "", errors.New("no such id in DB")
	}
	return result, nil
}

func (m *mapDBPlug) SaveURL(k, v string) error {
	/*
		Write info to the map by key - value.
	*/
	m.mutex.Lock()
	_, ok := m.mapURL[k]
	m.mutex.Unlock()
	if ok {
		return nil
	}

	m.mutex.Lock()
	m.mapURL[k] = v
	m.mutex.Unlock()
	return nil
}

func (m *mapDBPlug) Ping() error {
	log.Println("MAP IS READY!")
	return nil
}
