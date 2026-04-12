package huma

import (
	"html/template"
)

// WithConfig configuration option.
func WithConfig(config *Config) HandlerOption {
	return func(handler *Handler) {
		handler.config = config
	}
}

// WithRouter configuration option.
func WithRouter(router Router) HandlerOption {
	return func(handler *Handler) {
		handler.router = router
	}
}

// WithOperationWithHandler configuration option.
func WithOperationWithHandler(operations []OperationWithHandlerRegisterer) HandlerOption {
	return func(handler *Handler) {
		handler.operations = operations
	}
}

// WithMiddleware configuration option.
func WithMiddleware(middleware []MiddlewareFunc) HandlerOption {
	return func(handler *Handler) {
		handler.middleware = middleware
	}
}

// WithFormatter configuration option.
func WithFormatter(formatter []Formatter) HandlerOption {
	return func(handler *Handler) {
		handler.formatter = formatter
	}
}

// WithSecuritySchemes configuration option.
func WithSecuritySchemes(securityScheme []SecurityScheme) HandlerOption {
	return func(handler *Handler) {
		handler.securityScheme = securityScheme
	}
}

// WithAccessLogMiddlewareLogFunc configuration option.
func WithAccessLogMiddlewareLogFunc(fn LogFunc) AccessLogMiddlewareOption {
	return func(m *AccessLogMiddleware) {
		m.logFunc = fn
	}
}

// WithCorrelationIDMiddlewareCorrelationIDProvider configuration option.
func WithCorrelationIDMiddlewareCorrelationIDProvider(
	provider CorrelationIDProvider,
) CorrelationIDMiddlewareOption {
	return func(m *CorrelationIDMiddleware) {
		m.provider = provider
	}
}

// WithFormatTextHTMLTemplate configuration option.
func WithFormatTextHTMLTemplate(tpl string) FormatTextHTMLOption {
	return func(formatter *FormatTextHTML) {
		formatter.template = template.Must(
			template.New("text-html").Parse(tpl),
		)
	}
}
