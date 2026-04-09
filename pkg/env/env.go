package env

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

// MustBuild factory, use it if the error is not worth handling.
func MustBuild[T any]() *T {
	config, err := NewEnv[T]()
	if err != nil {
		panic(err)
	}

	return config
}

// MustBuildEnv factory, use it if the error is not worth handling.
//
// Deprecated: Use MustBuild.
func MustBuildEnv[T any]() *T {
	return MustBuild[T]()
}

// New factory.
func New[T any]() (*T, error) {
	var config T

	if e := envconfig.Process(context.Background(), &config); e != nil {
		return nil, fmt.Errorf("%w: %w", ErrEnv, e)
	}

	return &config, nil
}

// NewEnv factory.
//
// Deprecated: Use New.
func NewEnv[T any]() (*T, error) {
	return New[T]()
}
