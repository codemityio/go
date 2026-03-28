package json

import "encoding/json"

// RawMessageSequenceMap provides deterministic, per-key sequences of json.RawMessage values.
// It is useful for testing scenarios where each key should return a specific sequence
// of values on consecutive calls.
type RawMessageSequenceMap struct {
	data   map[string][]json.RawMessage // maps each key to a sequence of messages
	cursor map[string]int               // tracks the next index for each key
}

// Next returns the next json.RawMessage for the given key in the predefined sequence.
// If the key is unknown or has an empty sequence, it returns nil.
// If the sequence is exhausted, the last value is returned repeatedly.
func (s *RawMessageSequenceMap) Next(key string) json.RawMessage {
	list, ok := s.data[key]
	if !ok || len(list) == 0 {
		return nil
	}

	i := s.cursor[key]

	if i >= len(list) {
		return list[len(list)-1]
	}

	s.cursor[key] = i + 1

	return list[i]
}
