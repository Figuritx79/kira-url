package request

import (
	"encoding/json"
	"net/http"
)

func GetRequestBody[B any](r *http.Request, body *B) error {
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return err
	}
	return nil
}
