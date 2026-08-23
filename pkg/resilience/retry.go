package resilience

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// HandlerFunc is executed on each attempt made by a Retrier. It returns an error to
// trigger another attempt.
type HandlerFunc func(ctx context.Context) error

// Retrier executes a HandlerFunc, retrying with exponential backoff between attempts.
type Retrier struct {
	baseDelay time.Duration
	maxDelay  time.Duration
	retries   int
}

// IsRetriesExhausted reports whether err is or wraps the error Execute returns once its
// configured retries are used up without the handler succeeding.
func (r *Retrier) IsRetriesExhausted(err error) bool {
	return errors.Is(err, ErrRetriesExhausted)
}

// IsRetryWaitCanceled reports whether err is or wraps the error Execute returns when ctx is
// done before the backoff for the next retry attempt elapses.
func (r *Retrier) IsRetryWaitCanceled(err error) bool {
	return errors.Is(err, ErrRetryWaitCanceled)
}

// Execute runs handler, retrying with exponential backoff on error until it succeeds, the
// configured retries are used up, or ctx is done first.
func (r *Retrier) Execute(ctx context.Context, handler HandlerFunc) error {
	var lastErr error

	for attempt := range r.retries {
		if err := r.waitBeforeRetry(ctx, attempt); err != nil {
			return err
		}

		lastErr = handler(ctx)
		if lastErr == nil {
			return nil
		}
	}

	if lastErr == nil {
		return ErrRetriesExhausted
	}

	return fmt.Errorf("%w: %w", ErrRetriesExhausted, lastErr)
}

// waitBeforeRetry waits out the backoff for attempt, or returns early if ctx is done first
// (e.g. the caller is approaching its own deadline).
func (r *Retrier) waitBeforeRetry(ctx context.Context, attempt int) error {
	wait := r.delay(attempt)
	if wait <= 0 {
		return nil
	}

	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %w", ErrRetryWaitCanceled, ctx.Err())
	case <-timer.C:
		return nil
	}
}

// delay returns the wait before retry attempt, doubling from baseDelay and capped at
// maxDelay. The first attempt (0) has no delay.
func (r *Retrier) delay(attempt int) time.Duration {
	if attempt <= 0 {
		return 0
	}

	backoff := r.baseDelay << attempt
	if backoff <= 0 || backoff > r.maxDelay {
		return r.maxDelay
	}

	return backoff
}
