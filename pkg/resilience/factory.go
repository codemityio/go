package resilience

// NewRetrier creates a new Retrier, applying options on top of the defaults: BaseDelayDefault,
// MaxDelayDefault and RetriesDefault.
func NewRetrier(options ...Option) *Retrier {
	retrier := &Retrier{
		baseDelay: BaseDelayDefault,
		maxDelay:  MaxDelayDefault,
		retries:   RetriesDefault,
	}

	for _, opt := range options {
		opt(retrier)
	}

	return retrier
}
