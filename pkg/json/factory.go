package json

import "encoding/json"

// NewRawMessageSequenceMap initialises and returns a new instance of JSONSequence.
func NewRawMessageSequenceMap(input map[string][]json.RawMessage) *RawMessageSequenceMap {
	return &RawMessageSequenceMap{
		data:   input,
		cursor: make(map[string]int),
	}
}
