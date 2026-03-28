package sync

// Gatekeeper concurrency safe gate keeper.
type Gatekeeper interface {
	Wait()
	Open()
	Close()
	IsClosed() bool
}

// Locker represents an interface for acquiring and releasing locks to control access to shared resources.
type Locker interface {
	Lock()
	Unlock()
	IsLocked() bool
	ForceUnlock()
}

// SafeIndexer concurrency safe index.
type SafeIndexer interface {
	Set(index uint64, value any)
	Unset(index uint64)
	Get(index uint64) (any, bool)
	Indexes() []uint64
	HasIndex(index uint64) bool
}

// GenericSafeIndexer concurrency safe index.
type GenericSafeIndexer[T any] interface {
	Set(index uint64, value T)
	Unset(index uint64)
	Get(index uint64) (T, bool)
	Indexes() []uint64
	HasIndex(index uint64) bool
}

// GenericSafeListIndexer concurrency safe list index.
type GenericSafeListIndexer[T any] interface {
	Append(index uint64, value T)
	Get(index uint64) []T
	Exists(index uint64) bool
	Delete(index uint64)
}

// SafeMapper concurrency safe map.
type SafeMapper interface {
	Set(key string, value any)
	Unset(key string)
	Get(key string) (any, bool)
	Keys() []string
	HasKey(key string) bool
}

// GenericSafeMapper concurrency generic safe map.
type GenericSafeMapper[T any] interface {
	Set(key string, value T)
	Unset(key string)
	Get(key string) (T, bool)
	Keys() []string
	HasKey(key string) bool
}

// ReentrantLocker is a synchronisation primitive that allows a goroutine
// to acquire a lock multiple times without deadlocking.
// It maintains ownership information and supports reentrant locking by the same goroutine.
// The lock must be released as many times as it is acquired to fully unlock.
type ReentrantLocker interface {
	Lock()
	Unlock()
	IsLocked() bool
}
