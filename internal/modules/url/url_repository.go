package url

import "github.com/google/uuid"

type URLRepository interface {
	FindByShortCode(code string) (*URLResponse, error)
	Save(url CreatURL) error
	Update(url URL, id uuid.UUID) error
}
