package huma

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg error.
	ErrPkg = errors.New("huma")
	// ErrMarshal error.
	ErrMarshal = fmt.Errorf("%w: marshal", ErrPkg)
	// ErrUnmarshal error.
	ErrUnmarshal = fmt.Errorf("%w: unmarshal", ErrPkg)
)
