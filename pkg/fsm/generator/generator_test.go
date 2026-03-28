package generator

import (
	_ "embed"
	"encoding/json"
	"os"
	"testing"

	"github.com/codemityio/go/pkg/fsm/dsl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/workflow.json
var workflow []byte

//go:embed testdata/result.puml
var result []byte

// TestDefaultGenerate a Generator test.
func TestDefaultGenerate(t *testing.T) {
	g := New(WithColours("green", "red", "#f5f5f5", "#808080"))

	var w dsl.Workflow

	err := json.Unmarshal(workflow, &w)
	require.NoError(t, err)

	r, err := g.Generate(&w)
	require.NoError(t, err)

	assert.Equal(t, string(result), string(r))

	require.NoError(t, os.WriteFile("testdata/result.puml", r, 0o644)) // #nosec G306
}
