package huma

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecuritySchemeBasicAuth(t *testing.T) {
	s := NewSecuritySchemeBasicAuth()

	require.Equal(t, "BasicAuth", s.Key())
	require.Equal(t, "http", s.Type())
	require.Empty(t, s.Name())
	require.Equal(t, "Username and password based authentication method", s.Description())
	require.Equal(t, "basic", s.Scheme())
	require.Empty(t, s.In())
}

func TestSecuritySchemeAPIHeaderKeyAuth(t *testing.T) {
	s := NewSecuritySchemeAPIHeaderKeyAuth()

	require.Equal(t, "APIHeaderKeyAuth", s.Key())
	require.Equal(t, "apiKey", s.Type())
	require.Equal(t, "X-API-Key", s.Name())
	require.Equal(t, "Key based authentication method", s.Description())
	require.Empty(t, s.Scheme())
	require.Equal(t, "header", s.In())
}

func TestSecuritySchemeAPIQueryKeyAuth(t *testing.T) {
	s := NewSecuritySchemeAPIQueryKeyAuth()

	require.Equal(t, "APIQueryKeyAuth", s.Key())
	require.Equal(t, "apiKey", s.Type())
	require.Equal(t, "key", s.Name())
	require.Equal(t, "Key based authentication method", s.Description())
	require.Empty(t, s.Scheme())
	require.Equal(t, "query", s.In())
}

func TestSecuritySchemeTokenAuth(t *testing.T) {
	s := NewSecuritySchemeTokenAuth()

	require.Equal(t, "TokenAuth", s.Key())
	require.Equal(t, "http", s.Type())
	require.Equal(t, "token", s.Name())
	require.Equal(t, "Token based authentication method", s.Description())
	require.Equal(t, "bearer", s.Scheme())
	require.Equal(t, "header", s.In())
}
