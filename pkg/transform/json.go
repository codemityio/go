package transform

import (
	"encoding/json"
	"fmt"
)

// JSON transform type.
type JSON struct{}

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
