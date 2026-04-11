package transform

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg package error.
	ErrPkg = errors.New("transform")
	// ErrMarshal error.
	ErrMarshal = fmt.Errorf("%w: unable to marshal", ErrPkg)
	// ErrUnmarshal error.
	ErrUnmarshal = fmt.Errorf("%w: unable to unmarshal", ErrPkg)
)
