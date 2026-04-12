package token

import (
	"crypto"
	"crypto/rsa"
	"net/http"
)

// NewDefaultRoundTripper returns an instance of DefaultRoundTripper.
func NewDefaultRoundTripper(next http.RoundTripper, signer Signer) *DefaultRoundTripper {
	return &DefaultRoundTripper{
		next:   next,
		signer: signer,
	}
}

// NewDefaultSigner creates a new token signer.
func NewDefaultSigner(
	uuidProvider UUIDRandomiser,
	timeProvider TimeProvider,
	key *rsa.PrivateKey,
	keyID string, // md5 fingerprint of the x.509 raw certificate
	aud string,
) *DefaultSigner {
	return &DefaultSigner{
		uuidProvider: uuidProvider,
		timeProvider: timeProvider,
		key:          key,
		keyID:        keyID,
		aud:          aud,
		payloadHash:  crypto.SHA512,
	}
}

// NewUUIDProvider UUIDProvider factory.
func NewUUIDProvider() *DefaultUUIDProvider {
	return &DefaultUUIDProvider{}
}
