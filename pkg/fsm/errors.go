package fsm

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg a Broker package error.
	ErrPkg = errors.New("fsm")
	// ErrUnmarshal error.
	ErrUnmarshal = fmt.Errorf("%w: unmarshal", ErrPkg)
)
