package helpers

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// NewUUID returns new UUID
func NewUUID() (string, error) {
	rUUID, err := uuid.NewV4()
	if err != nil {
		return "", errors.WithMessage(err, "Couldn't generate new UUID: %s")
	}

	return rUUID.String(), nil
}
