package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeProvider_Now(t *testing.T) {
	t.Parallel()

	fixed := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		provider ProviderFunc
	}{
		{
			name:     "nil provider uses time.Now",
			provider: nil,
		},
		{
			name: "custom provider returns fixed time",
			provider: ProviderFunc(func() time.Time {
				return fixed
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			before := time.Now()
			got := tt.provider.Now()
			after := time.Now()

			if tt.provider == nil {
				// ensure time is within execution window
				require.False(t, got.Before(before))
				require.False(t, got.After(after))

				return
			}

			require.Equal(t, fixed, got)
		})
	}
}
