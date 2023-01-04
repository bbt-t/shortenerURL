package storage

import (
	"log"
	"sync"

	"github.com/gofrs/uuid"
)

type mapDBPlug struct {
	/*
		Simple DB stub.
	*/
	mapURL map[uuid.UUID]map[string]string
	mutex  *sync.RWMutex
}

func NewMapDBPlug() DatabaseRepository {
	/*
		return: object with an empty map to write data.
	*/
	return &mapDBPlug{
		mapURL: make(map[uuid.UUID]map[string]string),
		mutex:  new(sync.RWMutex),
	}
}

func (m *mapDBPlug) GetOriginalURL(k string) (string, error) {
	/*
		get info from the map by key.
	*/
	var result string
	defer m.mutex.RUnlock()
	m.mutex.RLock()

	for _, v := range m.mapURL {
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

func (m *mapDBPlug) GetURLArrayByUser(userID uuid.UUID) (map[string]string, error) {
	/*
		Take all saved urls.
	*/
	defer m.mutex.RUnlock()
	m.mutex.RLock()

	allURL, ok := m.mapURL[userID]
	if !ok {
		return nil, errDBEmpty
	}

	return allURL, nil
}

func (m *mapDBPlug) SaveShortURL(userID uuid.UUID, k, v string) error {
	/*
		Write info to the map by key - value.
	*/
	m.mutex.RLock()
	_, ok := m.mapURL[userID][k]
	m.mutex.RUnlock()
	if ok {
		return errHTTPConflict
	}

	m.mutex.Lock()
	m.mapURL[userID][k] = v
	m.mutex.Unlock()
	return nil
}

func (m *mapDBPlug) PingDB() error {
	log.Println("MAP IS READY!")
	return nil
}

func (m *mapDBPlug) NewUser(userID uuid.UUID) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	if m.mapURL[userID] == nil {
		m.mapURL[userID] = make(map[string]string)
	}
}
