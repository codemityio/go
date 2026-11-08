package validator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ Validator[[]byte]            = (*GenericValidator[[]byte])(nil)
	_ Validator[map[string]string] = (*GenericValidator[map[string]string])(nil)
)

type message struct {
	Message string `json:"message"`
}

func TestGenericValidator_Validate(t *testing.T) {
	rootTypeError := []Error{
		{
			Field:   "(root)",
			Message: "Invalid type. Expected: object, given: string",
		},
	}

	for name, test := range map[string]struct {
		validate       func() ([]Error, error)
		expectedError  error
		expectedResult []Error
	}{
		"bytes success": {
			validate: func() ([]Error, error) {
				return NewGeneric[[]byte]().Validate(schema, successInputJSON)
			},
			expectedError:  nil,
			expectedResult: nil,
		},
		"bytes error": {
			validate: func() ([]Error, error) {
				return NewGeneric[[]byte]().Validate(schema, errorInputJSON)
			},
			expectedError:  nil,
			expectedResult: rootTypeError,
		},
		"bytes invalid schema": {
			validate: func() ([]Error, error) {
				return NewGeneric[[]byte]().Validate([]byte("{"), successInputJSON)
			},
			expectedError:  ErrPkg,
			expectedResult: nil,
		},
		"raw message success": {
			validate: func() ([]Error, error) {
				return NewGeneric[json.RawMessage]().Validate(schema, successInputJSON)
			},
			expectedError:  nil,
			expectedResult: nil,
		},
		"raw message error": {
			validate: func() ([]Error, error) {
				return NewGeneric[json.RawMessage]().Validate(schema, errorInputJSON)
			},
			expectedError:  nil,
			expectedResult: rootTypeError,
		},
		"struct success": {
			validate: func() ([]Error, error) {
				return NewGeneric[message]().Validate(schema, message{Message: "success"})
			},
			expectedError:  nil,
			expectedResult: nil,
		},
		"map success": {
			validate: func() ([]Error, error) {
				return NewGeneric[map[string]string]().Validate(schema, map[string]string{"message": "success"})
			},
			expectedError:  nil,
			expectedResult: nil,
		},
		"string error": {
			validate: func() ([]Error, error) {
				return NewGeneric[string]().Validate(schema, "error")
			},
			expectedError:  nil,
			expectedResult: rootTypeError,
		},
	} {
		t.Run(name, func(t *testing.T) {
			result, err := test.validate()
			require.ErrorIs(t, err, test.expectedError)

			assert.Equal(t, test.expectedResult, result)
		})
	}
}
