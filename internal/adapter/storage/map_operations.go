package storage

import (
	"github.com/bbt-t/shortenerURL/internal/entity"

	"github.com/gofrs/uuid"
)

func deleteURLArrayMap(db map[uuid.UUID][]entity.DBMapFilling, uid uuid.UUID, inpURLs []string) map[uuid.UUID][]entity.DBMapFilling {
	for i, item := range db[uid] {
		for _, v := range inpURLs {
			if item.ShortURL == v {
				db[uid][i].Deleted = true
			}
		}
	}
	return db
}

func saveURLBatchMap(db map[uuid.UUID][]entity.DBMapFilling, uid uuid.UUID, urlBatch []entity.URLBatchInp) map[uuid.UUID][]entity.DBMapFilling {
	for _, v := range db[uid] {
		for _, item := range urlBatch {
			if v.OriginalURL != item.OriginalURL {
				db[uid] = append(db[uid], entity.DBMapFilling{
					OriginalURL: item.OriginalURL,
					ShortURL:    item.ShortURL,
					Deleted:     false,
				})
			}
		}
	}
	return db
}
