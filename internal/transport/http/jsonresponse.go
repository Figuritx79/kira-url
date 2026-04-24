package http

type ResponseType string

const (
	Success ResponseType = "success"
	Error   ResponseType = "error"
	Warn    ResponseType = "warn"
)

type JSONResponse[T any] struct {
	Type    ResponseType `json:"type"`
	Data    T            `json:"data"`
	Message string       `json:"message,omitempty"`
}
