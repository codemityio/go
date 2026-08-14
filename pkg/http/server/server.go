package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// DefaultServer is a web server.
type DefaultServer struct {
	http.Server

	logFunc   LogFunc
	config    *Config
	mu        sync.Mutex
	startOnce sync.Once
	stopOnce  sync.Once
	started   bool
	errors    chan error
	done      chan struct{}
}

// New a factory function to create a new server.
func New(options ...Option) *DefaultServer {
	ds := DefaultServer{
		Server: http.Server{ //nolint:exhaustruct_v5 // not necessary, use zero values
			ReadHeaderTimeout: 0,
			IdleTimeout:       0,
			ReadTimeout:       defaultTimeout * time.Second,
			WriteTimeout:      defaultTimeout * time.Second,
		},
		logFunc:   nil,
		config:    nil,
		mu:        sync.Mutex{},
		startOnce: sync.Once{},
		stopOnce:  sync.Once{},
		started:   false,
		errors:    nil,
		done:      nil,
	}

	for _, option := range options {
		option(&ds)
	}

	return &ds
}

// Start method to start a web server.
func (s *DefaultServer) Start() {
	s.startOnce.Do(func() {
		s.mu.Lock()
		s.errors = make(chan error, s.config.ErrorsBufferSize)
		s.started = true
		s.done = make(chan struct{})
		s.mu.Unlock()

		s.info("starting server...")

		if s.config.TLS == nil || (s.config.TLS.Cert == "" && s.config.TLS.Key == "") {
			go s.listenAndServe()

			return
		}

		if _, err := os.Stat(s.config.TLS.Cert); errors.Is(err, os.ErrNotExist) {
			s.errors <- fmt.Errorf(
				"%w: %w: with certificate path %s",
				ErrServerCertificateDoesNotExist, err, s.config.TLS.Cert,
			)
		}

		if _, err := os.Stat(s.config.TLS.Key); errors.Is(err, os.ErrNotExist) {
			s.errors <- fmt.Errorf(
				"%w: %w: with key path %s",
				ErrServerKeyDoesNotExist, err, s.config.TLS.Key,
			)
		}

		go s.listenAndServeTLS()
	})
}

// Stop method to stop a web server.
func (s *DefaultServer) Stop(ctx context.Context, timeout time.Duration) {
	s.stopOnce.Do(func() {
		s.mu.Lock()

		if !s.started {
			close(s.done)
			s.mu.Unlock()

			return
		}

		s.mu.Unlock()

		s.info(fmt.Sprintf("shutting down server with `%d` seconds timeout...", timeout))

		var cancel context.CancelFunc

		shutdownCtx := ctx

		if timeout > 0 {
			shutdownCtx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		if e := s.Shutdown(shutdownCtx); e != nil {
			s.errors <- fmt.Errorf("%w: %w", ErrServerShutdownFailure, e)

			const delay = 100

			time.Sleep(delay * time.Millisecond)
		}

		s.info("server gracefully shut down...")

		close(s.errors)
		close(s.done)
	})
}

// Done blocks until the server is fully stopped.
func (s *DefaultServer) Done() {
	if s.done == nil {
		return
	}

	<-s.done
}

// Errors subscribe to async errors.
func (s *DefaultServer) Errors() <-chan error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return nil
	}

	return s.errors
}

func (s *DefaultServer) listenAndServe() {
	s.info(fmt.Sprintf("starting HTTP server, listening on port %d", s.config.Port))

	if e := s.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
		s.errors <- fmt.Errorf("%w: %w", ErrServerUnableToListenAndServe, e)
	}
}

func (s *DefaultServer) listenAndServeTLS() {
	t, err := json.Marshal(s.config.TLS)
	if err != nil {
		s.errors <- fmt.Errorf("%w: %w: tls config", ErrServerUnableToMarshal, err)
	}

	s.info(
		fmt.Sprintf(
			"starting HTTP server, listening on port %d with tls: %s",
			s.config.Port,
			string(t),
		),
	)

	if e := s.ListenAndServeTLS(
		s.config.TLS.Cert,
		s.config.TLS.Key,
	); e != nil && !errors.Is(e, http.ErrServerClosed) {
		s.errors <- fmt.Errorf("%w: %w", ErrServerUnableToListenAndServe, e)
	}
}

func (s *DefaultServer) info(msg string) {
	if s.logFunc != nil {
		s.logFunc(msg)
	}
}
