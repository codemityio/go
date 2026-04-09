package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MyEnvConfig struct {
	Field string `env:"FIELD"`
}

func TestNew(t *testing.T) {
	t.Setenv("FIELD", "field")

	cfg, err := New[MyEnvConfig]()
	require.NoError(t, err)

	assert.Equal(t, "field", cfg.Field)
}

func TestMustBuild(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()

	t.Setenv("FIELD", "field")

	cfg := MustBuild[MyEnvConfig]()

	assert.Equal(t, "field", cfg.Field)
}
