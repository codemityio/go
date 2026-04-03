package container

import (
	"time"
)

// WithInfoFunc configuration option.
func WithInfoFunc(function LogFunc) Option {
	return func(service *DefaultContainer) {
		service.infoFunc = function
	}
}

// WithWarnFunc configuration option.
func WithWarnFunc(function LogFunc) Option {
	return func(service *DefaultContainer) {
		service.warnFunc = function
	}
}

// WithShutdownTimeout configuration option.
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(service *DefaultContainer) {
		service.shutdownTimeout = timeout
	}
}
