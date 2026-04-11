//nolint:lll
package huma

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/codemityio/go/pkg/env"
	"github.com/codemityio/go/pkg/http/client"
	"github.com/codemityio/go/pkg/http/server"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// health.
type TestHandlerHealthInput struct{}

type TestHandlerHealthOutputBody struct {
	Message string `json:"message"`
}

type TestHandlerHealthOutput struct {
	Body TestHandlerHealthOutputBody `required:"true"`
}

// hello.
type TestHandlerHelloInput struct {
	Value string `doc:"Value" example:"World" json:"value" path:"value" required:"true"`
}

type TestHandlerHelloOutputBody struct {
	Message string `json:"message"`
}

type TestHandlerHelloOutput struct {
	Body TestHandlerHelloOutputBody `required:"true"`
}

// request header.
type TestHandlerRequestHeaderInput struct {
	RequestHeader string `header:"X-Request-Header"` //nolint:tagliatelle
}

type TestHandlerRequestHeaderOutputBody struct {
	Message string `json:"message"`
}

type TestHandlerRequestHeaderOutput struct {
	Body TestHandlerRequestHeaderOutputBody `required:"true"`
}

// request header with middleware and response header.
type TestHandlerRequestHeaderWithMiddlewareAndResponseHeaderInput struct {
	RequestHeader string `header:"X-Request-Header-With-Middleware-And-Response-Header"` //nolint:tagliatelle
}

type TestHandlerRequestHeaderWithMiddlewareAndResponseHeaderOutputBody struct {
	Message string `json:"message"`
}

type TestHandlerRequestHeaderWithMiddlewareAndResponseHeaderOutput struct {
	ResponseHeader string                                                            `header:"X-Response-Header"` //nolint:tagliatelle
	Body           TestHandlerRequestHeaderWithMiddlewareAndResponseHeaderOutputBody `                           required:"true"`
}

// request validation.
type TestHandlerRequestValidationInputBody struct {
	Message string `json:"message" minLength:"3" required:"true"`
}

type TestHandlerRequestValidationInput struct {
	Body TestHandlerRequestValidationInputBody `required:"true"`
}

type TestHandlerRequestValidationOutputBody struct {
	Message string `json:"message" required:"true"`
}

type TestHandlerRequestValidationOutput struct {
	Body TestHandlerRequestValidationOutputBody `required:"true"`
}

// response validation.
type TestHandlerResponseValidationInput struct{}

type TestHandlerResponseValidationOutputBody struct {
	Message string `json:"message" minLength:"3" required:"true"`
}

type TestHandlerResponseValidationOutput struct {
	Body TestHandlerResponseValidationOutputBody `required:"true"`
}

func TestHandler(t *testing.T) {
	t.Setenv("PORT", "9091")
	t.Setenv("READ_TIMEOUT", "0")
	t.Setenv("WRITE_TIMEOUT", "0")
	t.Setenv("TITLE", "API")
	t.Setenv("VERSION", "v0.0.0")
	t.Setenv("OPEN_API_PATH", "/openapi")
	t.Setenv("DOCS_PATH", "/docs")
	t.Setenv("SCHEMAS_PATH", "/schemas")
	t.Setenv("DEFAULT_FORMAT", "application/json")

	handlerConfig, err := env.New[Config]()
	require.NoError(t, err)

	serverConfig, err := env.New[server.Config]()
	require.NoError(t, err)

	operationsWithHandler := []OperationWithHandlerRegisterer{
		&OperationWithHandler[TestHandlerHealthInput, TestHandlerHealthOutput]{
			Operation: huma.Operation{
				OperationID: "get-health",
				Method:      http.MethodGet,
				Path:        "/health",
				Summary:     "Health endpoint",
			},
			Handler: func(ctx context.Context, input *TestHandlerHealthInput) (*TestHandlerHealthOutput, error) {
				resp := &TestHandlerHealthOutput{}
				resp.Body.Message = "I am healthy!"

				return resp, nil
			},
		},
		&OperationWithHandler[TestHandlerHelloInput, TestHandlerHelloOutput]{
			Operation: huma.Operation{
				OperationID: "get-hello",
				Method:      http.MethodGet,
				Path:        "/hello/{value}",
				Summary:     "Hello endpoint",
			},
			Handler: func(ctx context.Context, input *TestHandlerHelloInput) (*TestHandlerHelloOutput, error) {
				resp := &TestHandlerHelloOutput{}
				resp.Body.Message = fmt.Sprintf("Hello %s!", input.Value)

				return resp, nil
			},
		},
		&OperationWithHandler[TestHandlerRequestHeaderInput, TestHandlerRequestHeaderOutput]{
			Operation: huma.Operation{
				OperationID: "get-request-header",
				Method:      http.MethodGet,
				Path:        "/request-header",
				Summary:     "Request header endpoint",
			},
			Handler: func(ctx context.Context, input *TestHandlerRequestHeaderInput) (*TestHandlerRequestHeaderOutput, error) {
				resp := &TestHandlerRequestHeaderOutput{}
				resp.Body.Message = input.RequestHeader

				return resp, nil
			},
		},
		&OperationWithHandler[
			TestHandlerRequestHeaderWithMiddlewareAndResponseHeaderInput,
			TestHandlerRequestHeaderWithMiddlewareAndResponseHeaderOutput,
		]{
			Operation: huma.Operation{
				OperationID: "get-request-header-with-middleware-and-response-header",
				Method:      http.MethodGet,
				Path:        "/request-header-with-middleware-and-response-header",
				Summary:     "Request header with middleware and response header endpoint",
			},
			Handler: func(
				ctx context.Context,
				input *TestHandlerRequestHeaderWithMiddlewareAndResponseHeaderInput,
			) (*TestHandlerRequestHeaderWithMiddlewareAndResponseHeaderOutput, error) {
				resp := &TestHandlerRequestHeaderWithMiddlewareAndResponseHeaderOutput{}
				resp.Body.Message = input.RequestHeader

				return resp, nil
			},
		},
		&OperationWithHandler[TestHandlerRequestValidationInput, TestHandlerRequestValidationOutput]{
			Operation: huma.Operation{
				OperationID: "post-validation",
				Method:      http.MethodPost,
				Path:        "/validation",
				Summary:     "Validation endpoint",
			},
			Handler: func(ctx context.Context, input *TestHandlerRequestValidationInput) (*TestHandlerRequestValidationOutput, error) {
				resp := &TestHandlerRequestValidationOutput{}
				resp.Body.Message = input.Body.Message

				return resp, nil
			},
		},
		&OperationWithHandler[TestHandlerResponseValidationInput, TestHandlerResponseValidationOutput]{
			Operation: huma.Operation{
				OperationID: "get-validation",
				Method:      http.MethodGet,
				Path:        "/validation",
				Summary:     "Validation endpoint",
			},
			Handler: func(ctx context.Context, input *TestHandlerResponseValidationInput) (*TestHandlerResponseValidationOutput, error) {
				resp := &TestHandlerResponseValidationOutput{}

				return resp, nil
			},
		},
	}

	handler := NewHandler(
		WithConfig(handlerConfig),
		WithOperationWithHandler(operationsWithHandler),
		WithRouter(chi.NewRouter()),
		WithMiddleware([]MiddlewareFunc{
			func(ctx huma.Context, api huma.API, next func(huma.Context)) {
				ctx.SetHeader(
					"X-Response-Header",
					ctx.Header("X-Request-Header-With-Middleware-And-Response-Header"),
				)

				next(ctx)
			},
		}),
		WithFormatter([]Formatter{
			NewFormatApplicationJSON(),
			NewFormatTextHTML(),
			NewFormatTextPlain(),
		}),
		WithSecuritySchemes([]SecurityScheme{
			NewSecuritySchemeBasicAuth(),
			NewSecuritySchemeAPIHeaderKeyAuth(),
			NewSecuritySchemeAPIQueryKeyAuth(),
			NewSecuritySchemeTokenAuth(),
		}),
	)

	srvr := server.New(
		server.WithConfig(serverConfig),
		server.WithHandler(handler),
		server.WithLogFunc(func(message string) {
			t.Log(message)
		}),
	)

	srvr.Start()

	require.Eventually(t, func() bool {
		req, err := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			"http://localhost:9091/health",
			nil,
		)
		if err != nil {
			return false
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return false
		}

		_ = resp.Body.Close()

		return resp.StatusCode == http.StatusOK
	}, 5*time.Second, 10*time.Millisecond)

	clnt := client.New()

	tests := []struct {
		name                 string
		request              *http.Request
		expectedHeaders      map[string]string
		expectedResponseBody string
	}{
		{
			name: "success-get-health",
			request: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "localhost:9091",
					Path:   "/health",
				},
			},
			expectedResponseBody: `{"$schema":"http://localhost:9091/schemas/TestHandlerHealthOutputBody.json","message":"I am healthy!"}`,
		},
		{
			name: "success-get-hello",
			request: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "localhost:9091",
					Path:   "/hello/World",
				},
			},
			expectedResponseBody: `{"$schema":"http://localhost:9091/schemas/TestHandlerHelloOutputBody.json","message":"Hello World!"}`,
		},
		{
			name: "success-get-request-header",
			request: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "localhost:9091",
					Path:   "/request-header",
				},
				Header: http.Header{
					"X-Request-Header": []string{"Request header value!"},
				},
			},
			expectedResponseBody: `{"$schema":"http://localhost:9091/schemas/TestHandlerRequestHeaderOutputBody.json","message":"Request header value!"}`,
		},
		{
			name: "success-get-request-header-with-middleware-and-response-header",
			request: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "localhost:9091",
					Path:   "/request-header-with-middleware-and-response-header",
				},
				Header: http.Header{
					"X-Request-Header-With-Middleware-And-Response-Header": []string{
						"Request header with middleware and response header value!",
					},
				},
			},
			expectedHeaders: map[string]string{
				"X-Response-Header": "Request header with middleware and response header value!",
			},
			expectedResponseBody: `{"$schema":"http://localhost:9091/schemas/TestHandlerRequestHeaderWithMiddlewareAndResponseHeaderOutputBody.json","message":"Request header with middleware and response header value!"}`,
		},
		{
			name: "success-post-validation",
			request: &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "localhost:9091",
					Path:   "/validation",
				},
				Body: io.NopCloser(bytes.NewBufferString(`{"message":"Hello World!"}`)),
			},
			expectedResponseBody: `{"$schema":"http://localhost:9091/schemas/TestHandlerRequestValidationOutputBody.json","message":"Hello World!"}`,
		},
		{
			name: "failure-post-validation",
			request: &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "localhost:9091",
					Path:   "/validation",
				},
				Body: io.NopCloser(bytes.NewBufferString(`{"message":""}`)),
			},
			expectedResponseBody: `{"$schema":"http://localhost:9091/schemas/ErrorModel.json","title":"Unprocessable Entity","status":422,"detail":"validation failed","errors":[{"message":"expected length >= 3","location":"body.message","value":""}]}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := clnt.Do(test.request)
			require.NoError(t, err)

			for name, value := range test.expectedHeaders {
				assert.Equal(t, value, resp.Header.Get(name))
			}

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(string(body)))
			require.NoError(t, resp.Body.Close())
		})
	}

	srvr.Stop(context.Background(), 10*time.Second)

	srvr.Done()
}
