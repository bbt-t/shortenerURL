package storage

type DBRepo interface {
	/*
		Interface for using DB. Save and get values.
	*/
	GetURL(shortURL string) (string, error)
	SaveURL(originalURL string, id string) error
	Ping() error
}