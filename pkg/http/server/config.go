package server

import "time"

// Config is a struct containing server configuration values.
type Config struct {
	Port             int           `env:"PORT,default=8080"              json:"port"`
	ReadTimeout      time.Duration `env:"READ_TIMEOUT,default=30s"       json:"readTimeout"`
	WriteTimeout     time.Duration `env:"WRITE_TIMEOUT,default=30s"      json:"writeTimeout"`
	ErrorsBufferSize int           `env:"ERRORS_BUFFER_SIZE,default=100" json:"errorsBufferSize"`
	TLS              *TLS          `env:"TLS,prefix=TLS_"                json:"tls"`
}

// TLS secure connection configuration.
type TLS struct {
	Cert   string `env:"CERT"    json:"cert"`
	Key    string `env:"KEY"     json:"key"`
	CACert string `env:"CA_CERT" json:"caCert"`
}
