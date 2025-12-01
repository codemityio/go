package validator

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// DefaultValidator is a Validator service.
type DefaultValidator struct{}

// ValidateJSON can be used to validate JSON value.
func (v *DefaultValidator) ValidateJSON(schema, input []byte) ([]Error, error) {
	res, err := gojsonschema.Validate(
		gojsonschema.NewBytesLoader(schema),
		gojsonschema.NewBytesLoader(input),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: validation failure: %w", ErrPkg, err)
	}

	return v.check(res), nil
}

// ValidateYAML can be used to validate JSON value.
func (v *DefaultValidator) ValidateYAML(schema, input []byte) ([]Error, error) {
	var value any

	if e := yaml.Unmarshal(input, &value); e != nil {
		return nil, fmt.Errorf("%w: validation failure: %w", ErrPkg, e)
	}

	input, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("%w: validation failure: %w", ErrPkg, err)
	}

	res, err := gojsonschema.Validate(
		gojsonschema.NewBytesLoader(schema),
		gojsonschema.NewBytesLoader(input),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: validation failure: %w", ErrPkg, err)
	}

	return v.check(res), nil
}

// ValidateAny can be used to validate any type.
func (v *DefaultValidator) ValidateAny(schema []byte, input any) ([]Error, error) {
	res, err := gojsonschema.Validate(
		gojsonschema.NewBytesLoader(schema),
		gojsonschema.NewGoLoader(input),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: validation failure: %w", ErrPkg, err)
	}

	return v.check(res), nil
}

func (v *DefaultValidator) check(result *gojsonschema.Result) []Error {
	var errors []Error

	if !result.Valid() {
		for _, d := range result.Errors() {
			errors = append(errors, Error{
				Field:   d.Field(),
				Message: d.Description(),
			})
		}
	}

	return errors
}
