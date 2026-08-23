package json

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerialiser_Serialise(t *testing.T) {
	type data struct {
		One   string `groups:"one"     json:"one"`
		Two   string `groups:"two"     json:"two"`
		Three string `groups:"one,two" json:"three"`
	}

	tests := map[string]struct {
		data           data
		groups         []string
		expectedResult string
		expectedError  error
	}{
		"success-one": {
			data: data{
				One:   "one",
				Two:   "two",
				Three: "three",
			},
			groups:         []string{"one"},
			expectedResult: `{"one":"one","three":"three"}`,
			expectedError:  nil,
		},
		"success-two": {
			data: data{
				One:   "one",
				Two:   "two",
				Three: "three",
			},
			groups:         []string{"two"},
			expectedResult: `{"two":"two","three":"three"}`,
			expectedError:  nil,
		},
		"success-three": {
			data: data{
				One:   "one",
				Two:   "two",
				Three: "three",
			},
			groups:         []string{"one", "two"},
			expectedResult: `{"one":"one","two":"two","three":"three"}`,
			expectedError:  nil,
		},
		"success-four": {
			data: data{
				One:   "one",
				Two:   "two",
				Three: "three",
			},
			groups:         []string{"four"},
			expectedResult: `{}`,
			expectedError:  nil,
		},
	}

	for i, test := range tests {
		t.Run(i, func(t *testing.T) {
			s := NewSerialiser()

			res, err := s.Serialise(test.data, test.groups)

			require.ErrorIs(t, err, test.expectedError)

			assert.JSONEq(t, test.expectedResult, string(res))
		})
	}
}

func TestSerialiser_IsMarshal(t *testing.T) {
	errBoom := errors.New("error")

	tests := map[string]struct {
		err  error
		want bool
	}{
		"true for the bare sentinel": {
			err:  ErrMarshal,
			want: true,
		},
		"true when wrapped": {
			err:  fmt.Errorf("%w: %w", ErrMarshal, errBoom),
			want: true,
		},
		"false for an unrelated error": {
			err:  errBoom,
			want: false,
		},
		"false for nil": {
			err:  nil,
			want: false,
		},
	}

	s := NewSerialiser()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, test.want, s.IsMarshal(test.err))
		})
	}
}
