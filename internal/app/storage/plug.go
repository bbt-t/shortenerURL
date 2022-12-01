package storage

import (
	"errors"
	"sync"
)

type MapDBPlug struct {
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
	return &MapDBPlug{mapURL: map[string]string{}}
}

func (m *MapDBPlug) GetURL(k string) (string, error) {
	/*
		Get info from the map by key.
	*/
	defer m.mutex.Unlock()

	m.mutex.Lock()
	result, IsOk := m.mapURL[k]

	if IsOk == false {
		return "", errors.New("no such id in DB")
	}
	return result, nil
}

func (m *MapDBPlug) SaveURL(k, v string) error {
	/*
		Write info to the map by key - value.
	*/
	m.mutex.Lock()
	_, IsOk := m.mapURL[k]
	m.mutex.Unlock()
	if IsOk == true {
		return nil
	}

	m.mutex.Lock()
	m.mapURL[k] = v
	m.mutex.Unlock()
	return nil
}
