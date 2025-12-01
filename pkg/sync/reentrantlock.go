package sync

import "sync"

// ReentrantLock is a synchronisation primitive that allows a goroutine
// to acquire a lock multiple times without deadlocking.
// It maintains ownership information and supports reentrant locking by the same goroutine.
// The lock must be released as many times as it is acquired to fully unlock.
type ReentrantLock struct {
	mu      sync.Mutex
	cond    *sync.Cond
	ownerID int64
	count   int
}

// Lock acquires the lock for the current goroutine, allowing reentrant locking by the same goroutine.
func (l *ReentrantLock) Lock() {
	gid := GID()

	l.mu.Lock()
	defer l.mu.Unlock()

	// reentrant lock: same goroutine
	if l.ownerID == gid {
		l.count++

		return
	}

	// wait until it's unlocked
	for l.count != 0 {
		l.cond.Wait()
	}

	// first lock by this goroutine
	l.ownerID = gid
	l.count = 1
}

// Unlock releases the lock held by the current goroutine, decrementing the reentrant lock count. Panics if unowned.
func (l *ReentrantLock) Unlock() {
	gid := GID()

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.ownerID != gid || l.count == 0 {
		panic("unlock of unowned reentrant lock")
	}

	l.count--

	if l.count == 0 {
		l.ownerID = 0

		l.cond.Signal()
	}
}

// IsLocked reports whether the lock is currently held by any goroutine.
// It does not block and is safe for concurrent use.
func (l *ReentrantLock) IsLocked() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.count > 0
}

// ForceUnlock releases the lock regardless of ownership. This forcibly resets
// the lock to an unlocked state and wakes all waiting goroutines.
//
// WARNING: This should only be used in system-level recovery situations.
// It breaks the reentrancy contract and may violate invariants for goroutines
// currently relying on lock ownership.
func (l *ReentrantLock) ForceUnlock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Reset ownership and lock count
	l.ownerID = 0
	l.count = 0

	// Wake all waiters — they will contend normally
	l.cond.Broadcast()
}
