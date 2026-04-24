package url

import "gorm.io/gorm"

var _ URLRepository = (*urlGormRepository)(nil)

type urlGormRepository struct {
	db *gorm.DB
}

func NewURLGormRepository(db *gorm.DB) *urlGormRepository {
	return &urlGormRepository{db: db}
}

func (repository *urlGormRepository) FindByShortURL(code string) (*URLResponse, error) {}
func (repository *urlGormRepository) Save(url CreatURL) error {
}
