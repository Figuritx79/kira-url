package url

type CreatURL struct {
	OriginalURL string `json:"original_url"`
	CustomCode  string `json:"custom_code,omitempty"`
	ShortURL    string `json:"-"`
}
