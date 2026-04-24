package url

import "log/slog"

type urlService struct {
	repository *urlGormRepository
	log        *slog.Logger
}

func NewURLService(repository *urlGormRepository, log *slog.Logger) *urlService {
	return &urlService{repository: repository, log: log}
}

func (service *urlService) FindByShortURL(code string) (*URLResponse, error) {
	return service.repository.FindByShortURL(code)
}

func (service *urlService) Save(url *CreatURL) (*URLCompleteResponse, error) {
	// return service.repository.Save(url)
}
