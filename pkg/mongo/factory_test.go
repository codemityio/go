package mongo

import (
	"testing"

	"github.com/codemityio/go/pkg/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MongoDBClient(t *testing.T) {
	cfg, err := env.New[Config]()
	require.NoError(t, err)

	client, err := NewClient(cfg)
	require.NoError(t, err)

	assert.NotNil(t, client)
}
