package json

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONSequence_Next(t *testing.T) {
	seq := NewRawMessageSequenceMap(map[string][]json.RawMessage{
		"one":   {json.RawMessage(`"first"`)},
		"two":   {json.RawMessage(`"a"`), json.RawMessage(`"b"`)},
		"empty": {},
	})

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{name: "single call returns first", key: "one", expected: `"first"`},
		{name: "first of two", key: "two", expected: `"a"`},
		{name: "second of two", key: "two", expected: `"b"`},
		{name: "repeat last when exhausted", key: "two", expected: `"b"`},
		{name: "unknown key returns nil", key: "unknown", expected: ``},
		{name: "empty list returns nil", key: "empty", expected: ``},
		{name: "repeating last value of single-element list", key: "one", expected: `"first"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seq.Next(tt.key)

			if tt.expected == "" {
				require.Nil(t, got, "expected nil for key %q", tt.key)
			} else {
				require.Equal(t, tt.expected, string(got), "unexpected value for key %q", tt.key)
			}
		})
	}
}
