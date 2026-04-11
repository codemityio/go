package huma

import (
	"html/template"
	"maps"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// NewHandler a factory function to create a new server.
func NewHandler(options ...HandlerOption) *Handler {
	handler := &Handler{
		config: &Config{
			Title:         "API",
			Version:       "v0.0.0",
			OpenAPIPath:   "/openapi",
			DocsPath:      "/docs",
			SchemasPath:   "/schemas",
			DefaultFormat: "application/json",
		},
		router:         nil,
		operations:     nil,
		middleware:     nil,
		formatter:      nil,
		securityScheme: nil,
	}

	for _, option := range options {
		option(handler)
	}

	// Chi router documentation https://github.com/go-chi/chi.
	if handler.router == nil {
		handler.router = chi.NewRouter()
	}

	config := huma.DefaultConfig(handler.config.Title, handler.config.Version)

	config.OpenAPIPath = handler.config.OpenAPIPath
	config.DocsPath = handler.config.DocsPath
	config.SchemasPath = handler.config.SchemasPath
	config.DefaultFormat = handler.config.DefaultFormat

	// clone formats to avoid shared map writes
	origFormats := config.Formats
	config.Formats = make(map[string]huma.Format, len(origFormats)+len(handler.formatter))

	maps.Copy(config.Formats, origFormats)

	// register formatters
	for i := range handler.formatter {
		config.Formats[handler.formatter[i].ContentType()] = huma.Format{
			Marshal:   handler.formatter[i].Marshal,
			Unmarshal: handler.formatter[i].Unmarshal,
		}
	}

	if handler.securityScheme != nil {
		config.Components.SecuritySchemes = make(map[string]*huma.SecurityScheme)
	}

	// register security schemes
	for i := range handler.securityScheme {
		config.Components.SecuritySchemes[handler.securityScheme[i].Key()] = &huma.SecurityScheme{
			Type:             handler.securityScheme[i].Type(),
			Description:      handler.securityScheme[i].Description(),
			Name:             handler.securityScheme[i].Name(),
			In:               handler.securityScheme[i].In(),
			Scheme:           handler.securityScheme[i].Scheme(),
			BearerFormat:     "",
			Flows:            nil,
			OpenIDConnectURL: "",
			Extensions:       nil,
		}
	}

	api := humachi.New(handler.router, config)

	// register middleware
	midFuncs := make([]func(ctx huma.Context, next func(huma.Context)), 0, len(handler.middleware))

	for _, mw := range handler.middleware {
		midFuncs = append(midFuncs, func(ctx huma.Context, next func(huma.Context)) {
			mw(ctx, api, next)
		})
	}

	api.UseMiddleware(midFuncs...)

	for _, v := range handler.operations {
		v.Register(api)
	}

	return handler
}

// NewAccessLogMiddleware factory function to crate a new access log middleware instance.
func NewAccessLogMiddleware(options ...AccessLogMiddlewareOption) *AccessLogMiddleware {
	middleware := &AccessLogMiddleware{
		logFunc: nil,
	}

	for _, option := range options {
		option(middleware)
	}

	return middleware
}

// NewCorrelationIDMiddleware factory function to crate a correlation ID handler middleware instance.
func NewCorrelationIDMiddleware(options ...CorrelationIDMiddlewareOption) *CorrelationIDMiddleware {
	middleware := &CorrelationIDMiddleware{
		provider: nil,
	}

	for _, option := range options {
		option(middleware)
	}

	return middleware
}

// NewUUIDCorrelationIDProvider factory function.
func NewUUIDCorrelationIDProvider() CorrelationIDProviderFunc {
	return uuid.NewString
}

// NewFormatTextPlain factory function.
func NewFormatTextPlain() *FormatTextPlain {
	return &FormatTextPlain{}
}

// NewFormatTextHTML factory function.
func NewFormatTextHTML(options ...FormatTextHTMLOption) *FormatTextHTML {
	formatter := &FormatTextHTML{
		template: template.Must(
			template.New("text-html").Parse(`<!doctype html>
<html>
  <body>
    <pre>{{ . }}</pre>
  </body>
</html>
`),
		),
	}

	for _, option := range options {
		option(formatter)
	}

	return formatter
}

// NewFormatApplicationJSON factory function.
func NewFormatApplicationJSON() *FormatApplicationJSON {
	return &FormatApplicationJSON{}
}

// NewSecuritySchemeBasicAuth creates a new basic auth security scheme.
func NewSecuritySchemeBasicAuth() *SecuritySchemeBasicAuth {
	return &SecuritySchemeBasicAuth{}
}

// NewSecuritySchemeAPIHeaderKeyAuth creates a new header-based API key security scheme.
func NewSecuritySchemeAPIHeaderKeyAuth() *SecuritySchemeAPIHeaderKeyAuth {
	return &SecuritySchemeAPIHeaderKeyAuth{}
}

// NewSecuritySchemeAPIQueryKeyAuth creates a new query-based API key security scheme.
func NewSecuritySchemeAPIQueryKeyAuth() *SecuritySchemeAPIQueryKeyAuth {
	return &SecuritySchemeAPIQueryKeyAuth{}
}

// NewSecuritySchemeTokenAuth creates a new token-based authentication scheme.
func NewSecuritySchemeTokenAuth() *SecuritySchemeTokenAuth {
	return &SecuritySchemeTokenAuth{}
}
