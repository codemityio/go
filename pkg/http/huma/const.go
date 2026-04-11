package huma

const (
	// SchemeHTTP scheme name.
	SchemeHTTP = "http"
	// SchemeHTTPS scheme name.
	SchemeHTTPS = "https"

	// HeaderCorrelationID correlation ID header name.
	HeaderCorrelationID = "X-Correlation-ID"
	// HeaderForwardedProto header name.
	HeaderForwardedProto = "X-Forwarded-Proto"
	// HeaderForwardedScheme header name.
	HeaderForwardedScheme = "X-Forwarded-Scheme"

	// ContextCorrelationID context correlation ID name.
	ContextCorrelationID = "correlation-id"

	ContentTypeTextPlain       = "text/plain"
	ContentTypeTextHTML        = "text/html"
	ContentTypeApplicationJSON = "application/json"
)
