//nolint:dupl // false positive
package sync

import "sync"

// SafeIndex is a thread-safe index with uint64 indexes.
type SafeIndex struct {
	indexes map[uint64]struct{} // internal set for deduplication
	list    []uint64            // cached list of indexes
	data    map[uint64]any
	mutex   sync.RWMutex
}

// Set sets a value under a specific uint64 index.
func (si *SafeIndex) Set(index uint64, value any) {
	si.mutex.Lock()
	defer si.mutex.Unlock()

	if _, exists := si.data[index]; !exists {
		si.list = append(si.list, index)
		si.indexes[index] = struct{}{}
	}

	si.data[index] = value
}

// Unset removes the index and its value from the index.
func (si *SafeIndex) Unset(index uint64) {
	si.mutex.Lock()
	defer si.mutex.Unlock()

	if _, exists := si.data[index]; exists {
		delete(si.data, index)
		delete(si.indexes, index)

		// rebuild list (cheap since only on delete)
		newList := make([]uint64, 0, len(si.indexes))
		for idx := range si.indexes {
			newList = append(newList, idx)
		}

		si.list = newList
	}
}

// Get returns the value stored under the index.
func (si *SafeIndex) Get(index uint64) (any, bool) {
	si.mutex.RLock()
	defer si.mutex.RUnlock()

	value, ok := si.data[index]

	return value, ok
}

// Indexes returns all the used uint64 indexes.
func (si *SafeIndex) Indexes() []uint64 {
	si.mutex.RLock()
	defer si.mutex.RUnlock()

	result := make([]uint64, len(si.list))
	copy(result, si.list)

	return result
}

// HasIndex checks if the index exists in the index.
func (si *SafeIndex) HasIndex(index uint64) bool {
	si.mutex.RLock()
	defer si.mutex.RUnlock()

	_, ok := si.data[index]

	return ok
}
