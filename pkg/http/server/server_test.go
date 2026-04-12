//nolint:lll
package server

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/codemityio/go/pkg/env"
	huma2 "github.com/codemityio/go/pkg/http/huma"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestDefaultServer(t *testing.T) {
	ctx := context.Background()

	t.Setenv("PORT", "8080")
	t.Setenv("READ_TIMEOUT", "10s")
	t.Setenv("WRITE_TIMEOUT", "10s")
	t.Setenv("TLS_CERT", "")
	t.Setenv("TLS_KEY", "")
	t.Setenv("TLS_CA_CERT", "")

	config, err := env.New[Config]()
	require.NoError(t, err)

	server := New(
		WithConfig(config),
	)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server.Start()

	require.Eventually(t, func() bool {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"http://localhost:8080/health",
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/health", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	server.Stop(ctx, 10*time.Second)

	server.Done()
}

func TestServerHuma(t *testing.T) {
	ctx := context.Background()

	t.Setenv("PORT", "8080")
	t.Setenv("READ_TIMEOUT", "0")
	t.Setenv("WRITE_TIMEOUT", "0")
	t.Setenv("TITLE", "API")
	t.Setenv("VERSION", "v0.0.0")
	t.Setenv("OPEN_API_PATH", "/openapi")
	t.Setenv("DOCS_PATH", "/docs")
	t.Setenv("SCHEMAS_PATH", "/schemas")
	t.Setenv("DEFAULT_FORMAT", "application/json")

	operationsWithHandler := []huma2.OperationWithHandlerRegisterer{
		&huma2.OperationWithHandler[
			GetHelloWorldTextPlainInput,
			GetHelloWorldTextPlainOutput,
		]{
			Operation: huma.Operation{
				OperationID: "text-plain",
				Method:      http.MethodGet,
				Path:        "/text-plain",
				Summary:     "Hello world!",
				Responses: map[string]*huma.Response{
					http.StatusText(http.StatusOK): {
						Content: map[string]*huma.MediaType{
							huma2.ContentTypeTextPlain: {},
						},
					},
				},
			},
			Handler: func(
				ctx context.Context,
				input *GetHelloWorldTextPlainInput,
			) (*GetHelloWorldTextPlainOutput, error) {
				return &GetHelloWorldTextPlainOutput{
					ContentType:  input.Accept,
					LastModified: time.Now(),
					Body:         `Hello world!`,
				}, nil
			},
		},
		&huma2.OperationWithHandler[
			GetHelloWorldTextHTMLInput,
			GetHelloWorldTextHTMLOutput,
		]{
			Operation: huma.Operation{
				OperationID: "text-html",
				Method:      http.MethodGet,
				Path:        "/text-html",
				Summary:     "Hello world!",
				Responses: map[string]*huma.Response{
					http.StatusText(http.StatusOK): {
						Content: map[string]*huma.MediaType{
							huma2.ContentTypeTextHTML: {},
						},
					},
				},
			},
			Handler: func(
				ctx context.Context,
				input *GetHelloWorldTextHTMLInput,
			) (*GetHelloWorldTextHTMLOutput, error) {
				return &GetHelloWorldTextHTMLOutput{
					ContentType:  input.Accept,
					LastModified: time.Now(),
					Body:         `<h1>Hello world!</h1>`,
				}, nil
			},
		},
		&huma2.OperationWithHandler[
			GetHelloWorldApplicationJSONInput,
			GetHelloWorldApplicationJSONOutput,
		]{
			Operation: huma.Operation{
				OperationID: "application-json",
				Method:      http.MethodGet,
				Path:        "/application-json",
				Summary:     "Hello world!",
				Responses: map[string]*huma.Response{
					http.StatusText(http.StatusOK): {
						Content: map[string]*huma.MediaType{
							huma2.ContentTypeApplicationJSON: {},
						},
					},
				},
			},
			Handler: func(
				ctx context.Context,
				input *GetHelloWorldApplicationJSONInput,
			) (*GetHelloWorldApplicationJSONOutput, error) {
				return &GetHelloWorldApplicationJSONOutput{
					ContentType:  input.Accept,
					LastModified: time.Now(),
					Body:         GetHelloWorldApplicationJSONOutputBody{Message: "Hello world!"},
				}, nil
			},
		},
	}

	serverConfig, err := env.New[Config]()
	require.NoError(t, err)

	handlerConfig, err := env.New[huma2.Config]()
	require.NoError(t, err)

	srvr := New(
		WithConfig(serverConfig),
		WithHandler(huma2.NewHandler(
			huma2.WithConfig(handlerConfig),
			huma2.WithOperationWithHandler(operationsWithHandler),
			huma2.WithRouter(chi.NewRouter()),
			huma2.WithMiddleware([]huma2.MiddlewareFunc{
				huma2.NewCorrelationIDMiddleware(
					huma2.WithCorrelationIDMiddlewareCorrelationIDProvider(
						huma2.CorrelationIDProviderFunc(func() string {
							return "correlation-id"
						}),
					),
				).Handler,
				huma2.NewAccessLogMiddleware(
					huma2.WithAccessLogMiddlewareLogFunc(
						func(message string, fields ...huma2.LogField) {
							t.Log(message)
							t.Log(fields)
						},
					),
				).Handler,
			}),
			huma2.WithFormatter([]huma2.Formatter{
				huma2.NewFormatApplicationJSON(),
				huma2.NewFormatTextHTML(),
				huma2.NewFormatTextPlain(),
			}),
			huma2.WithSecuritySchemes([]huma2.SecurityScheme{
				huma2.NewSecuritySchemeBasicAuth(),
				huma2.NewSecuritySchemeAPIHeaderKeyAuth(),
				huma2.NewSecuritySchemeAPIQueryKeyAuth(),
				huma2.NewSecuritySchemeTokenAuth(),
			}),
		)),
		WithLogFunc(func(message string) {
			t.Log(message)
		}),
		WithShutdownFunc(func() {
			t.Log("done")
		}),
	)

	srvr.Start()

	require.Eventually(t, func() bool {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"http://localhost:8080/text-plain",
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

	tests := []struct {
		name            string
		path            string
		accept          string
		wantStatus      int
		wantContentType string
		wantBody        string
	}{
		// -------- /text-plain --------
		{
			name:            "text-plain accept text/plain",
			path:            "/text-plain",
			accept:          "text/plain",
			wantStatus:      200,
			wantContentType: "text/plain",
			wantBody:        "\"Hello world!\"",
		},
		{
			name:            "text-plain accept text/html",
			path:            "/text-plain",
			accept:          "text/html",
			wantStatus:      200,
			wantContentType: "text/html",
			wantBody: `<!doctype html>
<html>
  <body>
    <pre>&#34;Hello world!&#34;</pre>
  </body>
</html>
`,
		},
		{
			name:            "text-plain accept application/json",
			path:            "/text-plain",
			accept:          "application/json",
			wantStatus:      200,
			wantContentType: "application/json",
			wantBody: `"Hello world!"
`,
		},

		// -------- /text-html --------
		{
			name:            "text-html accept text/plain",
			path:            "/text-html",
			accept:          "text/plain",
			wantStatus:      200,
			wantContentType: "text/plain",
			wantBody:        "\"\\u003ch1\\u003eHello world!\\u003c/h1\\u003e\"",
		},
		{
			name:            "text-html accept text/html",
			path:            "/text-html",
			accept:          "text/html",
			wantStatus:      200,
			wantContentType: "text/html",
			wantBody: `<!doctype html>
<html>
  <body>
    <pre>&#34;&lt;h1&gt;Hello world!&lt;/h1&gt;&#34;</pre>
  </body>
</html>
`,
		},
		{
			name:            "text-html accept application/json",
			path:            "/text-html",
			accept:          "application/json",
			wantStatus:      200,
			wantContentType: "application/json",
			wantBody: `"\u003ch1\u003eHello world!\u003c/h1\u003e"
`,
		},

		// -------- /application-json --------
		{
			name:            "application-json accept text/plain",
			path:            "/application-json",
			accept:          "text/plain",
			wantStatus:      200,
			wantContentType: "text/plain",
			wantBody:        `{"$schema":"http://localhost:8080/schemas/GetHelloWorldApplicationJSONOutputBody.json","message":"Hello world!"}`,
		},
		{
			name:            "application-json accept text/html",
			path:            "/application-json",
			accept:          "text/html",
			wantStatus:      200,
			wantContentType: "text/html",
			wantBody: `<!doctype html>
<html>
  <body>
    <pre>{&#34;$schema&#34;:&#34;http://localhost:8080/schemas/GetHelloWorldApplicationJSONOutputBody.json&#34;,&#34;message&#34;:&#34;Hello world!&#34;}</pre>
  </body>
</html>
`,
		},
		{
			name:            "application-json accept application/json",
			path:            "/application-json",
			accept:          "application/json",
			wantStatus:      200,
			wantContentType: "application/json",
			wantBody: `{"$schema":"http://localhost:8080/schemas/GetHelloWorldApplicationJSONOutputBody.json","message":"Hello world!"}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(
				t.Context(),
				http.MethodGet,
				"http://localhost:8080"+tt.path,
				nil,
			)
			require.NoError(t, err)

			req.Header.Set("Accept", tt.accept)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			defer func() {
				_ = res.Body.Close()
			}()

			bodyBytes, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			body := string(bodyBytes)

			require.Equal(t, tt.wantStatus, res.StatusCode)
			require.Contains(t, res.Header.Get("Content-Type"), tt.wantContentType)
			require.Equal(t, tt.wantBody, body)
		})
	}

	srvr.Stop(context.Background(), 10*time.Second)

	srvr.Done()
}
