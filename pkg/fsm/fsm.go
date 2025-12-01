package fsm

import (
	"fmt"
	"sort"
	"sync"

	"github.com/codemityio/go/pkg/fsm/dsl"
)

// FiniteStateMachine is a thread-safe finite state machine.
type FiniteStateMachine struct {
	mu       sync.RWMutex
	state    *dsl.State
	workflow *dsl.Workflow
}

// SetState sets the current state by name.
func (d *FiniteStateMachine) SetState(name string) error {
	if name == "" {
		return fmt.Errorf("%w: invalid state `%s`", ErrPkg, name)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	for i := range d.workflow.States {
		if name == d.workflow.States[i].Name {
			d.state = &d.workflow.States[i]

			return nil
		}
	}

	return fmt.Errorf("%w: invalid state `%s`", ErrPkg, name)
}

// GetState returns the current state.
func (d *FiniteStateMachine) GetState() *dsl.State {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.state
}

// GetInitialState finds initial state.
func (d *FiniteStateMachine) GetInitialState() *dsl.State {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// collect all destination states
	states := make(map[string]struct{})
	for _, t := range d.workflow.Edges {
		states[t.To] = struct{}{}
	}

	// the initial state is the one not found in any transition 'to'
	for _, s := range d.workflow.States {
		if _, state := states[s.Name]; !state {
			return &s
		}
	}

	return nil
}

// GetStates returns all the states.
func (d *FiniteStateMachine) GetStates() []*dsl.State {
	d.mu.RLock()
	defer d.mu.RUnlock()

	states := make([]*dsl.State, len(d.workflow.States))

	for i, s := range d.workflow.States {
		states[i] = &s
	}

	return states
}

// GetFinalStates returns all final states of the finite state machine. A state is final if no edges originate from it.
func (d *FiniteStateMachine) GetFinalStates() []*dsl.State {
	d.mu.RLock()
	defer d.mu.RUnlock()

	final := map[string]*dsl.State{}

	for _, v := range d.workflow.States {
		final[v.Name] = &v
	}

	for i := range d.workflow.Edges {
		delete(final, d.workflow.Edges[i].From)
	}

	keys := make([]string, 0, len(final))
	for k := range final {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	fin := make([]*dsl.State, len(keys))

	for i := range keys {
		fin[i] = final[keys[i]]
	}

	return fin
}

// GetPreviousStates returns the list of valid previous states that can transition into the current state.
func (d *FiniteStateMachine) GetPreviousStates() []*dsl.State { //nolint:dupl // accepted
	d.mu.RLock()
	defer d.mu.RUnlock()

	var states []*dsl.State

	if d.state == nil || d.workflow == nil {
		return states
	}

	for _, edge := range d.workflow.Edges {
		if d.state.Name == edge.To {
			for i := range d.workflow.States {
				if d.workflow.States[i].Name == edge.From {
					states = append(states, &d.workflow.States[i])
				}
			}
		}
	}

	return states
}

// GetNextStates returns the list of next valid states from the current state.
func (d *FiniteStateMachine) GetNextStates() []*dsl.State { //nolint:dupl // accepted
	d.mu.RLock()
	defer d.mu.RUnlock()

	var states []*dsl.State

	if d.state == nil || d.workflow == nil {
		return states
	}

	for _, edge := range d.workflow.Edges {
		if d.state.Name == edge.From {
			for i := range d.workflow.States {
				if d.workflow.States[i].Name == edge.To {
					states = append(states, &d.workflow.States[i])
				}
			}
		}
	}

	return states
}

// GetEdges returns all edges.
func (d *FiniteStateMachine) GetEdges() []*dsl.Edge {
	d.mu.RLock()
	defer d.mu.RUnlock()

	edges := make([]*dsl.Edge, len(d.workflow.Edges))

	for i, e := range d.workflow.Edges {
		edges[i] = &e
	}

	return edges
}

// GetNextEdges returns all edges that can be taken from the current state.
func (d *FiniteStateMachine) GetNextEdges() []*dsl.Edge {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var next []*dsl.Edge

	if d.state == nil || d.workflow == nil {
		return next
	}

	for _, edge := range d.workflow.Edges {
		if edge.From == d.state.Name {
			next = append(next, &edge)
		}
	}

	return next
}

// SetWorkflow sets the workflow and initialises the state to the workflow's initial state.
func (d *FiniteStateMachine) setWorkflow(workflow *dsl.Workflow) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.workflow = workflow
}
