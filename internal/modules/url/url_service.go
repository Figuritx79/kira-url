package url

import (
	"errors"
	"log/slog"

	"kira-url/internal/base62"
	"kira-url/internal/constants"
	"kira-url/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type urlService struct {
	repository URLRepository
	log        *slog.Logger
}

func newURLService(repository URLRepository, log *slog.Logger) *urlService {
	return &urlService{repository: repository, log: log}
}

func (service *urlService) FindByShortURL(code string) (*URLResponse, error) {
	shortURL, err := service.repository.FindByShortURL(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			service.log.Warn("URL not found")
			return nil, ErrURLNotFound
		}
		service.log.Error("Error", "URL", err)
		return nil, err
	}

	if shortURL == nil {
		service.log.Error("Error", "URL", err)
		return nil, err
	}

	return shortURL, nil
}

func (service *urlService) FindByURL(url string) (*ShortURLResponse, bool, error) {
	found, err := service.repository.FindByURL(url)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return found, true, nil
}

func (service *urlService) Save(url *CreatURL) (*URLCompleteResponse, error) {
	randNumber := base62.RandamBase62Number()

	code := base62.EncodeToBase62(randNumber)

	newURL := models.URL{
		ShortURL:    code,
		OriginalURL: url.OriginalURL,
	}

	err := service.repository.Save(newURL)
	if err != nil {
		return nil, err
	}

	return &URLCompleteResponse{
		ShortURL:    constants.BaseDomain + newURL.ShortURL,
		OriginalURL: newURL.OriginalURL,
	}, nil
}

func (service *urlService) Update(id uuid.UUID, visitCount int64) error {
	url := models.URL{
		VisitCount: visitCount + 1,
		ID:         id,
	}

	err := service.repository.Update(url, id)
	if err != nil {
		return err
	}
	return nil
}
