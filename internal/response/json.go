package response

import (
	"encoding/json"
	"net/http"
)

func JSON[T any](w http.ResponseWriter, status int, data T) error {
	return JSONWithHeader(w, status, data, nil)
}

func JSONWithHeader[T any](w http.ResponseWriter, status int, data T, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	w.Write(js)

	return nil
}
