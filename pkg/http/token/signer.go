package token

import (
	"crypto"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// DefaultSigner is a service to generate JWT token with signed payload.
type DefaultSigner struct {
	uuidProvider UUIDRandomiser
	timeProvider TimeProvider
	key          *rsa.PrivateKey
	keyID        string
	aud          string
	payloadHash  crypto.Hash
}

// IsSignerUnableToGenerateUUID reports whether err is or wraps the error Sign returns when it
// fails to generate the token's UUID.
func (ds *DefaultSigner) IsSignerUnableToGenerateUUID(err error) bool {
	return errors.Is(err, ErrSignerUnableToGenerateUUID)
}

// IsSignerUnableToSignToken reports whether err is or wraps the error Sign returns when it
// fails to sign the token.
func (ds *DefaultSigner) IsSignerUnableToSignToken(err error) bool {
	return errors.Is(err, ErrSignerUnableToSignToken)
}

// Sign function returns a JWT token.
func (ds *DefaultSigner) Sign(data []byte) (string, error) {
	var payloadHashHex string

	if len(data) > 0 {
		payloadHashHex = hex.EncodeToString(ds.payloadHash.New().Sum(data))
	}

	id, err := ds.uuidProvider.RandomUUID()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrSignerUnableToGenerateUUID, err)
	}

	claims := jwt.MapClaims{
		"jti": id,
		"iat": ds.timeProvider.Unix(),
		"aud": ds.aud,
	}

	if payloadHashHex != "" {
		claims["payload_hash"] = payloadHashHex
		claims["payload_hash_alg"] = ds.payloadHash.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	token.Header = map[string]any{
		"kid": ds.keyID,
		"typ": typ,
		"ver": ver,
		"alg": jwt.SigningMethodRS512.Name,
	}

	tkn, err := token.SignedString(ds.key)
	if err != nil {
		return "", fmt.Errorf("%w: %w: signature error", ErrSignerUnableToSignToken, err)
	}

	return tkn, nil
}
