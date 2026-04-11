package server

import (
	_ "embed"
	"testing"
	"time"

	"github.com/codemityio/go/pkg/env"
	"github.com/codemityio/go/pkg/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/config.json
var configJSON []byte

func TestConfig_Env(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("READ_TIMEOUT", "15s")
	t.Setenv("WRITE_TIMEOUT", "15s")
	t.Setenv("ERRORS_BUFFER_SIZE", "100")
	t.Setenv("TLS_CERT", "cert")
	t.Setenv("TLS_KEY", "key")
	t.Setenv("TLS_CA_CERT", "ca_cert")

	config, err := env.New[Config]()
	require.NoError(t, err)

	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, time.Duration(15000000000), config.ReadTimeout)
	assert.Equal(t, time.Duration(15000000000), config.WriteTimeout)
	assert.Equal(t, 100, config.ErrorsBufferSize)
	assert.Equal(t, "cert", config.TLS.Cert)
	assert.Equal(t, "key", config.TLS.Key)
	assert.Equal(t, "ca_cert", config.TLS.CACert)
}

func TestConfig_JSON(t *testing.T) {
	config, err := json.NewJSON[Config](configJSON)
	require.NoError(t, err)

	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, time.Duration(15000000000), config.ReadTimeout)
	assert.Equal(t, time.Duration(15000000000), config.WriteTimeout)
	assert.Equal(t, 100, config.ErrorsBufferSize)
	assert.Equal(t, "cert", config.TLS.Cert)
	assert.Equal(t, "key", config.TLS.Key)
	assert.Equal(t, "ca_cert", config.TLS.CACert)
}
