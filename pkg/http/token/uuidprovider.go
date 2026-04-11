package token

import (
	"fmt"

	"github.com/google/uuid"
)

// DefaultUUIDProvider provider.
type DefaultUUIDProvider struct{}

// RandomUUID return random UUID.
func (u *DefaultUUIDProvider) RandomUUID() (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%w: %w", ErrUUIDProviderNewRandom, err)
	}

	return id, nil
}
