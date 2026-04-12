package client

import (
	"net/http"
)

// DefaultClient a client struct.
type DefaultClient struct {
	http.Client

	config *Config
}

// New a factory function for the client.
func New(options ...Option) *DefaultClient {
	dc := DefaultClient{
		Client: http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		},
		config: nil,
	}

	for _, option := range options {
		option(&dc)
	}

	return &dc
}
