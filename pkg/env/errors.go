package env

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg package error.
	ErrPkg = errors.New("env")
	// ErrEnv error.
	ErrEnv = fmt.Errorf("%w: env error", ErrPkg)
)
