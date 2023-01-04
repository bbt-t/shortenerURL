package entity

type Shortener struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
