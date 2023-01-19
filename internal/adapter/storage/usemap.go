package storage

import (
	"log"
	"sync"

	"github.com/bbt-t/shortenerURL/internal/entity"

	"github.com/gofrs/uuid"
)

type mapDB struct {
	/*
		Simple DB stub.
	*/
	mapURL map[uuid.UUID]map[string]string
	mutex  *sync.RWMutex
}

func NewMapDB() DatabaseRepository {
	/*
		return: object with an empty map to write data.
	*/
	return &mapDB{
		mapURL: make(map[uuid.UUID]map[string]string),
		mutex:  new(sync.RWMutex),
	}
}

func (m *mapDB) NewUser(userID uuid.UUID) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	if nil == m.mapURL[userID] {
		m.mapURL[userID] = make(map[string]string)
	}
}

func (m *mapDB) GetOriginalURL(k string) (string, error) {
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

func (m *mapDB) GetURLArrayByUser(userID uuid.UUID, baseURL string) ([]map[string]string, error) {
	/*
		Take all saved urls.
	*/

	defer m.mutex.RUnlock()
	m.mutex.RLock()

	allURL, ok := m.mapURL[userID]
	if !ok || len(allURL) == 0 {
		return nil, errDBEmpty
	}
	result := convertToArrayMap(allURL, baseURL)

	return result, nil
}

func (m *mapDB) SaveShortURL(userID uuid.UUID, k, v string) error {
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

func (m *mapDB) PingDB() error {
	log.Println("MAP IS READY!")
	return nil
}

func (m *mapDB) DelURLArray(_ uuid.UUID, _ []byte) error {
	return nil
}

func (m *mapDB) SaveURLArray(_ uuid.UUID, _ []entity.URLBatchInp) error {
	return nil
}
