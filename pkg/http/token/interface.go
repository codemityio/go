package token

//go:generate mockgen -source interface.go -destination interface_mock_test.go -package token

import (
	"net/http"

	"github.com/google/uuid"
)

// Signer a JWT token signer.
type Signer interface {
	Sign(data []byte) (string, error)
}

// RoundTripper http round tripper.
type RoundTripper interface {
	RoundTrip(request *http.Request) (*http.Response, error)
}

// UUIDRandomiser randomised UUID generator.
type UUIDRandomiser interface {
	RandomUUID() (uuid.UUID, error)
}

// TimeProvider time provider.
type TimeProvider interface {
	Unix() int64
}
