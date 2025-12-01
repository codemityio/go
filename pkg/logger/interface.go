// Package logger contains a common logger.
package logger

// Field a custom field struct representation.
type Field struct {
	Key   string
	Value any
}

// Logger is a logger interface.
type Logger interface {
	DebugLogger
	InfoLogger
	WarnLogger
	ErrorLogger
	FatalLogger
	DebugfLogger
	InfofLogger
	WarnfLogger
	ErrorfLogger
	FatalfLogger
}

type DebugLogger interface {
	Debug(message string, fields ...Field)
}

type InfoLogger interface {
	Info(message string, fields ...Field)
}

type WarnLogger interface {
	Warn(message string, fields ...Field)
}

type ErrorLogger interface {
	Error(err error, message string, fields ...Field)
}

type FatalLogger interface {
	Fatal(message string, fields ...Field)
}

// DebugfLogger logger interface.
type DebugfLogger interface {
	Debugf(format string, v ...any)
}

// InfofLogger info logger.
type InfofLogger interface {
	Infof(format string, v ...any)
}

// WarnfLogger warn logger.
type WarnfLogger interface {
	Warnf(format string, v ...any)
}

// ErrorfLogger error logger.
type ErrorfLogger interface {
	Errorf(err error, format string, v ...any)
}

// FatalfLogger fatal logger.
type FatalfLogger interface {
	Fatalf(format string, v ...any)
}
