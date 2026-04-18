//nolint:lll
package huma

import (
	"bytes"
	"context"
	"encoding/json"
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

// request header with middleware and response header.
type TestMiddlewareRequestHeaderHeaderInput struct {
	HeaderRequest string `header:"X-Request-Header-With-Middleware-And-Response-Header"` //nolint:tagliatelle
}

type TestMiddlewareRequestHeaderHeaderOutputBody struct {
	Message string `json:"message"`
}

type TestMiddlewareRequestHeaderHeaderOutput struct {
	ResponseHeader string                                      `header:"X-Response-Header"` //nolint:tagliatelle
	Body           TestMiddlewareRequestHeaderHeaderOutputBody `                           required:"true"`
}

func TestMiddleware(t *testing.T) {
	ctx := context.Background()

	t.Setenv("PORT", "9092")
	t.Setenv("READ_TIMEOUT", "0")
	t.Setenv("WRITE_TIMEOUT", "0")
	t.Setenv("TITLE", "API")
	t.Setenv("VERSION", "v0.0.0")
	t.Setenv("OPEN_API_PATH", "/openapi")
	t.Setenv("DOCS_PATH", "/docs")
	t.Setenv("SCHEMAS_PATH", "/schemas")
	t.Setenv("DEFAULT_FORMAT", "application/json")

	operationsWithHandler := []OperationWithHandlerRegisterer{
		&OperationWithHandler[
			TestMiddlewareRequestHeaderHeaderInput,
			TestMiddlewareRequestHeaderHeaderOutput,
		]{
			Operation: huma.Operation{
				OperationID: "middleware",
				Method:      http.MethodGet,
				Path:        "/middleware",
				Summary:     "",
			},
			Handler: func(
				ctx context.Context,
				input *TestMiddlewareRequestHeaderHeaderInput,
			) (*TestMiddlewareRequestHeaderHeaderOutput, error) {
				resp := &TestMiddlewareRequestHeaderHeaderOutput{}
				resp.Body.Message = input.HeaderRequest

				return resp, nil
			},
		},
	}

	serverConfig, err := env.New[server.Config]()
	require.NoError(t, err)

	handlerConfig, err := env.New[Config]()
	require.NoError(t, err)

	bfr := bytes.Buffer{}

	accessLogMiddleware := NewAccessLogMiddleware(
		WithAccessLogMiddlewareLogFunc(func(message string, fields ...LogField) {
			bfr.WriteString(message)

			fStr, e := json.Marshal(fields)
			require.NoError(t, e)

			_, e = fmt.Fprintf(&bfr, "%s\n", string(fStr))
			require.NoError(t, e)
		}),
	)
	correlationIDMiddleware := NewCorrelationIDMiddleware(
		WithCorrelationIDMiddlewareCorrelationIDProvider(CorrelationIDProviderFunc(func() string {
			return "correlation-id-value"
		})),
	)

	handler := NewHandler(
		WithConfig(handlerConfig),
		WithOperationWithHandler(operationsWithHandler),
		WithRouter(chi.NewRouter()),
		WithMiddleware([]MiddlewareFunc{
			correlationIDMiddleware.Handler,
			accessLogMiddleware.Handler,
		}),
	)

	srvr := server.New(
		server.WithConfig(serverConfig),
		server.WithHandler(handler),
		server.WithLogFunc(func(message string) {
			_, e := fmt.Fprintf(&bfr, "%s\n", message)

			require.NoError(t, e)
		}),
	)

	srvr.Start()

	require.Eventually(t, func() bool {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"http://localhost:9092/middleware",
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
			name: "success-middleware",
			request: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "localhost:9092",
					Path:   "/middleware",
				},
				Header: http.Header{},
			},
			expectedHeaders:      map[string]string{},
			expectedResponseBody: `{"$schema":"http://localhost:9092/schemas/TestMiddlewareRequestHeaderHeaderOutputBody.json","message":""}`,
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

			assert.Equal(t,
				test.expectedResponseBody,
				strings.TrimSpace(string(body)),
			)
			require.NoError(t, resp.Body.Close())
		})
	}

	srvr.Stop(ctx, 0)

	srvr.Done()

	// assert logs
	l := bfr.String()

	assert.Contains(t, l, "starting server...")
	assert.Contains(t, l, "starting HTTP server, listening on port 9092")
	assert.Contains(
		t,
		l,
		`/middleware[{"Key":"method","Value":"GET"},{"Key":"scheme","Value":"http"},{"Key":"host","Value":"localhost"},{"Key":"port","Value":9092},{"Key":"path","Value":"/middleware"},{"Key":"correlation-id","Value":"correlation-id-value"}]`,
	)
	assert.Contains(t, l, "shutting down server with `0` seconds timeout...")
	assert.Contains(t, l, "server gracefully shut down...")
}
