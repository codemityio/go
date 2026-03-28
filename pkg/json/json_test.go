package json

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MyJSONConfig struct {
	Field string `json:"field"`
}

func TestNewJSON(t *testing.T) {
	cfg, err := NewJSON[MyJSONConfig](json.RawMessage(`{"field":"field"}`))
	require.NoError(t, err)

	assert.Equal(t, "field", cfg.Field)
}

func TestMustBuildJSON(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()

	cfg := MustBuildJSON[MyJSONConfig](json.RawMessage(`{"field":"field"}`))

	assert.Equal(t, "field", cfg.Field)
}
