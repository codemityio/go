package sync

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeMap(t *testing.T) {
	sf := NewSafeMap()

	sf.Set("key_string", "value")
	sf.Set("key_int", 1)
	sf.Set("key_bool", true)
	sf.Set("key_float", 1.23)

	val, ok := sf.Get("key_string")

	assert.True(t, ok)
	assert.Equal(t, "value", val)

	val, ok = sf.Get("key_int")

	assert.True(t, ok)
	assert.Equal(t, 1, val)

	val, ok = sf.Get("key_bool")

	assert.True(t, ok)
	assert.Equal(t, true, val)

	val, ok = sf.Get("key_float")

	assert.True(t, ok)
	assert.InDelta(t, 1.23, val, 0.0001)

	keys := sf.Keys()
	sort.Strings(keys)

	kr, err := json.Marshal(keys)
	require.NoError(t, err)

	ke, err := json.Marshal([]string{
		"key_bool",
		"key_float",
		"key_int",
		"key_string",
	})
	require.NoError(t, err)

	assert.JSONEq(t, string(ke), string(kr))

	sf.Unset("key_int")

	keys = sf.Keys()
	sort.Strings(keys)

	kr, err = json.Marshal(keys)
	require.NoError(t, err)

	ke, err = json.Marshal([]string{
		"key_bool",
		"key_float",
		"key_string",
	})
	require.NoError(t, err)

	assert.JSONEq(t, string(ke), string(kr))

	assert.True(t, sf.HasKey("key_bool"))
	assert.False(t, sf.HasKey("key_int"))
}

func BenchmarkSafeMap_Set(b *testing.B) {
	m := NewSafeMap()

	for i := range b.N {
		m.Set(fmt.Sprintf("key_%d", i), i)
	}
}

func BenchmarkSafeMap_Get(b *testing.B) {
	m := NewSafeMap()
	for i := range 1000 {
		m.Set(fmt.Sprintf("key_%d", i), i)
	}

	b.ResetTimer()

	for range b.N {
		_, _ = m.Get("key_500") // fixed key to avoid cache miss effects
	}
}

func BenchmarkSafeMap_Keys(b *testing.B) {
	m := NewSafeMap()
	for i := range 1000 {
		m.Set(fmt.Sprintf("key_%d", i), i)
	}

	b.ResetTimer()

	for range b.N {
		_ = m.Keys()
	}
}

func BenchmarkSafeMap_Unset(b *testing.B) {
	m := NewSafeMap()

	for i := range b.N {
		key := fmt.Sprintf("key_%d", i)
		m.Set(key, i)
		m.Unset(key)
	}
}
