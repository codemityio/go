package token

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// DefaultRoundTripper signs request using signer function.
type DefaultRoundTripper struct {
	next   http.RoundTripper
	signer Signer
}

// RoundTrip transparently handles request body signing.
func (t *DefaultRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	var (
		body []byte
		err  error
	)

	if request.Body != nil {
		if body, err = io.ReadAll(request.Body); err != nil {
			return nil, fmt.Errorf("round tripper error: unable to read body: %w", err)
		}
	}

	signature, err := t.signer.Sign(body)
	if err != nil {
		return nil, fmt.Errorf("unable to sign: %w", err)
	}

	request.Header.Set("Authorization", signature)

	if body != nil {
		if err = request.Body.Close(); err != nil {
			return nil, fmt.Errorf("round tripper error: unable to close the body: %w", err)
		}

		request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	r, err := t.next.RoundTrip(request)
	if err != nil {
		return nil, fmt.Errorf("http: unable to handle round trip: %w", err)
	}

	return r, nil
}
