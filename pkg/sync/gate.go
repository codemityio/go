package sync

import (
	"sync"
)

// Gate gatekeeper.
type Gate struct {
	mu     sync.RWMutex
	open   chan struct{}
	opened bool
}

// Wait blocks until the gate is open.
func (g *Gate) Wait() {
	g.mu.RLock()
	open := g.open
	g.mu.RUnlock()
	<-open
}

// Open marks the gate as open, allowing any waiting routines to proceed. This method is idempotent.
func (g *Gate) Open() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.opened {
		return
	}

	close(g.open)
	g.opened = true
}

// Close closes the gate, making Wait block; it resets the gate if it is currently open. It is idempotent.
func (g *Gate) Close() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.opened {
		return
	}

	g.openGate() // reset
}

// IsClosed checks if the gate is currently closed and returns true if it is, otherwise false.
func (g *Gate) IsClosed() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return !g.opened
}

func (g *Gate) openGate() {
	g.open = make(chan struct{})
	g.opened = false
}
