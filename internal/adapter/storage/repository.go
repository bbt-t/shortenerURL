package storage

import (
	"github.com/bbt-t/shortenerURL/internal/entity"
	"github.com/gofrs/uuid"
)

type DatabaseRepository interface {
	/*
		Interface for using DB.
	*/
	NewUser(userID uuid.UUID)
	GetOriginalURL(shortURL string) (string, error)
	GetURLArrayByUser(userID uuid.UUID, baseURL string) ([]map[string]string, error)
	SaveShortURL(userID uuid.UUID, shortURL, originalURL string) error
	PingDB() error
	DelURLArray(userID uuid.UUID, inpJSON []byte) error
	SaveURLArray(uid uuid.UUID, inpURL []entity.URLBatchInp) error
}
