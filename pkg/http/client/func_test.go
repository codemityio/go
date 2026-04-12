package client

import (
	"testing"
)

//go:generate sh -c "./jwt.sh"
//go:generate sh -c "./mtls.sh"

func TestWithConfig(t *testing.T) {
	cfg := Config{
		TLS: &TLS{
			Cert:   "testdata/mtls/client.crt",
			Key:    "testdata/mtls/client.key",
			CACert: "testdata/mtls/ca.crt",
		},
		Transport: &Transport{
			MaxConnsPerHost:     0,
			MaxIdleConns:        0,
			MaxIdleConnsPerHost: 0,
			Timeout:             0,
		},
		JWT: &JWT{
			ID:      "key-id",
			KeyPath: "testdata/private.pem",
		},
	}

	fn := WithConfig(&cfg)

	client := DefaultClient{
		config: nil,
	}

	fn(&client)
}
