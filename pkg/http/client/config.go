package client

import "time"

// Config is a struct containing client configuration values.
type Config struct {
	TLS       *TLS       `env:"TLS,prefix=TLS_"             json:"tls"`
	Transport *Transport `env:"TRANSPORT,prefix=TRANSPORT_" json:"transport"`
	JWT       *JWT       `env:"JWT,prefix=JWT_"             json:"jwt"`
}

// TLS secure connection configuration.
type TLS struct {
	Cert   string `env:"CERT"    json:"cert"`
	Key    string `env:"KEY"     json:"key"`
	CACert string `env:"CA_CERT" json:"caCert"`
}

// Transport contains client transport configuration.
type Transport struct {
	MaxIdleConns        int           `env:"MAX_IDLE_CONNS,default=0"          json:"maxIdleConns"`
	MaxConnsPerHost     int           `env:"MAX_CONNS_PER_HOST,default=0"      json:"maxConnsPerHost"`
	MaxIdleConnsPerHost int           `env:"MAX_IDLE_CONNS_PER_HOST,default=2" json:"maxIdleConnsPerHost"`
	Timeout             time.Duration `env:"TIMEOUT,default=10s"               json:"timeout"`
}

// JWT a JWT token signer configuration.
type JWT struct {
	ID       string `env:"ID"       json:"id"`
	KeyPath  string `env:"KEY_PATH" json:"keyPath"`
	Audience string `env:"AUDIENCE" json:"audience"`
}
