package json

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/liip/sheriff/v2"
)

// Serialiser is a serialiser service.
type Serialiser struct{}

// NewSerialiser a group specific factory.
func NewSerialiser() *Serialiser {
	return &Serialiser{}
}

// IsMarshal reports whether err is or wraps the error Serialise returns when marshalling fails.
func (s *Serialiser) IsMarshal(err error) bool {
	return errors.Is(err, ErrMarshal)
}

// Serialise a simple serialiser with capability to serialise for specific group of interest.
func (s *Serialiser) Serialise(input any, groups []string) ([]byte, error) {
	var err error

	data := input

	var out []byte

	if len(groups) > 0 {
		opts := &sheriff.Options{ //nolint:exhaustruct_v5
			ApiVersion:      nil,
			Groups:          groups,
			IncludeEmptyTag: false,
		}

		data, err = sheriff.Marshal(opts, input)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrMarshal, err)
		}
	}

	out, err = json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrMarshal, err)
	}

	return out, nil
}
