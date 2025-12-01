package sync

import (
	"encoding/json"
	"fmt"
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeIndex(t *testing.T) {
	si := NewSafeIndex()

	si.Set(1, "value")
	si.Set(2, 123)
	si.Set(3, true)
	si.Set(4, 1.23)

	val, ok := si.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "value", val)

	val, ok = si.Get(2)
	assert.True(t, ok)
	assert.Equal(t, 123, val)

	val, ok = si.Get(3)
	assert.True(t, ok)
	assert.Equal(t, true, val)

	val, ok = si.Get(4)
	assert.True(t, ok)
	assert.InDelta(t, 1.23, val, 0.0001)

	indexes := si.Indexes()
	slices.Sort(indexes)

	ir, err := json.Marshal(indexes)
	require.NoError(t, err)

	expected := []uint64{1, 2, 3, 4}
	ie, err := json.Marshal(expected)
	require.NoError(t, err)

	assert.JSONEq(t, string(ie), string(ir))

	si.Unset(2)

	indexes = si.Indexes()
	slices.Sort(indexes)

	ir, err = json.Marshal(indexes)
	require.NoError(t, err)

	expected = []uint64{1, 3, 4}
	ie, err = json.Marshal(expected)
	require.NoError(t, err)

	assert.JSONEq(t, string(ie), string(ir))

	assert.True(t, si.HasIndex(3))
	assert.False(t, si.HasIndex(2))
}

func BenchmarkSafeIndex_Set(b *testing.B) {
	m := NewSafeIndex()

	for i := range b.N {
		m.Set(uint64(i), i)
	}
}

func BenchmarkSafeIndex_Get(b *testing.B) {
	m := NewSafeIndex()
	for i := range 1000 {
		m.Set(uint64(i), i)
	}

	b.ResetTimer()

	for range b.N {
		_, _ = m.Get(500)
	}
}

func BenchmarkSafeIndex_Indexes(b *testing.B) {
	m := NewSafeIndex()
	for i := range 1000 {
		m.Set(uint64(i), i)
	}

	b.ResetTimer()

	for range b.N {
		_ = m.Indexes()
	}
}

func BenchmarkSafeIndex_Unset(b *testing.B) {
	m := NewSafeIndex()
	for i := range b.N {
		m.Set(uint64(i), i)
		m.Unset(uint64(i))
	}
}

func TestSafeIndex_ConcurrentAccess(t *testing.T) {
	si := NewSafeIndex()

	var wg sync.WaitGroup

	// number of concurrent goroutines
	const (
		goroutines       = 100
		keysPerGoroutine = 100
	)

	// run concurrent writers

	for g := range goroutines {
		wg.Add(1)

		go func(offset uint64) {
			defer wg.Done()

			for i := range uint64(keysPerGoroutine) {
				idx := offset*keysPerGoroutine + i
				si.Set(idx, fmt.Sprintf("val-%d", idx))
			}
		}(uint64(g))
	}

	wg.Wait()

	// verify all expected values exist
	expectedCount := goroutines * keysPerGoroutine
	indexes := si.Indexes()
	require.Len(t, indexes, expectedCount)

	for _, idx := range indexes {
		val, ok := si.Get(idx)
		require.True(t, ok)
		require.NotNil(t, val)
	}
}
