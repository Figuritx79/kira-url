package funcs

import (
	"github.com/google/uuid"
)

func GenerateUUID() (uuid.UUID, error) {
	newUUID, err := uuid.NewV7()
	return newUUID, err
}
