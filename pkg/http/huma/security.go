package huma

// SecuritySchemeBasicAuth represents a basic HTTP authentication scheme.
type SecuritySchemeBasicAuth struct{}

// Key returns the unique identifier for this security scheme.
func (s *SecuritySchemeBasicAuth) Key() string { return "BasicAuth" }

// Type returns the security scheme type.
func (s *SecuritySchemeBasicAuth) Type() string { return "http" }

// Name returns the parameter name for this scheme, if any.
func (s *SecuritySchemeBasicAuth) Name() string { return "" }

// Description returns a short description of this security scheme.
func (s *SecuritySchemeBasicAuth) Description() string {
	return "Username and password based authentication method"
}

// Scheme returns the HTTP authentication scheme.
func (s *SecuritySchemeBasicAuth) Scheme() string { return "basic" }

// In returns the location of the authentication parameter.
func (s *SecuritySchemeBasicAuth) In() string { return "" }

// SecuritySchemeAPIHeaderKeyAuth represents an API key authentication scheme using a header.
type SecuritySchemeAPIHeaderKeyAuth struct{}

// Key returns the unique identifier for this security scheme.
func (s *SecuritySchemeAPIHeaderKeyAuth) Key() string { return "APIHeaderKeyAuth" }

// Type returns the security scheme type.
func (s *SecuritySchemeAPIHeaderKeyAuth) Type() string { return "apiKey" }

// Name returns the parameter name for this scheme.
func (s *SecuritySchemeAPIHeaderKeyAuth) Name() string { return "X-API-Key" }

// Description returns a short description of this security scheme.
func (s *SecuritySchemeAPIHeaderKeyAuth) Description() string {
	return "Key based authentication method"
}

// Scheme returns the HTTP authentication scheme, if any.
func (s *SecuritySchemeAPIHeaderKeyAuth) Scheme() string { return "" }

// In returns the location of the authentication parameter.
func (s *SecuritySchemeAPIHeaderKeyAuth) In() string { return "header" }

// SecuritySchemeAPIQueryKeyAuth represents an API key authentication scheme using a query parameter.
type SecuritySchemeAPIQueryKeyAuth struct{}

// Key returns the unique identifier for this security scheme.
func (s *SecuritySchemeAPIQueryKeyAuth) Key() string { return "APIQueryKeyAuth" }

// Type returns the security scheme type.
func (s *SecuritySchemeAPIQueryKeyAuth) Type() string { return "apiKey" }

// Name returns the parameter name for this scheme.
func (s *SecuritySchemeAPIQueryKeyAuth) Name() string { return "key" }

// Description returns a short description of this security scheme.
func (s *SecuritySchemeAPIQueryKeyAuth) Description() string {
	return "Key based authentication method"
}

// Scheme returns the HTTP authentication scheme, if any.
func (s *SecuritySchemeAPIQueryKeyAuth) Scheme() string { return "" }

// In returns the location of the authentication parameter.
func (s *SecuritySchemeAPIQueryKeyAuth) In() string { return "query" }

// SecuritySchemeTokenAuth represents a bearer token authentication scheme.
type SecuritySchemeTokenAuth struct{}

// Key returns the unique identifier for this security scheme.
func (s *SecuritySchemeTokenAuth) Key() string { return "TokenAuth" }

// Type returns the security scheme type.
func (s *SecuritySchemeTokenAuth) Type() string { return "http" }

// Name returns the parameter name for this scheme.
func (s *SecuritySchemeTokenAuth) Name() string { return "token" }

// Description returns a short description of this security scheme.
func (s *SecuritySchemeTokenAuth) Description() string {
	return "Token based authentication method"
}

// Scheme returns the HTTP authentication scheme.
func (s *SecuritySchemeTokenAuth) Scheme() string { return "bearer" }

// In returns the location of the authentication parameter.
func (s *SecuritySchemeTokenAuth) In() string { return "header" }
