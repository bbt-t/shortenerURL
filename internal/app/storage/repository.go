package storage

type DBRepo interface {
	GetURL(shortURL string) (string, error)
	SaveURL(originalURL string, id string) error
}
