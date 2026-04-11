package huma

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

// HandlerOption function.
type HandlerOption func(p *Handler)

// AccessLogMiddlewareOption function.
type AccessLogMiddlewareOption func(p *AccessLogMiddleware)

// CorrelationIDMiddlewareOption function.
type CorrelationIDMiddlewareOption func(p *CorrelationIDMiddleware)

// OperationWithHandler an operation type.
type OperationWithHandler[I, O any] struct {
	huma.Operation

	Handler func(context.Context, *I) (*O, error)
}

// Register to register operation in huma API.
func (op *OperationWithHandler[I, O]) Register(api huma.API) {
	huma.Register(api, op.Operation, op.Handler)
}

// MiddlewareFunc middleware function type.
type MiddlewareFunc func(ctx huma.Context, api huma.API, next func(huma.Context))

// HTMLErrorView type.
type HTMLErrorView struct {
	Status int
	Title  string
	Detail string
}

// FormatTextHTMLOption function.
type FormatTextHTMLOption func(p *FormatTextHTML)

// LogField a custom field struct representation.
type LogField struct {
	Key   string
	Value any
}

// LogFunc to allow log warnings and infos.
type LogFunc func(message string, fields ...LogField)
