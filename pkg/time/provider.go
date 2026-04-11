package time

import "time"

// ProviderFunc provides the current time.
// It can be replaced in tests to control time behaviour.
type ProviderFunc func() time.Time

// Now returns the current time using the provider.
// If the provider is nil, time.Now is used.
func (t ProviderFunc) Now() time.Time {
	if t == nil {
		return time.Now()
	}

	return t()
}

// Unix returns t as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
func (t ProviderFunc) Unix() int64 {
	if t == nil {
		return time.Now().Unix()
	}

	return t().Unix()
}
