package server

import (
	"context"
	"net/http"
	"time"
)

// Server http server.
type Server interface {
	Start()
	Stop(ctx context.Context, timeout time.Duration)
	Done()
	Errors() <-chan error
}

// Handler http handler.
type Handler interface {
	http.Handler
}
