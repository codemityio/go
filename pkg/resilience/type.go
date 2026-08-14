package resilience

// Option to be used to build a Retrier with optional config.
type Option func(retrier *Retrier)
