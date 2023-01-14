package storage

import "github.com/gofrs/uuid"

type DatabaseRepository interface {
	/*
		Interface for using DB.
	*/
	NewUser(userID uuid.UUID)
	GetOriginalURL(shortURL string) (string, error)
	GetURLArrayByUser(userID uuid.UUID, baseURL string) ([]map[string]string, error)
	SaveShortURL(userID uuid.UUID, shortURL, originalURL string) error
	PingDB() error
	DelURLArray(inpJSON []byte, userID string) error
}
