package sync

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGate_InitialState(t *testing.T) {
	g := NewGate()

	assert.True(t, g.IsClosed(), "Gate should be closed initially")

	waiter := make(chan struct{})

	go func() {
		g.Wait()
		close(waiter)
	}()

	select {
	case <-waiter:
		t.Fatal("Wait() should block when gate is closed")
	case <-time.After(50 * time.Millisecond):
		// expected
	}
}

func TestGate_OpenUnblocksWait(t *testing.T) {
	gate := NewGate()

	waiter := make(chan struct{})

	go func() {
		gate.Wait()
		close(waiter)
	}()

	time.Sleep(20 * time.Millisecond) // ensure goroutine is waiting
	gate.Open()

	select {
	case <-waiter:
		// ✅ expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait() should unblock after Open()")
	}

	assert.False(t, gate.IsClosed(), "Gate should be open after Open()")
}

func TestGate_CloseBlocksAgain(t *testing.T) {
	gate := NewGate()
	gate.Open()

	assert.False(t, gate.IsClosed(), "Gate should be open")

	gate.Close()

	assert.True(t, gate.IsClosed(), "Gate should be closed after Close()")

	waiter := make(chan struct{})

	go func() {
		gate.Wait()
		close(waiter)
	}()

	select {
	case <-waiter:
		t.Fatal("Wait() should block after Close()")
	case <-time.After(50 * time.Millisecond):
		// ✅ expected
	}
}

func TestGate_OpenIsIdempotent(t *testing.T) {
	gate := NewGate()
	gate.Open()
	assert.False(t, gate.IsClosed(), "Gate should be open")

	assert.NotPanics(t, func() {
		gate.Open()
	}, "Open() should be idempotent")
}

func TestGate_CloseIsIdempotent(t *testing.T) {
	gate := NewGate()

	// Closed by default, so first Close() is no-op
	assert.True(t, gate.IsClosed(), "Gate should be closed")

	assert.NotPanics(t, func() {
		gate.Close()
	}, "Close() should be idempotent")

	// Open, then Close again
	gate.Open()
	gate.Close()

	assert.True(t, gate.IsClosed(), "Gate should be closed again")
}

func TestGate_ConcurrentWaiters(t *testing.T) {
	gate := NewGate()

	const nw = 5

	waiters := make(chan struct{}, nw)

	for range nw {
		go func() {
			gate.Wait()

			waiters <- struct{}{}
		}()
	}

	time.Sleep(20 * time.Millisecond)

	// No waiters should have passed yet
	select {
	case <-waiters:
		t.Fatal("Wait() should block when gate is closed")
	default:
		// good
	}

	gate.Open()

	timeout := time.After(100 * time.Millisecond)

	for i := range nw {
		select {
		case <-waiters:
			// good
		case <-timeout:
			t.Fatalf("Waiter %d did not unblock in time", i)
		}
	}
}
