package huma

// CorrelationIDProviderFunc correlation ID provider.
type CorrelationIDProviderFunc func() string

// String method to get correlation ID.
func (f CorrelationIDProviderFunc) String() string {
	return f()
}
