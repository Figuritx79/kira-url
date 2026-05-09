package url

import "kira-url/internal/database/models"

type URLRepository interface {
	FindByShortURL(code string) (*URLResponse, error)
	Save(url models.URL) error
	FindByURL(url string) (*ShortURLResponse, error)
	// Update(url Creat, id uuid.UUID) error
}
