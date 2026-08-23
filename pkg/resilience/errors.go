package resilience

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg package error.
	ErrPkg = errors.New("resilience")
	// ErrRetryWaitCanceled is returned when ctx is done before the backoff for the next retry attempt elapses.
	ErrRetryWaitCanceled = fmt.Errorf("%w: retry wait canceled", ErrPkg)
	// ErrRetriesExhausted is returned when a Retrier's configured attempts are used up without the handler succeeding.
	ErrRetriesExhausted = fmt.Errorf("%w: retries exhausted", ErrPkg)
)
