package usecase

import "github.com/gofrs/uuid"

type DatabaseRepository interface {
	NewUser(userID uuid.UUID)
	GetOriginalURL(shortURL string) (string, error)
	SaveShortURL(userID uuid.UUID, hashURL, originalURL string) error
	PingDB() error
}

type ShortenerService struct {
	repo DatabaseRepository
}

func NewShortener(r DatabaseRepository) *ShortenerService {
	return &ShortenerService{
		repo: r,
	}
}

func (s ShortenerService) GetOriginalURL(shortURL string) (string, error) {
	result, err := s.repo.GetOriginalURL(shortURL)
	return result, err
}

func (s ShortenerService) SaveShortURL(userID uuid.UUID, hashURL, originalURL string) error {
	return s.repo.SaveShortURL(userID, hashURL, originalURL)
}

func (s ShortenerService) PingDB() error {
	return s.repo.PingDB()
}

func (s ShortenerService) NewUser(userID uuid.UUID) {
	s.repo.NewUser(userID)
}
