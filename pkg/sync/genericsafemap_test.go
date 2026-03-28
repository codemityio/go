package sync

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericSafeMap_String(t *testing.T) {
	sf := NewGenericSafeMap[string]()

	sf.Set("key_string", "value")

	v, s := sf.Get("key_string")

	assert.True(t, s)
	assert.Equal(t, "value", v)

	keys := sf.Keys()
	sort.Strings(keys)

	kr, err := json.Marshal(keys)
	require.NoError(t, err)

	ke, err := json.Marshal([]string{
		"key_string",
	})
	require.NoError(t, err)

	assert.JSONEq(t, string(ke), string(kr))

	sf.Unset("key_string")

	keys = sf.Keys()
	sort.Strings(keys)

	kr, err = json.Marshal(keys)
	require.NoError(t, err)

	ke, err = json.Marshal([]string{})
	require.NoError(t, err)

	assert.JSONEq(t, string(ke), string(kr))
}

func TestGenericSafeMap_Int(t *testing.T) {
	sf := NewGenericSafeMap[int]()

	sf.Set("key_int", 1)

	v, s := sf.Get("key_int")

	assert.True(t, s)
	assert.Equal(t, 1, v)

	ke, err := json.Marshal([]string{
		"key_int",
	})
	require.NoError(t, err)

	keys := sf.Keys()
	sort.Strings(keys)

	kr, err := json.Marshal(keys)
	require.NoError(t, err)

	assert.JSONEq(t, string(ke), string(kr))

	sf.Unset("key_int")

	keys = sf.Keys()
	sort.Strings(keys)

	kr, err = json.Marshal(keys)
	require.NoError(t, err)

	ke, err = json.Marshal([]string{})
	require.NoError(t, err)

	assert.JSONEq(t, string(ke), string(kr))
}
