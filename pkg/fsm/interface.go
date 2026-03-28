// Package fsm contains finite state machine source code.
package fsm

import "github.com/codemityio/go/pkg/fsm/dsl"

// StateSetter defines the ability to transition the FSM to a new state by name.
// Returns an error if the provided state name is invalid or not part of the FSM definition.
type StateSetter interface {
	SetState(name string) error
}

// StateGetter defines the ability to retrieve the current active state of the FSM.
// Returns a pointer to the current state, or nil if no state has been set.
type StateGetter interface {
	GetState() *dsl.State
}

// InitialStateGetter returns the initial state.
type InitialStateGetter interface {
	GetInitialState() *dsl.State
}

// StatesGetter returns all states.
type StatesGetter interface {
	GetStates() []*dsl.State
}

// FinalStateGetter returns all the final states.
type FinalStateGetter interface {
	GetFinalStates() []*dsl.State
}

// PreviousStatesGetter returns all previous states.
type PreviousStatesGetter interface {
	GetPreviousStates() []*dsl.State
}

// NextStatesGetter defines the ability to retrieve all valid next states
// that can be transitioned to from the current state.
type NextStatesGetter interface {
	GetNextStates() []*dsl.State
}

// EdgesGetter returns all edges.
type EdgesGetter interface {
	GetEdges() []*dsl.Edge
}

// NextEdgesGetter defines the ability to retrieve all valid outgoing edges
// from the current state, representing possible transitions in the FSM.
type NextEdgesGetter interface {
	GetNextEdges() []*dsl.Edge
}

// FSM represents a finite state machine abstraction.
// It provides methods to set and get the current state, as well as query valid next states and transitions.
type FSM interface {
	StateSetter
	StateGetter
	NextStatesGetter
	NextEdgesGetter
}
