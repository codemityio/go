package logger

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultLogger(t *testing.T) {
	exitFunc = func(code int) {}

	tests := map[string]struct {
		level          Level
		format         Format
		error          error
		message        string
		fields         []Field
		expectedResult string
	}{
		"debug-text": {
			level:   LevelDebug,
			format:  FormatText,
			error:   nil,
			message: "one",
			fields: []Field{
				{Key: "int", Value: 1},
				{Key: "string", Value: "one"},
				{Key: "bool", Value: true},
			},
			expectedResult: `bool=true;int=1;level=debug;message=one;string=one`,
		},
		"debug-json": {
			level:   LevelDebug,
			format:  FormatJSON,
			error:   nil,
			message: "one",
			fields: []Field{
				{Key: "int", Value: 1},
				{Key: "string", Value: "one"},
				{Key: "bool", Value: true},
			},
			expectedResult: `{"bool":true,"int":1,"level":"debug","message":"one","string":"one"}`,
		},
		"info-text": {
			level:   LevelInfo,
			format:  FormatText,
			error:   nil,
			message: "one",
			fields: []Field{
				{Key: "int", Value: 1},
				{Key: "string", Value: "one"},
				{Key: "bool", Value: true},
			},
			expectedResult: `bool=true;int=1;level=info;message=one;string=one`,
		},
		"info-json": {
			level:   LevelInfo,
			format:  FormatJSON,
			error:   nil,
			message: "one",
			fields: []Field{
				{Key: "int", Value: 1},
				{Key: "string", Value: "one"},
				{Key: "bool", Value: true},
			},
			expectedResult: `{"bool":true,"int":1,"level":"info","message":"one","string":"one"}`,
		},
		"warn-text": {
			level:   LevelWarn,
			format:  FormatText,
			error:   nil,
			message: "one",
			fields: []Field{
				{Key: "int", Value: 1},
				{Key: "string", Value: "one"},
				{Key: "bool", Value: true},
			},
			expectedResult: `bool=true;int=1;level=warning;message=one;string=one`,
		},
		"warn-json": {
			level:   LevelWarn,
			format:  FormatJSON,
			error:   nil,
			message: "one",
			fields: []Field{
				{Key: "int", Value: 1},
				{Key: "string", Value: "one"},
				{Key: "bool", Value: true},
			},
			expectedResult: `{"bool":true,"int":1,"level":"warning","message":"one","string":"one"}`,
		},
		"error-text": {
			level:   LevelError,
			format:  FormatText,
			error:   errors.New("error"),
			message: "one",
			fields: []Field{
				{Key: "int", Value: 1},
				{Key: "string", Value: "one"},
				{Key: "bool", Value: true},
			},
			expectedResult: `bool=true;error=error;int=1;level=error;message=one;string=one`,
		},
		"error-json": {
			level:   LevelError,
			format:  FormatJSON,
			error:   errors.New("error"),
			message: "one",
			fields: []Field{
				{Key: "int", Value: 1},
				{Key: "string", Value: "one"},
				{Key: "bool", Value: true},
			},
			expectedResult: `{"bool":true,"int":1,"level":"error","message":"one","error":"error","string":"one"}`,
		},
		"fatal-text": {
			level:   LevelFatal,
			format:  FormatText,
			error:   nil,
			message: "one",
			fields: []Field{
				{Key: "int", Value: 1},
				{Key: "string", Value: "one"},
				{Key: "bool", Value: true},
			},
			expectedResult: `bool=true;int=1;level=fatal;message=one;string=one`,
		},
		"fatal-json": {
			level:   LevelFatal,
			format:  FormatJSON,
			error:   nil,
			message: "one",
			fields: []Field{
				{Key: "int", Value: 1},
				{Key: "string", Value: "one"},
				{Key: "bool", Value: true},
			},
			expectedResult: `{"bool":true,"int":1,"level":"fatal","message":"one","string":"one"}`,
		},
	}

	for i, test := range tests {
		t.Run(i, func(t *testing.T) {
			// Create a buffer to capture the output
			var buf bytes.Buffer

			lgr := New(
				WithLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)),
				WithLevel(test.level),
				WithFormat(test.format),
				WithOutput(&buf),
				WithBufferSize(100),
			)

			switch test.level {
			case LevelDebug:
				lgr.Debug(test.message, test.fields...)
			case LevelInfo:
				lgr.Info(test.message, test.fields...)
			case LevelWarn:
				lgr.Warn(test.message, test.fields...)
			case LevelError:
				lgr.Error(test.error, test.message, test.fields...)
			case LevelFatal:
				lgr.Fatal(test.message, test.fields...)
			default:
				t.Fatalf("Wrong %s level used", test.level)
			}

			pattern := `^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} `

			re, err := regexp.Compile(pattern)
			require.NoError(t, err)

			time.Sleep(10 * time.Millisecond)

			lgr.Close()
			lgr.Done()

			trimmed := re.ReplaceAllString(strings.TrimSpace(buf.String()), "")

			if test.format == FormatText {
				mp := make(map[string]string)

				pairs := strings.SplitSeq(trimmed, ";")

				for pair := range pairs {
					parts := strings.SplitN(pair, "=", 2)
					if len(parts) == 2 {
						mp[parts[0]] = parts[1]
					}
				}

				keys := make([]string, 0, len(mp))
				for key := range mp {
					keys = append(keys, key)
				}

				sort.Strings(keys)

				result := make([]string, 0, len(keys))
				for _, key := range keys {
					result = append(result, fmt.Sprintf("%s=%s", key, mp[key]))
				}

				trimmed = strings.Join(result, ";")

				assert.Equal(t, test.expectedResult, trimmed)
			} else {
				assert.JSONEq(t, test.expectedResult, trimmed)
			}
		})
	}
}

func TestDefaultLoggerFormatted(t *testing.T) {
	exitFunc = func(code int) {}

	tests := map[string]struct {
		level          Level
		format         Format
		error          error
		message        string
		vals           []any
		expectedResult string
	}{
		"debug-text": {
			level:   LevelDebug,
			format:  FormatText,
			error:   nil,
			message: "one;int=%d;string=%s;bool=%t",
			vals: []any{
				1,
				"one",
				true,
			},
			expectedResult: `bool=true;int=1;level=debug;message=one;string=one`,
		},
		"debug-json": {
			level:   LevelDebug,
			format:  FormatJSON,
			error:   nil,
			message: "one;int=%d;string=%s;bool=%t",
			vals: []any{
				1,
				"one",
				true,
			},
			expectedResult: `{"level":"debug", "message":"one;int=1;string=one;bool=true"}`,
		},
		"info-text": {
			level:   LevelInfo,
			format:  FormatText,
			error:   nil,
			message: "one;int=%d;string=%s;bool=%t",
			vals: []any{
				1,
				"one",
				true,
			},
			expectedResult: `bool=true;int=1;level=info;message=one;string=one`,
		},
		"info-json": {
			level:   LevelInfo,
			format:  FormatJSON,
			error:   nil,
			message: "one;int=%d;string=%s;bool=%t",
			vals: []any{
				1,
				"one",
				true,
			},
			expectedResult: `{"level":"info", "message":"one;int=1;string=one;bool=true"}`,
		},
		"warn-text": {
			level:   LevelWarn,
			format:  FormatText,
			error:   nil,
			message: "one;int=%d;string=%s;bool=%t",
			vals: []any{
				1,
				"one",
				true,
			},
			expectedResult: `bool=true;int=1;level=warning;message=one;string=one`,
		},
		"warn-json": {
			level:   LevelWarn,
			format:  FormatJSON,
			error:   nil,
			message: "one;int=%d;string=%s;bool=%t",
			vals: []any{
				1,
				"one",
				true,
			},
			expectedResult: `{"level":"warning", "message":"one;int=1;string=one;bool=true"}`,
		},
		"error-text": {
			level:   LevelError,
			format:  FormatText,
			error:   errors.New("error"),
			message: "one;int=%d;string=%s;bool=%t",
			vals: []any{
				1,
				"one",
				true,
			},
			expectedResult: `bool=true;error=error;int=1;level=error;message=one;string=one`,
		},
		"error-json": {
			level:   LevelError,
			format:  FormatJSON,
			error:   errors.New("error"),
			message: "one;int=%d;string=%s;bool=%t",
			vals: []any{
				1,
				"one",
				true,
			},
			expectedResult: `{"error":"error", "level":"error", "message":"one;int=1;string=one;bool=true"}`,
		},
		"fatal-text": {
			level:   LevelFatal,
			format:  FormatText,
			error:   nil,
			message: "one;int=%d;string=%s;bool=%t",
			vals: []any{
				1,
				"one",
				true,
			},
			expectedResult: `bool=true;int=1;level=fatal;message=one;string=one`,
		},
		"fatal-json": {
			level:   LevelFatal,
			format:  FormatJSON,
			error:   nil,
			message: "one;int=%d;string=%s;bool=%t",
			vals: []any{
				1,
				"one",
				true,
			},
			expectedResult: `{"level":"fatal", "message":"one;int=1;string=one;bool=true"}`,
		},
	}

	for i, test := range tests {
		t.Run(i, func(t *testing.T) {
			// Create a buffer to capture the output
			var buf bytes.Buffer

			lgr := New(
				WithLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)),
				WithLevel(test.level),
				WithFormat(test.format),
				WithOutput(&buf),
				WithBufferSize(100),
			)

			switch test.level {
			case LevelDebug:
				lgr.Debugf(test.message, test.vals...)
			case LevelInfo:
				lgr.Infof(test.message, test.vals...)
			case LevelWarn:
				lgr.Warnf(test.message, test.vals...)
			case LevelError:
				lgr.Errorf(test.error, test.message, test.vals...)
			case LevelFatal:
				lgr.Fatalf(test.message, test.vals...)
			default:
				t.Fatalf("Wrong %s level used", test.level)
			}

			pattern := `^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} `

			re, err := regexp.Compile(pattern)
			require.NoError(t, err)

			time.Sleep(10 * time.Millisecond)

			lgr.Close()
			lgr.Done()

			trimmed := re.ReplaceAllString(strings.TrimSpace(buf.String()), "")

			if test.format == FormatText {
				mp := make(map[string]string)

				pairs := strings.SplitSeq(trimmed, ";")

				for pair := range pairs {
					parts := strings.SplitN(pair, "=", 2)
					if len(parts) == 2 {
						mp[parts[0]] = parts[1]
					}
				}

				keys := make([]string, 0, len(mp))
				for key := range mp {
					keys = append(keys, key)
				}

				sort.Strings(keys)

				result := make([]string, 0, len(keys))
				for _, key := range keys {
					result = append(result, fmt.Sprintf("%s=%s", key, mp[key]))
				}

				trimmed = strings.Join(result, ";")

				assert.Equal(t, test.expectedResult, trimmed)
			} else {
				assert.JSONEq(t, test.expectedResult, trimmed)
			}
		})
	}
}

func TestDefaultLogger_CloseAndDone(t *testing.T) {
	tests := map[string]struct {
		level   Level
		format  Format
		message string
	}{
		"close-done-text": {
			level:   LevelInfo,
			format:  FormatText,
			message: "close test %d",
		},
		"close-done-json": {
			level:   LevelInfo,
			format:  FormatJSON,
			message: "close test %d",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer

			logger := New(
				WithLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)),
				WithLevel(test.level),
				WithFormat(test.format),
				WithOutput(&buf),
				WithBufferSize(10),
			)

			// log a few messages
			for i := range 3 {
				logger.Infof(test.message, i)
			}

			// close and wait for shutdown
			logger.Close()

			// done should return immediately after Close
			doneReturned := make(chan struct{})

			go func() {
				logger.Done()
				close(doneReturned)
			}()

			select {
			case <-doneReturned:
				// good
			case <-time.After(time.Second):
				t.Fatal("logger.Done() did not return after Close()")
			}

			// check logs
			output := buf.String()

			for i := range 3 {
				expected := fmt.Sprintf("close test %d", i)
				if !strings.Contains(output, expected) {
					t.Errorf("Expected log message %q in output, got:\n%s", expected, output)
				}
			}
		})
	}
}

func BenchmarkDefaultLogger(b *testing.B) {
	tests := map[string]struct {
		level  Level
		format Format
	}{
		"debug-text": {
			level:  LevelDebug,
			format: FormatText,
		},
		"debug-json": {
			level:  LevelDebug,
			format: FormatJSON,
		},
		"info-text": {
			level:  LevelInfo,
			format: FormatText,
		},
		"info-json": {
			level:  LevelInfo,
			format: FormatJSON,
		},
		"error-text": {
			level:  LevelError,
			format: FormatText,
		},
		"error-json": {
			level:  LevelError,
			format: FormatJSON,
		},
	}

	for name, test := range tests {
		b.Run(name, func(b *testing.B) {
			var buf bytes.Buffer

			logger := New(
				WithLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)),
				WithLevel(test.level),
				WithFormat(test.format),
				WithOutput(&buf),
				WithBufferSize(100),
			)

			b.ResetTimer()

			for i := range b.N {
				switch test.level {
				case LevelDebug:
					logger.Debugf("benchmarking debug level %d", i)
				case LevelInfo:
					logger.Infof("benchmarking info level %d", i)
				case LevelError:
					logger.Errorf(nil, "benchmarking error level %d", i)
				case LevelWarn, LevelFatal:
				default:
					b.Fatalf("Wrong %s level used", test.level)
				}
			}

			b.StopTimer()

			logger.Close()
			logger.Done()
		})
	}
}
