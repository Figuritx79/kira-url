package url

import (
	"kira-url/internal/database/models"
	dtourl "kira-url/internal/dto/url"
)

type URLRepository interface {
	FindByShortURL(code string) (*dtourl.URLResponse, error)
	Save(url *models.URL) error
	FindByURL(url string) (*dtourl.ShortURLResponse, error)
	Update(url models.URL, code string) error
	Updates(urls []models.URL) error
}
