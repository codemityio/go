package sync

import "sync"

// Lock is a normal non-reentrant lock.
type Lock struct {
	mu   sync.Mutex
	cond *sync.Cond
	held bool
}

// Lock acquires the lock, blocking if it is already held.
func (l *Lock) Lock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for l.held {
		l.cond.Wait()
	}

	l.held = true
}

// Unlock releases the lock. Panics if it is not held.
func (l *Lock) Unlock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.held {
		panic("unlock of unlocked lock")
	}

	l.held = false
	l.cond.Signal()
}

// IsLocked reports whether the lock is currently held.
func (l *Lock) IsLocked() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.held
}

// ForceUnlock releases the lock regardless of who holds it.
func (l *Lock) ForceUnlock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.held = false
	l.cond.Broadcast()
}
