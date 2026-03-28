//nolint:dupl // false positive
package sync

import "sync"

// SafeMap thread safe map.
type SafeMap struct {
	keys  map[string]struct{} // internal set for deduplication
	list  []string            // cached list of keys
	data  map[string]any
	mutex sync.RWMutex
}

// Set sets value under a specific key.
func (sm *SafeMap) Set(key string, value any) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if _, exists := sm.data[key]; !exists {
		sm.list = append(sm.list, key)
		sm.keys[key] = struct{}{}
	}

	sm.data[key] = value
}

// Unset unsets the key.
func (sm *SafeMap) Unset(key string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if _, exists := sm.data[key]; exists {
		delete(sm.data, key)
		delete(sm.keys, key)

		// rebuild list (cheap since only on delete)
		newList := make([]string, 0, len(sm.keys))

		for k := range sm.keys {
			newList = append(newList, k)
		}

		sm.list = newList
	}
}

// Get returns the value stored in under the key.
func (sm *SafeMap) Get(key string) (any, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	value, ok := sm.data[key]

	return value, ok
}

// Keys lists all the used keys.
func (sm *SafeMap) Keys() []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	result := make([]string, len(sm.list))

	copy(result, sm.list)

	return result
}

// HasKey checks the key exists.
func (sm *SafeMap) HasKey(key string) bool {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	_, ok := sm.data[key]

	return ok
}
