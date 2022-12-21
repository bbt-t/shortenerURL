package storage

import (
	"log"
	"sync"
)

type mapDBPlug struct {
	/*
		Simple DB stub.
	*/
	mapURL  map[string]string
	mapUser map[string]interface{}
	mutex   sync.Mutex
}

func NewMapDBPlug() DBRepo {
	/*
		return: object with an empty map to write data.
	*/
	return &mapDBPlug{
		mapURL: map[string]string{},
		mutex:  sync.Mutex{},
	}
}

func (m *mapDBPlug) GetURL(k string) (string, error) {
	/*
		get info from the map by key.
	*/
	defer m.mutex.Unlock()
	m.mutex.Lock()

	result, ok := m.mapURL[k]
	if !ok {
		return "", errDBUnknownID
	}

	return result, nil
}

func (m *mapDBPlug) GetAllURL() ([]map[string]string, error) {
	/*
		Take all saved urls.
	*/
	defer m.mutex.Unlock()
	m.mutex.Lock()

	var allUrl []map[string]string
	if len(m.mapURL) == 0 {
		return nil, errDBEmpty
	}
	for k, v := range m.mapURL {
		temp := make(map[string]string)
		temp[k] = v
		allUrl = append(allUrl, temp)
	}

	return allUrl, nil
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

func (m *mapDBPlug) SaveUser(k string, v interface{}) error {
	m.mutex.Lock()
	_, ok := m.mapUser[k]
	m.mutex.Unlock()
	if ok {
		return nil
	}

	m.mutex.Lock()
	m.mapUser[k] = v
	m.mutex.Unlock()
	return nil
}

func (m *mapDBPlug) Ping() error {
	log.Println("MAP IS READY!")
	return nil
}
