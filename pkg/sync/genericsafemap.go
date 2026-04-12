//nolint:dupl
package sync

import "sync"

// GenericSafeMap thread generic safe map.
type GenericSafeMap[T any] struct {
	keys  map[string]struct{} // internal set for deduplication
	list  []string            // cached list of keys
	data  map[string]T
	mutex sync.RWMutex
}

// Set sets value under a specific key.
func (smg *GenericSafeMap[T]) Set(key string, value T) {
	smg.mutex.Lock()
	defer smg.mutex.Unlock()

	if _, exists := smg.data[key]; !exists {
		smg.list = append(smg.list, key)
		smg.keys[key] = struct{}{}
	}

	smg.data[key] = value
}

// Unset unsets the key.
func (smg *GenericSafeMap[T]) Unset(key string) {
	smg.mutex.Lock()
	defer smg.mutex.Unlock()

	if _, exists := smg.data[key]; exists {
		delete(smg.data, key)
		delete(smg.keys, key)

		// rebuild list (cheap since only on delete)
		newList := make([]string, 0, len(smg.keys))

		for k := range smg.keys {
			newList = append(newList, k)
		}

		smg.list = newList
	}
}

// Get returns the value stored in under the key.
func (smg *GenericSafeMap[T]) Get(key string) (T, bool) { //nolint:ireturn
	smg.mutex.RLock()
	defer smg.mutex.RUnlock()

	value, ok := smg.data[key]

	return value, ok
}

// Keys lists all the used keys.
func (smg *GenericSafeMap[T]) Keys() []string {
	smg.mutex.RLock()
	defer smg.mutex.RUnlock()

	result := make([]string, len(smg.list))

	copy(result, smg.list)

	return result
}

// HasKey checks the key exists.
func (smg *GenericSafeMap[T]) HasKey(key string) bool {
	smg.mutex.RLock()
	defer smg.mutex.RUnlock()

	_, ok := smg.data[key]

	return ok
}
