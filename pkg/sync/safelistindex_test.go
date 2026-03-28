package sync

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafeListIndex(t *testing.T) {
	list := NewSafeListIndex[string]()

	// add multiple values under the same index
	list.Append(1, "one-a")
	list.Append(1, "one-b")
	list.Append(2, "two-a")

	// get values for index 1
	values := list.Get(1)
	require.Len(t, values, 2)
	require.Contains(t, values, "one-a")
	require.Contains(t, values, "one-b")

	// get values for index 2
	values = list.Get(2)
	require.Len(t, values, 1)
	require.Equal(t, "two-a", values[0])

	// check existence
	require.True(t, list.Exists(1))
	require.False(t, list.Exists(99))

	// delete index
	list.Delete(1)
	require.False(t, list.Exists(1))
	require.Empty(t, list.Get(1))
}

func TestSafeListIndex_ConcurrentAccess(t *testing.T) {
	list := NewSafeListIndex[int]()

	var wg sync.WaitGroup

	const (
		index = 42
		total = 1000
	)

	// spawn multiple goroutines to append concurrently

	for i := range total {
		wg.Add(1)

		go func(val int) {
			defer wg.Done()

			list.Append(index, val)
		}(i)
	}

	wg.Wait()

	values := list.Get(index)
	require.Len(t, values, total)

	// check all values are unique and present
	seen := make(map[int]struct{})
	for _, v := range values {
		seen[v] = struct{}{}
	}

	require.Len(t, seen, total)
}
