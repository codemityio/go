package huma

import (
	"net"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
)

// AccessLogMiddleware middleware to log access events.
type AccessLogMiddleware struct {
	logFunc LogFunc
}

// Handler middleware handler.
func (m *AccessLogMiddleware) Handler(ctx huma.Context, _ huma.API, next func(huma.Context)) {
	url := ctx.URL()

	scheme := m.deriveScheme(ctx)

	host, port := m.deriveHostAndPort(ctx, scheme)

	var id string

	if value, ok := ctx.Context().Value(ContextCorrelationID).(string); ok {
		id = value
	}

	if m.logFunc != nil {
		m.logFunc(url.String(), []LogField{
			{Key: "method", Value: ctx.Method()},
			{Key: "scheme", Value: scheme},
			{Key: "host", Value: host},
			{Key: "port", Value: port},
			{Key: "path", Value: url.Path},
			{Key: "correlation-id", Value: id},
		}...)
	}

	next(ctx)
}

func (m *AccessLogMiddleware) deriveScheme(ctx huma.Context) string {
	scheme := ctx.Header(HeaderForwardedProto)

	if scheme == "" {
		scheme = ctx.Header(HeaderForwardedScheme)
	}

	if scheme == "" {
		if ctx.TLS() != nil {
			scheme = SchemeHTTPS
		} else {
			scheme = SchemeHTTP
		}
	}

	return scheme
}

func (m *AccessLogMiddleware) deriveHostAndPort(ctx huma.Context, scheme string) (string, int) {
	hostHeader := ctx.Host() // e.g. "localhost:8080" or "example.com"

	var (
		host, prt string
		port      int
		err       error
	)

	if host, prt, err = net.SplitHostPort(hostHeader); err == nil {
		port, err = strconv.Atoi(prt)
		if err != nil {
			port = 0
		}
	} else {
		// no explicit port provided, assume default
		host = hostHeader

		if scheme == SchemeHTTPS {
			port = 443
		} else {
			port = 80
		}
	}

	return host, port
}

// CorrelationIDMiddleware middleware to handle correlation ID.
type CorrelationIDMiddleware struct {
	provider CorrelationIDProvider
}

// Handler middleware handler.
func (m *CorrelationIDMiddleware) Handler(ctx huma.Context, _ huma.API, next func(huma.Context)) {
	id := ctx.Header(HeaderCorrelationID)

	if id == "" {
		id = m.provider.String()
	}

	ctx = huma.WithValue(ctx, ContextCorrelationID, id)

	next(ctx)
}
