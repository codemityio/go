package env

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

// MustBuildEnv factory, use it if the error is not worth handling.
func MustBuildEnv[T any]() *T {
	config, err := NewEnv[T]()
	if err != nil {
		panic(err)
	}

	return config
}

// NewEnv factory.
func NewEnv[T any]() (*T, error) {
	var config T

	if e := envconfig.Process(context.Background(), &config); e != nil {
		return nil, fmt.Errorf("%w: %w", ErrEnv, e)
	}

	return &config, nil
}
