package url

import "github.com/google/uuid"

type URLResponse struct {
	OriginalURL string
}
type URLCompleteResponse struct {
	ID          uuid.UUID
	ShortURL    string
	OriginalURL string
}

type CreatURL struct {
	OriginalURL string `json:"original_url"`
}
