package url

import "errors"

var (
	ErrInvalidURL         = errors.New("invalid URL")
	ErrURLNotFound        = errors.New("URL not found")
	ErrMinRunesCustomCode = errors.New("slug min 6 characters")
	ErrMaxRunesCustomCode = errors.New("slug max 40 characters")
	ErrInvalidCustomCode  = errors.New("invalid slug code")
)
