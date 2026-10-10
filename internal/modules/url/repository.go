package url

import (
	"context"
	"kira-url/internal/database/models"
	dtourl "kira-url/internal/dto/url"
)

type URLRepository interface {
	FindByShortURL(ctx context.Context, code string) (*dtourl.URLResponse, error)
	Save(ctx context.Context, url *models.URL) error
	FindByURL(ctx context.Context, url string) (*dtourl.ShortURLResponse, error)
	Update(ctx context.Context, url models.URL, code string) error
	Updates(urls []models.URL) error
}
