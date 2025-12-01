package sync

import "sync"

// SafeListIndex indexed list.
type SafeListIndex[T any] struct {
	mu   sync.RWMutex
	data map[uint64][]T
}

// Append appends to a list.
func (s *SafeListIndex[T]) Append(index uint64, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[index] = append(s.data[index], value)
}

// Get returns a list.
func (s *SafeListIndex[T]) Get(index uint64) []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data[index]
}

// Exists checks if the list exists.
func (s *SafeListIndex[T]) Exists(index uint64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[index]

	return ok
}

// Delete deletes a list.
func (s *SafeListIndex[T]) Delete(index uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, index)
}
