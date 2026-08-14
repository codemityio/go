package transform

import (
	"encoding/json"
	"errors"
	"fmt"
)

// JSON transform type.
type JSON struct{}

// IsMarshal reports whether err is or wraps the error Transform returns when marshalling fails.
func (j *JSON) IsMarshal(err error) bool {
	return errors.Is(err, ErrMarshal)
}

// IsUnmarshal reports whether err is or wraps the error Transform returns when unmarshalling fails.
func (j *JSON) IsUnmarshal(err error) bool {
	return errors.Is(err, ErrUnmarshal)
}

// Transform method to transform types by using json marshal/unmarshal functions.
func (j *JSON) Transform(input any, output any) error {
	serial, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrMarshal, err)
	}

	if e := json.Unmarshal(serial, output); e != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshal, e)
	}

	return nil
}
