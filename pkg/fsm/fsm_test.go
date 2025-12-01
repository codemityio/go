package fsm

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/codemityio/go/pkg/fsm/dsl"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/workflow.json
var workflow []byte

func TestFiniteStateMachine(t *testing.T) {
	type testCase struct {
		state          string
		states         []*dsl.State
		edges          []*dsl.Edge
		initialState   *dsl.State
		finalStates    []*dsl.State
		currentState   *dsl.State
		previousStates []*dsl.State
		nextStates     []*dsl.State
		nextEdges      []*dsl.Edge
		expectedError  error
	}

	states := []*dsl.State{
		{
			Name:        "zero",
			Description: "state zero description",
		},
		{
			Name:        "one",
			Description: "state one description",
		},
		{
			Name:        "two",
			Description: "state two\\nmulti line description",
		},
		{
			Name: "three",
		},
		{
			Name: "four",
		},
		{
			Name:        "five",
			Description: "state five description",
		},
		{
			Name: "six",
		},
	}

	edges := []*dsl.Edge{
		{
			From: "zero",
			To:   "one",
		},
		{
			From: "zero",
			To:   "two",
		},
		{
			From:        "zero",
			To:          "three",
			Description: "edge zero-three\\nmulti line description",
		},
		{
			From: "one",
			To:   "four",
		},
		{
			From:        "two",
			To:          "four",
			Description: "edge two-four description",
		},
		{
			From: "three",
			To:   "five",
		},
		{
			From: "three",
			To:   "six",
		},
		{
			From: "six",
			To:   "three",
		},
	}

	initialState := &dsl.State{
		Name:        "zero",
		Description: "state zero description",
	}

	finalStates := []*dsl.State{
		{
			Name:        "five",
			Description: "state five description",
		},
		{
			Name: "four",
		},
	}

	tests := map[string]testCase{
		"one": {
			state:        "",
			states:       states,
			edges:        edges,
			initialState: initialState,
			finalStates:  finalStates,
			currentState: &dsl.State{
				Description: "state zero description",
				Name:        "zero",
			},
			previousStates: nil,
			nextStates: []*dsl.State{
				{Name: "one", Description: "state one description"},
				{Name: "two", Description: "state two\\nmulti line description"},
				{Name: "three", Description: ""},
			},
			nextEdges: []*dsl.Edge{
				{From: "zero", To: "one"},
				{From: "zero", To: "two"},
				{
					From:        "zero",
					To:          "three",
					Description: "edge zero-three\\nmulti line description",
				},
			},
			expectedError: fmt.Errorf("%w: invalid state `%s`", ErrPkg, ""),
		},
		"two": {
			state:        "zero",
			states:       states,
			edges:        edges,
			initialState: initialState,
			finalStates:  finalStates,
			currentState: &dsl.State{
				Description: "state zero description",
				Name:        "zero",
			},
			previousStates: nil,
			nextStates: []*dsl.State{
				{Name: "one", Description: "state one description"},
				{Name: "two", Description: "state two\\nmulti line description"},
				{Name: "three", Description: ""},
			},
			nextEdges: []*dsl.Edge{
				{From: "zero", To: "one"},
				{From: "zero", To: "two"},
				{
					From:        "zero",
					To:          "three",
					Description: "edge zero-three\\nmulti line description",
				},
			},
			expectedError: nil,
		},
		"three": {
			state:        "one",
			states:       states,
			edges:        edges,
			initialState: initialState,
			finalStates:  finalStates,
			currentState: &dsl.State{
				Description: "state one description",
				Name:        "one",
			},
			previousStates: []*dsl.State{
				{Description: "state zero description", Name: "zero"},
			},
			nextStates: []*dsl.State{
				{Name: "four"},
			},
			nextEdges: []*dsl.Edge{
				{From: "one", To: "four"},
			},
			expectedError: nil,
		},
		"four": {
			state:        "three",
			states:       states,
			edges:        edges,
			initialState: initialState,
			finalStates:  finalStates,
			currentState: &dsl.State{
				Name: "three",
			},
			previousStates: []*dsl.State{
				{Description: "state zero description", Name: "zero"},
				{Name: "six"},
			},
			nextStates: []*dsl.State{
				{Name: "five", Description: "state five description"},
				{Name: "six"},
			},
			nextEdges: []*dsl.Edge{
				{From: "three", To: "five"},
				{From: "three", To: "six"},
			},
			expectedError: nil,
		},
		"five": {
			state:        "four",
			states:       states,
			edges:        edges,
			initialState: initialState,
			finalStates:  finalStates,
			currentState: &dsl.State{
				Name: "four",
			},
			previousStates: []*dsl.State{
				{Name: "one", Description: "state one description"},
				{Name: "two", Description: "state two\\nmulti line description"},
			},
			nextStates:    nil,
			nextEdges:     nil,
			expectedError: nil,
		},
	}

	var w dsl.Workflow
	require.NoError(t, json.Unmarshal(workflow, &w))

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fsm := New(&w)

			initial := fsm.GetInitialState()
			require.NotNil(t, initial)

			require.NoError(t, fsm.SetState(initial.Name))

			var wg sync.WaitGroup

			const goroutines = 10

			for range goroutines {
				wg.Go(func() {
					err := fsm.SetState(test.state)
					require.Equal(t, test.expectedError, err)
					require.Equal(t, test.states, fsm.GetStates())
					require.Equal(t, test.currentState, fsm.GetState())
					require.Equal(t, test.previousStates, fsm.GetPreviousStates())
					require.Equal(t, test.nextStates, fsm.GetNextStates())
					require.Equal(t, test.nextEdges, fsm.GetNextEdges())
					require.Equal(t, test.initialState, fsm.GetInitialState())
					require.Equal(t, test.finalStates, fsm.GetFinalStates())
				})
			}

			wg.Wait()
		})
	}
}
