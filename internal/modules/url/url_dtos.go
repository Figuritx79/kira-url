package url

import "github.com/google/uuid"

// This struct is used when the user wants the original url by code
type URLResponse struct {
	ID          uuid.UUID `json:"-"`
	VisitCount  int64     `json:"-"`
	OriginalURL string    `json:"original_url"`
}

// This struct is used when we find a duplicate record in database and we send the found record
type ShortURLResponse struct {
	ShortURL         string `json:"-"`
	CompleteShortURL string `json:"short_url"`
}

type URLCompleteResponse struct {
	ID          uuid.UUID `json:"id,omitempty"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

type CreatURL struct {
	OriginalURL string `json:"original_url"`
	CustomCode  string `json:"custom_code,omitempty"`
	ShortURL    string `json:"-"`
}
