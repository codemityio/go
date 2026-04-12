package mongo

import (
	"testing"

	"github.com/codemityio/go/pkg/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_LoadENV(t *testing.T) {
	t.Setenv("DSN", "mongodb://localhost:27017/database?retryWrites=false")
	t.Setenv("USERNAME", "username")
	t.Setenv("PASSWORD", "password")

	cfg, err := env.New[Config]()
	require.NoError(t, err)

	assert.Equal(t, "mongodb://localhost:27017/database?retryWrites=false", cfg.DSN)
	assert.Equal(t, "username", cfg.Username)
	assert.Equal(t, "password", cfg.Password)
}
