package httperrors

import (
	"fmt"
	"net/http"
	"strings"

	"kira-url/internal/response"
	"kira-url/internal/transport/httptransport"
)

func errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, header http.Header) {
	message = strings.ToLower(message[:1] + message[1:])

	body := httptransport.JSONResponse[any]{
		Type:    httptransport.Error,
		Message: message,
	}
	err := response.JSONWithHeader[httptransport.JSONResponse[any]](w, status, body, header)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	errorMessage(w, r, http.StatusInternalServerError, err.Error(), nil)
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
