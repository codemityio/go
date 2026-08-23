package token

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// DefaultUUIDProvider provider.
type DefaultUUIDProvider struct{}

// IsUUIDProviderUnableToGenerateUUID reports whether err is or wraps the error RandomUUID returns
// when it fails to generate a new UUID.
func (u *DefaultUUIDProvider) IsUUIDProviderUnableToGenerateUUID(err error) bool {
	return errors.Is(err, ErrUUIDProviderUnableToGenerateUUID)
}

// RandomUUID return random UUID.
func (u *DefaultUUIDProvider) RandomUUID() (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%w: %w", ErrUUIDProviderUnableToGenerateUUID, err)
	}

	return id, nil
}
