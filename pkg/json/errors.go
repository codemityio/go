package json

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg package error.
	ErrPkg = errors.New("json")
	// ErrJSON error.
	ErrJSON = fmt.Errorf("%w: JSON error", ErrPkg)
)
