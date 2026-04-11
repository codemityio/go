package mongo

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg package error.
	ErrPkg = errors.New("mongodb")
	// ErrClient error.
	ErrClient = fmt.Errorf("%w: client error", ErrPkg)
)
