package fsm

import (
	"testing"

	"github.com/codemityio/go/pkg/fsm/dsl"
	"github.com/stretchr/testify/require"
)

func TestNewWorkflow(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		expectedError  error
		expectedResult *dsl.Workflow
	}{
		{
			name: "valid workflow with one state and one edge",
			input: []byte(`{
				"states": [
					{"name": "Created", "description": "initial", "colour": "blue"}
				],
				"edges": [
					{"from": "Created", "to": "Ready", "description": "transition"}
				]
			}`),
			expectedError: nil,
			expectedResult: &dsl.Workflow{
				States: []dsl.State{
					{Name: "Created", Description: "initial", Colour: "blue"},
				},
				Edges: []dsl.Edge{
					{From: "Created", To: "Ready", Description: "transition"},
				},
			},
		},
		{
			name:           "invalid JSON input",
			input:          []byte(`{ invalid json }`),
			expectedError:  ErrUnmarshal,
			expectedResult: nil,
		},
		{
			name:           "empty JSON input",
			input:          []byte(``),
			expectedError:  ErrUnmarshal,
			expectedResult: nil,
		},
		{
			name:          "valid empty workflow",
			input:         []byte(`{"states": [], "edges": []}`),
			expectedError: nil,
			expectedResult: &dsl.Workflow{
				States: []dsl.State{},
				Edges:  []dsl.Edge{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wf, err := NewWorkflow(tt.input)

			require.ErrorIs(t, err, tt.expectedError)

			require.Equal(t, tt.expectedResult, wf)
		})
	}
}
