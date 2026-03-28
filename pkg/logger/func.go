package logger

import (
	"io"
	"log"
)

// WithLogger set logger.
func WithLogger(logger *log.Logger) Option {
	return func(lgr *DefaultLogger) {
		lgr.logger = logger
	}
}

// WithLevel set level value (e.g. `logger.New(logger.WithLevel(logger.LevelDebug))`).
func WithLevel(level Level) Option {
	return func(lgr *DefaultLogger) {
		var lev uint

		switch level {
		case LevelDebug:
			lev = d | i | w | e | f
		case LevelInfo:
			lev = i | w | e | f
		case LevelWarn:
			lev = w | e | f
		case LevelError:
			lev = e | f
		case LevelFatal:
			lev = f
		}

		lgr.level = lev
	}
}

// WithFormat set format value.
func WithFormat(format Format) Option {
	return func(lgr *DefaultLogger) {
		lgr.format = format
	}
}

// WithOutput set output.
func WithOutput(output io.Writer) Option {
	return func(lgr *DefaultLogger) {
		lgr.loggerMu.Lock()
		defer lgr.loggerMu.Unlock()

		lgr.logger.SetOutput(output)
	}
}

// WithBufferSize set output.
func WithBufferSize(bufferSize int) Option {
	return func(lgr *DefaultLogger) {
		lgr.ch = make(chan logEntry, bufferSize)
	}
}
