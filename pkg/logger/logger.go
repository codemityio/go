package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

var exitFunc = os.Exit //nolint:gochecknoglobals // non exported variable

// DefaultLogger a common logger.
type DefaultLogger struct {
	level    uint
	format   Format
	logger   *log.Logger
	loggerMu sync.Mutex
	ch       chan logEntry
	mu       sync.Mutex
	done     chan struct{}
	once     sync.Once
}

// Debug a common logger debug method.
func (l *DefaultLogger) Debug(message string, fields ...Field) {
	l.log(d, LevelDebug, message, nil, fields...)
}

// Info a common logger debug method.
func (l *DefaultLogger) Info(message string, fields ...Field) {
	l.log(i, LevelInfo, message, nil, fields...)
}

// Warn a common logger debug method.
func (l *DefaultLogger) Warn(message string, fields ...Field) {
	l.log(w, LevelWarn, message, nil, fields...)
}

// Error a common logger debug method.
func (l *DefaultLogger) Error(err error, message string, fields ...Field) {
	l.log(e, LevelError, message, err, fields...)
}

// Fatal a common logger debug method.
func (l *DefaultLogger) Fatal(message string, fields ...Field) {
	l.log(f, LevelFatal, message, nil, fields...)

	exitFunc(1)
}

// Debugf a common logger debug method.
func (l *DefaultLogger) Debugf(format string, v ...any) {
	l.log(d, LevelDebug, fmt.Sprintf(format, v...), nil)
}

// Infof a common logger info method.
func (l *DefaultLogger) Infof(format string, v ...any) {
	l.log(i, LevelInfo, fmt.Sprintf(format, v...), nil)
}

// Warnf a common logger warn method.
func (l *DefaultLogger) Warnf(format string, v ...any) {
	l.log(w, LevelWarn, fmt.Sprintf(format, v...), nil)
}

// Errorf a common logger error method.
func (l *DefaultLogger) Errorf(err error, format string, v ...any) {
	l.log(e, LevelError, fmt.Sprintf(format, v...), err)
}

// Fatalf a common logger error method.
func (l *DefaultLogger) Fatalf(format string, v ...any) {
	l.log(f, LevelFatal, fmt.Sprintf(format, v...), nil)

	exitFunc(1)
}

// Close signal the worker to stop.
func (l *DefaultLogger) Close() {
	l.once.Do(func() {
		l.mu.Lock()
		close(l.ch)
		l.mu.Unlock()
	})
}

// Done signal the logger is done.
func (l *DefaultLogger) Done() {
	<-l.done
}

func (l *DefaultLogger) worker() {
	defer close(l.done)

	for entry := range l.ch {
		if l.level&entry.level > 0 {
			l.loggerMu.Lock()
			l.logger.Println(l.fmt(entry.label, entry.message, entry.err, entry.fields...))
			l.loggerMu.Unlock()
		}
	}
}

func (l *DefaultLogger) fmt(level Level, message string, e error, fields ...Field) string {
	switch l.format {
	case FormatText:
		return l.formatText(level, message, e, fields...)
	case FormatJSON:
		return l.formatJSON(level, message, e, fields...)
	default:
		return ""
	}
}

func (l *DefaultLogger) formatJSON(level Level, message string, e error, fields ...Field) string {
	msg := map[string]any{
		labelLevel:   string(level),
		labelMessage: message,
	}

	if e != nil {
		msg[labelError] = e.Error()
	}

	for _, field := range fields {
		msg[field.Key] = field.Value
	}

	output, err := json.Marshal(msg)
	if err != nil {
		l.Error(err, "unable to marshal JSON log message")

		return ""
	}

	return string(output)
}

func (l *DefaultLogger) formatText(level Level, message string, e error, fields ...Field) string {
	var errorStr string

	if e != nil {
		errorStr = fmt.Sprintf(";%s=%v", labelError, e)
	}

	msg := fmt.Sprintf("%s=%s;%s=%s%s", labelLevel, string(level), labelMessage, message, errorStr)

	output := l.formatFieldsText(fields...)

	if output != "" {
		msg = fmt.Sprintf("%s;%s", msg, output)
	}

	return msg
}

func (l *DefaultLogger) formatFieldsText(fields ...Field) string {
	if len(fields) == 0 {
		return ""
	}

	keys := make([]string, len(fields))
	keyToValue := make(map[string]any, len(fields))

	for index, field := range fields {
		keys[index] = field.Key
		keyToValue[field.Key] = field.Value
	}

	sort.Strings(keys)

	result := make([]string, len(keys))
	for index, key := range keys {
		result[index] = fmt.Sprintf("%s=%v", key, keyToValue[key])
	}

	return strings.Join(result, ";")
}

func (l *DefaultLogger) log(level uint, label Level, message string, err error, fields ...Field) {
	entry := logEntry{
		level:   level,
		label:   label,
		message: message,
		err:     err,
		fields:  fields,
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	select {
	case l.ch <- entry:
		// success
	default:
		if r := recover(); r != nil {
			l.loggerMu.Lock()
			l.logger.Println(l.fmt(LevelDebug, "logged on closed channel", nil))
			l.loggerMu.Unlock()
		}

		if l.level&level > 0 {
			l.loggerMu.Lock()
			l.logger.Println(l.fmt(label, message, err, fields...))
			l.loggerMu.Unlock()
		}
	}
}
