package entity

import "github.com/gofrs/uuid"

type DBMapFilling struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	Deleted     bool   `json:"-"`
}

type URLBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"-"`
	ShortURL      string `json:"short_url"`
}

type URLBatchInp struct {
	UserID        uuid.UUID `db:"user_id" json:"-"`
	ID            uuid.UUID `db:"id" json:"-"`
	CorrelationID string    `json:"correlation_id"`
	OriginalURL   string    `db:"original_url" json:"original_url"`
	ShortURL      string    `db:"short_url" json:"short_url"`
}

type URLs struct {
	OriginalURL string `db:"original_url" json:"original_url"`
	ShortURL    string `db:"short_url" json:"short_url"`
}

type CheckURL struct {
	Deleted     bool   `db:"deleted"`
	OriginalURL string `db:"original_url"`
}
