package json

import (
	"encoding/json"
	"fmt"

	"github.com/liip/sheriff/v2"
)

// Serialiser is a serialiser service.
type Serialiser struct{}

// NewSerialiser a group specific factory.
func NewSerialiser() *Serialiser {
	return &Serialiser{}
}

// Serialise a simple serialiser with capability to serialise for specific group of interest.
func (s *Serialiser) Serialise(input any, groups []string) ([]byte, error) {
	var err error

	data := input

	var out []byte

	if len(groups) > 0 {
		opts := &sheriff.Options{ //nolint:exhaustruct
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
