package sync

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg package error.
	ErrPkg = errors.New("auxiliary")
	// ErrSafeBufferWrite error.
	ErrSafeBufferWrite = fmt.Errorf("%w: safe buffer write", ErrPkg)
)
