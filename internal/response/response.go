package response

import (
	"fmt"
	"kira-url/internal/transport/httptransport"
	"net/http"
	"strings"
)

func OK[T any](w http.ResponseWriter, data *T, message string) {
	responseBody := httptransport.JSONResponse[T]{
		Message: message,
		Data:    data,
		Type:    httptransport.Success,
	}
	err := JSON(w, http.StatusOK, responseBody)
	if err != nil {
		w.WriteHeader(http.StatusOK)
	}
}

func Found[T any](w http.ResponseWriter, data *T, message, originalURL string) {
	headers := make(http.Header)
	headers["Location"] = []string{originalURL}
	body := httptransport.JSONResponse[T]{
		Message: message,
		Data:    data,
		Type:    httptransport.Success,
	}
	err := JSONWithHeader(w, http.StatusFound, body, headers)
	if err != nil {
		w.WriteHeader(http.StatusFound)
		w.Header().Set("Location", originalURL)
	}
}

func errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, header http.Header) {
	message = strings.ToLower(message[:1] + message[1:])

	body := httptransport.JSONResponse[any]{
		Type:    httptransport.Error,
		Message: message,
	}
	err := JSONWithHeader(w, status, body, header)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerError(w http.ResponseWriter, r *http.Request) {
	errorMessage(w, r, http.StatusInternalServerError, "internal server error", nil)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method  is not supported for this resource", r.Method)
	errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	errorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	errorMessage(w, r, http.StatusNotFound, message, nil)
}
