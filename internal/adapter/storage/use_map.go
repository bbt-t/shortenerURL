package storage

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/bbt-t/shortenerURL/internal/entity"
	"github.com/bbt-t/shortenerURL/pkg"

	"github.com/gofrs/uuid"
)

type mapDB struct {
	/*
		Simple DB stub.
	*/
	mutex  *sync.RWMutex
	mapURL map[uuid.UUID][]entity.DBMapFilling
}

func NewMapDB() DatabaseRepository {
	/*
		return: object with an empty map to write data.
	*/
	return &mapDB{
		mapURL: make(map[uuid.UUID][]entity.DBMapFilling),
		mutex:  new(sync.RWMutex),
	}
}

func (m *mapDB) NewUser(uid uuid.UUID) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	if _, ok := m.mapURL[uid]; !ok {
		m.mapURL[uid] = []entity.DBMapFilling{}
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
		for _, val := range v {
			if k == val.ShortURL && val.Deleted {
				return "", ErrDBUnknownID
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
		return "", ErrDBUnknownID
	}
	return result, nil
}

func (m *mapDB) GetURLArrayByUser(uid uuid.UUID, baseURL string) ([]map[string]string, error) {
	/*
		Take all saved urls.
	*/

	defer m.mutex.RUnlock()
	m.mutex.RLock()

	allURL, ok := m.mapURL[uid]
	if !ok || len(allURL) == 0 {
		return nil, ErrDBEmpty
	}

	convInfo := make(map[string]string)
	for _, item := range allURL {
		convInfo[item.ShortURL] = item.OriginalURL
	}
	result := pkg.ConvertToArrayMap(convInfo, baseURL)

	return result, nil
}

func (m *mapDB) SaveShortURL(uid uuid.UUID, k, v string) error {
	/*
		Write info to the map by key - value.
	*/
	m.mutex.RLock()
	for _, v := range m.mapURL[uid] {
		if v.ShortURL == k {
			return ErrHTTPConflict
		}
	}
	m.mutex.RUnlock()

	m.mutex.Lock()
	m.mapURL[uid] = append(m.mapURL[uid], entity.DBMapFilling{
		OriginalURL: v,
		ShortURL:    k,
		Deleted:     false,
	})
	m.mutex.Unlock()
	return nil
}

func (m *mapDB) PingDB() error {
	log.Println("MAP IS READY!")
	return nil
}

func (m *mapDB) DelURLArray(ctx context.Context, uid uuid.UUID, inpURLs []string) error {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	m.mapURL = deleteURLArrayMap(m.mapURL, uid, inpURLs)
	ctx.Done()
	return nil
}

func (m *mapDB) SaveURLArray(ctx context.Context, uid uuid.UUID, urlBatch []entity.URLBatchInp) error {
	for i, item := range urlBatch {
		temp := strings.Split(item.ShortURL, "/")
		urlBatch[i].ShortURL = temp[len(temp)-1]
	}
	defer m.mutex.Unlock()
	m.mutex.Lock()
	m.mapURL = saveURLBatchMap(m.mapURL, uid, urlBatch)
	ctx.Done()
	return nil
}
