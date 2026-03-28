package json

import (
	"encoding/json"
	"fmt"
)

// MustBuildJSON factory, use it if the error is not worth handling.
func MustBuildJSON[T any](data json.RawMessage) *T {
	config, err := NewJSON[T](data)
	if err != nil {
		panic(err)
	}

	return config
}

// NewJSON factory.
func NewJSON[T any](data json.RawMessage) (*T, error) {
	var config T

	if e := json.Unmarshal(data, &config); e != nil {
		return nil, fmt.Errorf("%w: %w", ErrJSON, e)
	}

	return &config, nil
}
