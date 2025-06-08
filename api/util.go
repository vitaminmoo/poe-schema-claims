package api

import (
	"fmt"

	"github.com/google/uuid"
)

func ID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("generating UUID: %w", err)
	}
	return id.String(), nil
}
