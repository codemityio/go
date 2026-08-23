package resilience

import "time"

const (
	// BaseDelayDefault default wait before the second attempt.
	BaseDelayDefault = 100 * time.Millisecond
	// MaxDelayDefault default cap on the backoff between attempts.
	MaxDelayDefault = 30 * time.Second
	// RetriesDefault default maximum number of attempts, including the first.
	RetriesDefault = 3
)
