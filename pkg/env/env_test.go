package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MyEnvConfig struct {
	Field string `env:"FIELD"`
}

func TestNewEnv(t *testing.T) {
	t.Setenv("FIELD", "field")

	cfg, err := NewEnv[MyEnvConfig]()
	require.NoError(t, err)

	assert.Equal(t, "field", cfg.Field)
}

func TestMustBuildEnv(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()

	t.Setenv("FIELD", "field")

	cfg := MustBuildEnv[MyEnvConfig]()

	assert.Equal(t, "field", cfg.Field)
}
