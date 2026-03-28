package sync

import (
	"bytes"
	"sync"
)

// NewGate creates a new open gate.
func NewGate() *Gate {
	gate := &Gate{
		mu:     sync.RWMutex{},
		open:   nil,
		opened: false,
	}

	gate.openGate()

	return gate
}

// NewLock creates a new Lock.
func NewLock() *Lock {
	lock := &Lock{
		mu:   sync.Mutex{},
		cond: nil,
		held: false,
	}

	lock.cond = sync.NewCond(&lock.mu)

	return lock
}

// NewReentrantLock initialises and returns a new instance of ReentrantLock, supporting reentrant capabilities.
func NewReentrantLock() *ReentrantLock {
	lock := &ReentrantLock{
		mu:      sync.Mutex{},
		cond:    nil,
		ownerID: 0,
		count:   0,
	}

	lock.cond = sync.NewCond(&lock.mu)

	return lock
}

// NewGenericSafeIndex factory function.
func NewGenericSafeIndex[T any]() *GenericSafeIndex[T] {
	return &GenericSafeIndex[T]{
		indexes: make(map[uint64]struct{}),
		list:    make([]uint64, 0),
		data:    make(map[uint64]T),
		mutex:   sync.RWMutex{},
	}
}

// NewGenericSafeMap factory function.
func NewGenericSafeMap[T any]() *GenericSafeMap[T] {
	return &GenericSafeMap[T]{
		keys:  make(map[string]struct{}),
		list:  make([]string, 0),
		data:  make(map[string]T),
		mutex: sync.RWMutex{},
	}
}

// NewSafeBuffer factory function.
func NewSafeBuffer() *SafeBuffer {
	return &SafeBuffer{
		mu: sync.Mutex{},
		b:  bytes.Buffer{},
	}
}

// NewSafeIndex factory function.
func NewSafeIndex() *SafeIndex {
	return &SafeIndex{
		indexes: make(map[uint64]struct{}),
		list:    make([]uint64, 0),
		data:    make(map[uint64]any),
		mutex:   sync.RWMutex{},
	}
}

// NewSafeListIndex factory function.
func NewSafeListIndex[T any]() *SafeListIndex[T] {
	return &SafeListIndex[T]{
		mu:   sync.RWMutex{},
		data: make(map[uint64][]T),
	}
}

// NewSafeMap factory function.
func NewSafeMap() *SafeMap {
	return &SafeMap{
		keys:  make(map[string]struct{}),
		list:  make([]string, 0),
		data:  make(map[string]any),
		mutex: sync.RWMutex{},
	}
}
