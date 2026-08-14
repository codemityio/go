package client

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/codemityio/go/pkg/http/token"
	time2 "github.com/codemityio/go/pkg/time"
)

// WithConfig configuration option.
func WithConfig(config *Config) Option {
	return func(service *DefaultClient) {
		tlsConfig := &tls.Config{ //nolint:exhaustruct_v5
			MinVersion: tls.VersionTLS12,
		}

		// use CA certificate
		configCACert(tlsConfig, config)

		// use client certificate and key
		configCertKey(tlsConfig, config)

		// use JWT token key if provided
		ds := configJWT(config)

		var transport http.RoundTripper

		transport = &http.Transport{ //nolint:exhaustruct_v5
			TLSClientConfig: tlsConfig,
		}

		if config != nil && config.Transport != nil {
			transport = &http.Transport{ //nolint:exhaustruct_v5
				MaxIdleConns:        config.Transport.MaxIdleConns,
				MaxConnsPerHost:     config.Transport.MaxConnsPerHost,
				MaxIdleConnsPerHost: config.Transport.MaxIdleConnsPerHost,
				TLSClientConfig:     tlsConfig,
			}
		}

		if ds != nil {
			transport = token.NewDefaultRoundTripper(transport, ds)
		}

		service.Transport = transport

		if config != nil && config.Transport != nil {
			service.Timeout = config.Transport.Timeout
		}

		service.config = config
	}
}

func configCACert(tlsConfig *tls.Config, config *Config) {
	if config == nil || config.TLS == nil || config.TLS.CACert == "" {
		return
	}

	if _, err := os.Stat(config.TLS.CACert); err != nil {
		panic(fmt.Sprintf(
			"CA certificate file `%s` is unavailable: %v",
			config.TLS.CACert, err,
		))
	}

	var caCert []byte

	caCert, err := os.ReadFile(config.TLS.CACert)
	if err != nil {
		panic(err)
	}

	caCertPool := x509.NewCertPool()

	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		panic("unable to append certificate")
	}

	tlsConfig.RootCAs = caCertPool
}

func configCertKey(tlsConfig *tls.Config, config *Config) {
	if config == nil || config.TLS == nil || config.TLS.Cert == "" || config.TLS.Key == "" {
		return
	}

	if _, err := os.Stat(config.TLS.Cert); err != nil {
		panic(fmt.Sprintf(
			"TLS client certificate file `%s` is unavailable: %v",
			config.TLS.Cert, err,
		))
	}

	if _, err := os.Stat(config.TLS.Key); err != nil {
		panic(fmt.Sprintf(
			"TLS client key file `%s` is unavailable: %v",
			config.TLS.Key, err,
		))
	}

	var cert tls.Certificate

	cert, err := tls.LoadX509KeyPair(config.TLS.Cert, config.TLS.Key)
	if err != nil {
		panic(err)
	}

	tlsConfig.Certificates = []tls.Certificate{cert}
}

func configJWT(config *Config) *token.DefaultSigner {
	if config == nil || config.JWT == nil || config.JWT.KeyPath == "" || config.JWT.ID == "" {
		return nil
	}

	if _, err := os.Stat(config.JWT.KeyPath); err != nil {
		panic("error: JWT key file `" + config.JWT.KeyPath + "` is unavailable: " + err.Error())
	}

	var pemBytes []byte

	pemBytes, err := os.ReadFile(config.JWT.KeyPath)
	if err != nil {
		panic("error: unable to read JWT key file " + config.JWT.KeyPath)
	}

	block, _ := pem.Decode(pemBytes)

	var key *rsa.PrivateKey

	// Try to parse PKCS#1 (RSA) private key.
	// key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	// If that fails, try to parse a PKCS#8 key
	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(fmt.Sprintf("error: unable to parse decoded JWT key: %v", err))
	}

	key, ok := pk.(*rsa.PrivateKey)
	if !ok {
		panic("error: unable to assert private key")
	}

	return token.NewDefaultSigner(
		token.NewUUIDProvider(),
		time2.ProviderFunc(time.Now),
		key,
		config.JWT.ID,
		config.JWT.Audience,
	)
}
