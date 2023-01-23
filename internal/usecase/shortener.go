package usecase

import (
	"context"

	"github.com/bbt-t/shortenerURL/internal/entity"

	"github.com/gofrs/uuid"
)

type DatabaseRepository interface {
	NewUser(userID uuid.UUID)
	GetOriginalURL(shortURL string) (string, error)
	GetURLArrayByUser(userID uuid.UUID, baseURL string) ([]map[string]string, error)
	SaveShortURL(userID uuid.UUID, shortURL, originalURL string) error
	PingDB() error
	DelURLArray(ctx context.Context, userID uuid.UUID, inpJSON []byte) error
	SaveURLArray(ctx context.Context, uid uuid.UUID, inpURL []entity.URLBatchInp) error
}

type ShortenerService struct {
	repo DatabaseRepository
}

func NewShortener(r DatabaseRepository) *ShortenerService {
	return &ShortenerService{
		repo: r,
	}
}

func (s ShortenerService) NewUser(userID uuid.UUID) {
	s.repo.NewUser(userID)
}

func (s ShortenerService) GetOriginalURL(shortURL string) (string, error) {
	result, err := s.repo.GetOriginalURL(shortURL)
	return result, err
}

func (s ShortenerService) GetURLArrayByUser(userID uuid.UUID, baseURL string) ([]map[string]string, error) {
	result, err := s.repo.GetURLArrayByUser(userID, baseURL)
	return result, err
}

func (s ShortenerService) SaveShortURL(userID uuid.UUID, shortURL, originalURL string) error {
	return s.repo.SaveShortURL(userID, shortURL, originalURL)
}

func (s ShortenerService) PingDB() error {
	return s.repo.PingDB()
}

func (s ShortenerService) DelURLArray(ctx context.Context, userID uuid.UUID, inpJSON []byte) error {
	return s.repo.DelURLArray(ctx, userID, inpJSON)
}

func (s ShortenerService) SaveURLArray(ctx context.Context, uid uuid.UUID, inpURL []entity.URLBatchInp) error {
	return s.repo.SaveURLArray(ctx, uid, inpURL)
}
