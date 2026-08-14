package resilience

import "time"

// WithBaseDelay sets the wait before the second attempt, doubling on each subsequent
// attempt up to maxDelay.
func WithBaseDelay(baseDelay time.Duration) Option {
	return func(retrier *Retrier) {
		retrier.baseDelay = baseDelay
	}
}

// WithMaxDelay sets the cap on the backoff between attempts.
func WithMaxDelay(maxDelay time.Duration) Option {
	return func(retrier *Retrier) {
		retrier.maxDelay = maxDelay
	}
}

// WithRetries sets the maximum number of attempts, including the first.
func WithRetries(retries int) Option {
	return func(retrier *Retrier) {
		retrier.retries = retries
	}
}
