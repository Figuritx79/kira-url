package url

import (
	"context"

	"kira-url/internal/database"
	"kira-url/internal/database/models"
	dtourl "kira-url/internal/dto/url"

	"gorm.io/gorm"
)

type URLRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (repository *URLRepository) FindByShortURL(code string) (*dtourl.URLResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), database.DEFAULT_TIMEOUT)

	defer cancel()

	var url *dtourl.URLResponse

	err := repository.db.WithContext(ctx).
		Model(&models.URL{}).
		Select("urls.original_url, urls.visit_count as visit_count, urls.id as id").
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

func (repository *URLRepository) Save(url *models.URL) error {
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

func (repository *URLRepository) FindByURL(url string) (*dtourl.ShortURLResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), database.DEFAULT_TIMEOUT)

	defer cancel()

	var shortURL *dtourl.ShortURLResponse

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

func (repository *URLRepository) Update(updateURL models.URL, code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), database.DEFAULT_TIMEOUT)
	defer cancel()

	err := repository.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&models.URL{}).Where("short_url= ?", code).Updates(updateURL).Error; err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

func (repository *URLRepository) Updates(urls []models.URL) error {
	ctx, cancel := context.WithTimeout(context.Background(), database.DEFAULT_TIMEOUT)
	defer cancel()

	err := repository.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			for _, url := range urls {
				if err := tx.Model(&models.URL{}).Where("short_url= ?", url.ShortURL).UpdateColumn("visit_count", gorm.Expr("visit_count + ?", url.VisitCount)).Error; err != nil {
					return err
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}
