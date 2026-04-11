package huma

import (
	"io"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
)

// Router common router interface.
type Router interface {
	chi.Router
}

// OperationWithHandlerRegisterer an operation with handler registerer.
type OperationWithHandlerRegisterer interface {
	Register(api huma.API)
}

// MiddlewareHandler interface.
type MiddlewareHandler interface {
	Handler(ctx huma.Context, api huma.API, next func(huma.Context))
}

// CorrelationIDProvider correlation ID provider.
type CorrelationIDProvider interface {
	String() string
}

// Formatter input/output formatter.
type Formatter interface {
	ContentType() string
	Marshal(writer io.Writer, v any) error
	Unmarshal(data []byte, v any) error
}

// SecurityScheme interface.
type SecurityScheme interface {
	Key() string
	Type() string
	Name() string
	Description() string
	Scheme() string
	In() string
}
