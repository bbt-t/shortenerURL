package storage

import (
	"context"

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
	DelURLArray(ctx context.Context, userID uuid.UUID, inpURLs []string) error
	SaveURLArray(ctx context.Context, uid uuid.UUID, inpURL []entity.URLBatchInp) error
}
