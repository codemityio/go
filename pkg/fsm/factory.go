package fsm

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/codemityio/go/pkg/fsm/dsl"
)

// New is a factory function.
func New(workflow *dsl.Workflow) *FiniteStateMachine {
	fsm := FiniteStateMachine{
		mu:       sync.RWMutex{},
		state:    nil,
		workflow: nil,
	}

	fsm.setWorkflow(workflow)

	return &fsm
}

// NewWorkflow a workflow factory.
func NewWorkflow(input []byte) (*dsl.Workflow, error) {
	var workflow dsl.Workflow

	if e := json.Unmarshal(input, &workflow); e != nil {
		return nil, fmt.Errorf("%w: unable to unmarshal input: %w", ErrUnmarshal, e)
	}

	return &workflow, nil
}
