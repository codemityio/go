package transform

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type typeA struct {
	FieldOne string `json:"fieldOne"`
}

type typeB struct {
	FieldOne string `json:"fieldOne"`
}

func TestSecurity_CreateUser(t *testing.T) {
	tests := map[string]struct {
		input         any
		output        any
		expected      any
		expectedError error
	}{
		"success": {
			input:         &typeA{FieldOne: "fieldOne"},
			output:        &typeB{},
			expected:      &typeB{FieldOne: "fieldOne"},
			expectedError: nil,
		},
		"failure-marshal": {
			input:         make(chan int),
			output:        &typeB{},
			expected:      &typeB{},
			expectedError: ErrMarshal,
		},
		"failure": {
			input:         nil,
			output:        nil,
			expected:      nil,
			expectedError: ErrUnmarshal,
		},
	}

	for i, test := range tests {
		t.Run(i, func(t *testing.T) {
			transform := NewJSON()

			require.ErrorIs(t, transform.Transform(test.input, test.output), test.expectedError)

			require.Equal(t, test.expected, test.output)
		})
	}
}

func TestJSON_IsMarshalIsUnmarshal(t *testing.T) {
	errBoom := errors.New("error")
	transform := NewJSON()

	require.True(t, transform.IsMarshal(ErrMarshal))
	require.True(t, transform.IsMarshal(fmt.Errorf("%w: %w", ErrMarshal, errBoom)))
	require.False(t, transform.IsMarshal(errBoom))
	require.False(t, transform.IsMarshal(ErrUnmarshal))
	require.False(t, transform.IsMarshal(nil))

	require.True(t, transform.IsUnmarshal(ErrUnmarshal))
	require.True(t, transform.IsUnmarshal(fmt.Errorf("%w: %w", ErrUnmarshal, errBoom)))
	require.False(t, transform.IsUnmarshal(errBoom))
	require.False(t, transform.IsUnmarshal(ErrMarshal))
	require.False(t, transform.IsUnmarshal(nil))
}
