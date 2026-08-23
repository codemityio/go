package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// WithConfig configuration option.
func WithConfig(config *Config) Option {
	return func(service *DefaultServer) {
		var tlsConfig *tls.Config

		if config != nil && config.TLS != nil && config.TLS.CACert != "" {
			if _, err := os.Stat(config.TLS.CACert); errors.Is(err, os.ErrNotExist) {
				panic(fmt.Sprintf("CA certificate file does not exist: %v", err))
			}

			// load CA certificate file and add it to list of client CAs
			caCertFile, err := os.ReadFile(config.TLS.CACert)
			if err != nil {
				panic(err)
			}

			certPool := x509.NewCertPool()
			certPool.AppendCertsFromPEM(caCertFile)

			tlsConfig = &tls.Config{ //nolint:exhaustruct_v5 // not necessary, use zero values
				MinVersion: tls.VersionTLS13,
				ClientAuth: tls.RequireAndVerifyClientCert,
				ClientCAs:  certPool,
			}
		}

		service.config = config

		service.ReadTimeout = config.ReadTimeout
		service.WriteTimeout = config.WriteTimeout
		service.Addr = fmt.Sprintf(":%d", config.Port)
		service.TLSConfig = tlsConfig
	}
}

// WithShutdownFunc configuration option.
func WithShutdownFunc(function func()) Option {
	return func(service *DefaultServer) {
		service.RegisterOnShutdown(function)
	}
}

// WithLogFunc configuration option.
func WithLogFunc(fn LogFunc) Option {
	return func(service *DefaultServer) {
		service.logFunc = fn
	}
}

// WithHandler configuration option.
func WithHandler(handler http.Handler) Option {
	return func(service *DefaultServer) {
		service.Handler = handler
	}
}
