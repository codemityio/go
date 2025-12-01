package sync

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReentrantLock_BasicLockUnlock(t *testing.T) {
	lock := NewReentrantLock()

	require.NotPanics(t, func() {
		lock.Lock()
		lock.Unlock() //nolint:staticcheck
	})
}

func TestReentrantLock_Reentrancy(t *testing.T) {
	lock := NewReentrantLock()

	require.NotPanics(t, func() {
		lock.Lock()
		lock.Lock()   // reentrant
		lock.Unlock() //nolint:staticcheck
		lock.Unlock()
	})
}

func TestReentrantLock_SecondGoroutineBlocksUntilUnlocked(t *testing.T) {
	lock := NewReentrantLock()

	lock.Lock()

	ch := make(chan struct{})

	go func() {
		lock.Lock()
		lock.Unlock() //nolint:staticcheck

		ch <- struct{}{}
	}()

	select {
	case <-ch:
		t.Fatal("second goroutine should be blocked until first unlocks")
	case <-time.After(50 * time.Millisecond):
		// ok: still blocked
	}

	lock.Unlock() // now the second goroutine should proceed

	select {
	case <-ch:
		// passed
	case <-time.After(time.Second):
		t.Fatal("second goroutine did not proceed after unlock")
	}
}

func TestReentrantLock_PanicsOnWrongUnlock(t *testing.T) {
	lock := NewReentrantLock()

	require.Panics(t, func() {
		lock.Unlock()
	})
}

func TestReentrantLock_IsLocked(t *testing.T) {
	lock := NewReentrantLock()

	t.Run("initially unlocked", func(t *testing.T) {
		require.False(t, lock.IsLocked())
	})

	t.Run("locked once", func(t *testing.T) {
		lock.Lock()
		require.True(t, lock.IsLocked())
		lock.Unlock()
	})

	t.Run("reentrantly locked", func(t *testing.T) {
		lock.Lock()
		lock.Lock()
		require.True(t, lock.IsLocked())
		lock.Unlock()
		lock.Unlock()
	})

	t.Run("unlocked again", func(t *testing.T) {
		require.False(t, lock.IsLocked())
	})
}

func TestReentrantLock_ForceUnlock(t *testing.T) {
	lock := NewReentrantLock()

	// Lock once
	lock.Lock()

	// make sure it's locked
	require.True(t, lock.IsLocked(), "lock should be held before ForceUnlock")

	// a goroutine waiting on lock.Lock() will be released after ForceUnlock()
	wg := sync.WaitGroup{}
	wg.Add(1)

	started := make(chan struct{})
	finished := make(chan struct{})

	go func() {
		defer wg.Done()

		close(started) // goroutine started

		lock.Lock()   // should unblock after ForceUnlock
		lock.Unlock() //nolint:staticcheck

		close(finished)
	}()

	// wait for goroutine to be actually waiting
	<-started

	// give it a moment to block
	time.Sleep(50 * time.Millisecond)
	require.True(t, lock.IsLocked(), "lock should still be locked before force unlock")

	// act: forcibly unlock
	require.NotPanics(t, func() {
		lock.ForceUnlock()
	})

	// after ForceUnlock: lock must report unlocked
	require.False(t, lock.IsLocked(), "lock should be unlocked after force unlock")

	// waiting goroutine must proceed
	select {
	case <-finished:
		// ok
	case <-time.After(time.Second):
		t.Fatal("waiting goroutine did not resume after ForceUnlock")
	}

	// the lock must be usable again
	require.NotPanics(t, func() {
		lock.Lock()
		lock.Unlock() //nolint:staticcheck
	})
}
