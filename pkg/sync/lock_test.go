package sync

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLock_BasicLockUnlock(t *testing.T) {
	lock := NewLock()

	require.False(t, lock.IsLocked(), "lock should initially be unlocked")

	lock.Lock()
	require.True(t, lock.IsLocked(), "lock should be held after Lock()")

	lock.Unlock()
	require.False(t, lock.IsLocked(), "lock should be released after Unlock()")
}

func TestLock_UnlockPanicsWhenNotHeld(t *testing.T) {
	lock := NewLock()

	require.Panics(t, func() {
		lock.Unlock()
	}, "unlocking an unlocked lock should panic")
}

func TestLock_BlocksUntilUnlocked(t *testing.T) {
	lock := NewLock()

	lock.Lock() // hold the lock

	ch := make(chan struct{})

	go func() {
		lock.Lock()   // should block
		lock.Unlock() //nolint:staticcheck // empty critical section is intentional, testing lock acquisition

		ch <- struct{}{}
	}()

	select {
	case <-ch:
		t.Fatal("goroutine should be blocked until lock is released")
	case <-time.After(50 * time.Millisecond):
		// OK: still blocked
	}

	lock.Unlock() // release the lock

	select {
	case <-ch:
		// passed: goroutine acquired lock and unlocked
	case <-time.After(time.Second):
		t.Fatal("goroutine did not proceed after unlock")
	}
}

func TestLock_ForceUnlock(t *testing.T) {
	lock := NewLock()

	lock.Lock()
	require.True(t, lock.IsLocked(), "lock should be held")

	lock.ForceUnlock()
	require.False(t, lock.IsLocked(), "lock should be released after ForceUnlock()")

	// now another goroutine can acquire the lock
	ch := make(chan struct{})

	go func() {
		lock.Lock()
		lock.Unlock() //nolint:staticcheck // empty critical section is intentional, testing lock acquisition

		ch <- struct{}{}
	}()

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("goroutine could not acquire lock after ForceUnlock")
	}
}

func TestLock_IsLocked(t *testing.T) {
	lock := NewLock()

	require.False(t, lock.IsLocked(), "lock should initially be unlocked")

	lock.Lock()
	require.True(t, lock.IsLocked(), "lock should be held after Lock()")

	lock.Unlock()
	require.False(t, lock.IsLocked(), "lock should be released after Unlock()")
}

func TestLock_ConcurrentLocking(t *testing.T) {
	lock := NewLock()

	var (
		wg      sync.WaitGroup
		counter int
	)

	const goroutines = 5
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()

			lock.Lock()
			tmp := counter

			time.Sleep(10 * time.Millisecond) // simulate work

			counter = tmp + 1
			lock.Unlock()
		}()
	}

	wg.Wait()
	require.Equal(t, goroutines, counter, "all goroutines should increment counter safely")
}
