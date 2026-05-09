package url

import (
	"context"

	"kira-url/internal/database"
	"kira-url/internal/database/models"

	"gorm.io/gorm"
)

var _ URLRepository = (*urlGormRepository)(nil)

type urlGormRepository struct {
	db *gorm.DB
}

func NewURLGormRepository(db *gorm.DB) *urlGormRepository {
	return &urlGormRepository{db: db}
}

func (repository *urlGormRepository) FindByShortURL(code string) (*URLResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), database.DEFAULT_TIMEOUT)

	defer cancel()

	var url *URLResponse

	err := repository.db.WithContext(ctx).
		Model(&models.URL{}).
		Select("urls.original_url").
		Where("urls.short_url= ?", code).
		First(&url).
		Error
	if err != nil {
		return nil, err
	}

	if url == nil {
		return nil, err
	}
	return url, nil
}

func (repository *urlGormRepository) Save(url models.URL) error {
	ctx, cancel := context.WithTimeout(context.Background(), database.DEFAULT_TIMEOUT)
	defer cancel()
	err := repository.db.WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			url.BeforeCreate(tx)
			if err := tx.Model(&models.URL{}).Create(&url).Error; err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

func (repository *urlGormRepository) FindByURL(url string) (*ShortURLResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), database.DEFAULT_TIMEOUT)

	defer cancel()

	var shortURL *ShortURLResponse

	err := repository.db.WithContext(ctx).
		Model(&models.URL{}).
		Select("urls.short_url").
		Where("urls.original_url= ?", url).
		First(&shortURL).
		Error
	if err != nil {
		return nil, err
	}

	if shortURL == nil {
		return nil, err
	}
	return shortURL, nil
}
