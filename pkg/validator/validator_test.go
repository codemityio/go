package validator

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/schema.json
	schema []byte

	//go:embed testdata/success-input.json
	successInputJSON []byte

	//go:embed testdata/error-input.json
	errorInputJSON []byte

	//go:embed testdata/success-input.yaml
	successInputYAML []byte

	//go:embed testdata/error-input.yaml
	errorInputYAML []byte
)

func TestDefaultValidator_ValidateJSON(t *testing.T) {
	validator := New()

	for name, test := range map[string]struct {
		input          []byte
		schema         []byte
		expectedError  error
		expectedResult []Error
	}{
		"success": {
			input:          successInputJSON,
			schema:         schema,
			expectedError:  nil,
			expectedResult: nil,
		},
		"error": {
			input:         errorInputJSON,
			schema:        schema,
			expectedError: nil,
			expectedResult: []Error{
				{
					Field:   "(root)",
					Message: "Invalid type. Expected: object, given: string",
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			result, err := validator.ValidateJSON(test.schema, test.input)
			require.ErrorIs(t, err, test.expectedError)

			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func TestDefaultValidator_ValidateYAML(t *testing.T) {
	validator := New()

	for name, test := range map[string]struct {
		input          []byte
		schema         []byte
		expectedError  error
		expectedResult []Error
	}{
		"success": {
			input:          successInputYAML,
			schema:         schema,
			expectedError:  nil,
			expectedResult: nil,
		},
		"error": {
			input:         errorInputYAML,
			schema:        schema,
			expectedError: nil,
			expectedResult: []Error{
				{
					Field:   "(root)",
					Message: "Invalid type. Expected: object, given: string",
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			result, err := validator.ValidateYAML(test.schema, test.input)
			require.ErrorIs(t, err, test.expectedError)

			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func TestDefaultValidator_ValidateAny(t *testing.T) {
	validator := New()

	for name, test := range map[string]struct {
		input          any
		schema         []byte
		expectedError  error
		expectedResult []Error
	}{
		"success": {
			input: map[string]string{
				"message": "success",
			},
			schema:         schema,
			expectedError:  nil,
			expectedResult: nil,
		},
		"error": {
			input:         "error",
			schema:        schema,
			expectedError: nil,
			expectedResult: []Error{
				{
					Field:   "(root)",
					Message: "Invalid type. Expected: object, given: string",
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			result, err := validator.ValidateAny(test.schema, test.input)
			require.ErrorIs(t, err, test.expectedError)

			assert.Equal(t, test.expectedResult, result)
		})
	}
}
