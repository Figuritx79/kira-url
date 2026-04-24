package url

type URLRepository interface {
	FindByShortURL(code string) (*URLResponse, error)
	Save(url CreatURL) error
	// Update(url Creat, id uuid.UUID) error
}
