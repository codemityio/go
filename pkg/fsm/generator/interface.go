// Package generator contains finite state machine diagram generator source code.
package generator

import "github.com/codemityio/go/pkg/fsm/dsl"

// Generator an FSM diagram generator interface.
type Generator interface {
	Generate(dsl *dsl.Workflow) ([]byte, error)
}
