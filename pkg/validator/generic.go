package validator

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// GenericValidator is a type-safe Validator implementation for values of type T.
//
// Byte-based inputs ([]byte and json.RawMessage) are treated as raw JSON documents, any other value is validated
// as a Go value (marshalled to JSON using its json tags).
type GenericValidator[T any] struct {
	base DefaultValidator
}

// Validate can be used to validate a value of type T.
func (v *GenericValidator[T]) Validate(schema []byte, input T) ([]Error, error) {
	var loader gojsonschema.JSONLoader

	switch value := any(input).(type) {
	case []byte:
		loader = gojsonschema.NewBytesLoader(value)
	case json.RawMessage:
		loader = gojsonschema.NewBytesLoader(value)
	default:
		loader = gojsonschema.NewGoLoader(input)
	}

	res, err := gojsonschema.Validate(gojsonschema.NewBytesLoader(schema), loader)
	if err != nil {
		return nil, fmt.Errorf("%w: validation failure: %w", ErrPkg, err)
	}

	return v.base.check(res), nil
}
