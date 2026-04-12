package huma

import (
	_ "embed"
	"testing"

	"github.com/codemityio/go/pkg/env"
	"github.com/codemityio/go/pkg/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/config.json
var configJSON []byte

func testConfig(t *testing.T) {
	t.Helper()

	config, err := env.New[Config]()
	require.NoError(t, err)

	assert.Equal(t, "API", config.Title)
	assert.Equal(t, "v0.0.0", config.Version)
	assert.Equal(t, "/openapi", config.OpenAPIPath)
	assert.Equal(t, "/docs", config.DocsPath)
	assert.Equal(t, "/schemas", config.SchemasPath)
	assert.Equal(t, "application/json", config.DefaultFormat)
}

func TestConfig_Env(t *testing.T) {
	t.Setenv("TITLE", "API")
	t.Setenv("VERSION", "v0.0.0")
	t.Setenv("OPEN_API_PATH", "/openapi")
	t.Setenv("DOCS_PATH", "/docs")
	t.Setenv("SCHEMAS_PATH", "/schemas")
	t.Setenv("DEFAULT_FORMAT", "application/json")

	testConfig(t)
}

func TestConfig_Env_Defaults(t *testing.T) {
	testConfig(t)
}

func TestConfig_JSON(t *testing.T) {
	config, err := json.NewJSON[Config](configJSON)
	require.NoError(t, err)

	assert.Equal(t, "API", config.Title)
	assert.Equal(t, "v0.0.0", config.Version)
	assert.Equal(t, "/openapi", config.OpenAPIPath)
	assert.Equal(t, "/docs", config.DocsPath)
	assert.Equal(t, "/schemas", config.SchemasPath)
	assert.Equal(t, "application/json", config.DefaultFormat)
}
