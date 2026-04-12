package token

import (
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate sh -c "./jwt.sh"

//go:embed testdata/rsa.pem
var privateKey []byte

func TestDefaultSigner_Sign(t *testing.T) {
	block, _ := pem.Decode(privateKey)
	require.NotNil(t, block, "failed to decode PEM block containing private key")

	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	require.NoError(t, err)

	key, ok := pk.(*rsa.PrivateKey)
	require.True(t, ok, "not an RSA private key")

	namespace := uuid.UUID{}
	hash := md5.New()
	hash.Write([]byte(""))
	defaultJTI := uuid.NewMD5(namespace, hash.Sum(nil))

	tests := []struct {
		name       string
		data       []byte
		jti        uuid.UUID
		unixTime   int64
		keyID      string
		wantIat    int64
		wantAud    string
		wantClaims func(t *testing.T, wantIat int64, wantAud string, claims jwt.MapClaims)
	}{
		{
			name:     "signs payload and includes payload_hash claims",
			data:     []byte("Hello, World!"),
			jti:      defaultJTI,
			unixTime: 1234567890,
			wantIat:  int64(1234567890),
			wantAud:  "audience",
			wantClaims: func(t *testing.T, wantIat int64, wantAud string, claims jwt.MapClaims) {
				t.Helper()
				assert.Equal(t, defaultJTI.String(), claims["jti"])
				iat, ok := claims["iat"].(float64)
				require.True(t, ok, "iat claim is not a float64")
				assert.Equal(t, wantIat, int64(iat))
				assert.Equal(t, wantAud, claims["aud"])
				assert.NotEmpty(t, claims["payload_hash"])
				assert.NotEmpty(t, claims["payload_hash_alg"])
			},
		},
		{
			name:     "omits payload_hash claims for empty payload",
			data:     []byte{},
			jti:      defaultJTI,
			unixTime: 1234567890,
			wantIat:  int64(1234567890),
			wantAud:  "key-id",
			wantClaims: func(t *testing.T, wantIat int64, wantAud string, claims jwt.MapClaims) {
				t.Helper()
				assert.Equal(t, defaultJTI.String(), claims["jti"])
				iat, ok := claims["iat"].(float64)
				require.True(t, ok, "iat claim is not a float64")
				assert.Equal(t, wantIat, int64(iat))
				assert.Equal(t, wantAud, claims["aud"])
				assert.Nil(t, claims["payload_hash"])
				assert.Nil(t, claims["payload_hash_alg"])
			},
		},
		{
			name:     "uses provided key ID in header",
			data:     []byte("Hello, World!"),
			jti:      defaultJTI,
			unixTime: 1234567890,
			wantIat:  int64(1234567890),
			wantAud:  "another-key-id",
			wantClaims: func(t *testing.T, wantIat int64, wantAud string, claims jwt.MapClaims) {
				t.Helper()
				assert.Equal(t, defaultJTI.String(), claims["jti"])
				iat, ok := claims["iat"].(float64)
				require.True(t, ok, "iat claim is not a float64")
				assert.Equal(t, wantIat, int64(iat))
				assert.Equal(t, wantAud, claims["aud"])
				assert.NotEmpty(t, claims["payload_hash"])
				assert.NotEmpty(t, claims["payload_hash_alg"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uuidProvider := NewMockUUIDRandomiser(ctrl)
			uuidProvider.EXPECT().RandomUUID().Return(tt.jti, nil).Times(1)

			timeProvider := NewMockTimeProvider(ctrl)
			timeProvider.EXPECT().Unix().Return(tt.unixTime).Times(1)

			signer := NewDefaultSigner(
				uuidProvider,
				timeProvider,
				key,
				tt.keyID,
				tt.wantAud,
			)

			res, err := signer.Sign(tt.data)
			require.NoError(t, err)
			require.NotEmpty(t, res)

			parsed, err := jwt.Parse(res, func(tok *jwt.Token) (any, error) {
				_, ok := tok.Method.(*jwt.SigningMethodRSA)
				require.True(t, ok, "unexpected signing method: %v", tok.Header["alg"])

				return &key.PublicKey, nil
			})
			require.NoError(t, err)
			require.True(t, parsed.Valid)

			assert.Equal(t, "RS512", parsed.Header["alg"])
			assert.Equal(t, tt.keyID, parsed.Header["kid"])
			assert.Equal(t, typ, parsed.Header["typ"])
			assert.Equal(t, ver, parsed.Header["ver"])

			claims, ok := parsed.Claims.(jwt.MapClaims)
			require.True(t, ok, "unexpected claims type")

			tt.wantClaims(t, tt.wantIat, tt.wantAud, claims)
		})
	}
}
