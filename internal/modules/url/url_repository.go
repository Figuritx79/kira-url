package url

import (
	"kira-url/internal/database/models"

	"github.com/google/uuid"
)

type URLRepository interface {
	FindByShortURL(code string) (*URLResponse, error)
	Save(url models.URL) error
	FindByURL(url string) (*ShortURLResponse, error)
	Update(url models.URL, code string) error
}
