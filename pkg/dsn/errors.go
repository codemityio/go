package dsn

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg package error.
	ErrPkg = errors.New("dsn")
	// ErrInvalidDSN error.
	ErrInvalidDSN = fmt.Errorf("%w: invalid DSN", ErrPkg)
	// ErrInvalidQuery error.
	ErrInvalidQuery = fmt.Errorf("%w: invalid query", ErrPkg)
	// ErrInvalidPort error.
	ErrInvalidPort = fmt.Errorf("%w: invalid port", ErrPkg)
)
