package token

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultUUIDProvider_RandomUUID(t *testing.T) {
	provider := &DefaultUUIDProvider{}

	id, err := provider.RandomUUID()
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

func TestDefaultUUIDProvider_IsUUIDProviderUnableToGenerateUUID(t *testing.T) {
	errBoom := errors.New("error")

	tests := map[string]struct {
		err  error
		want bool
	}{
		"true for the bare sentinel": {
			err:  ErrUUIDProviderUnableToGenerateUUID,
			want: true,
		},
		"true when wrapped": {
			err:  fmt.Errorf("%w: %w", ErrUUIDProviderUnableToGenerateUUID, errBoom),
			want: true,
		},
		"false for an unrelated error": {
			err:  errBoom,
			want: false,
		},
		"false for nil": {
			err:  nil,
			want: false,
		},
	}

	provider := &DefaultUUIDProvider{}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, test.want, provider.IsUUIDProviderUnableToGenerateUUID(test.err))
		})
	}
}
