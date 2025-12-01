package logger

import (
	"log"
	"os"
	"sync"
)

// New a common logger factory.
func New(options ...Option) *DefaultLogger {
	lgr := &DefaultLogger{
		level:    d | i | w | e | f,
		format:   FormatText,
		logger:   log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC|log.Lmicroseconds),
		loggerMu: sync.Mutex{},
		ch:       make(chan logEntry, BufferSizeDefault),
		mu:       sync.Mutex{},
		done:     make(chan struct{}),
		once:     sync.Once{},
	}

	for _, opt := range options {
		opt(lgr)
	}

	go lgr.worker()

	return lgr
}
