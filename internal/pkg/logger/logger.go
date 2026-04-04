package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Level represents log level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns string representation of log level
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a structured logger
type Logger struct {
	level   Level
	output  io.Writer
	fields  Fields
	mu      sync.Mutex
	service string
}

// Fields represents log fields
type Fields map[string]interface{}

// New creates a new logger
func New(service string, level Level) *Logger {
	return &Logger{
		level:   level,
		output:  os.Stdout,
		fields:  make(Fields),
		service: service,
	}
}

// WithFields adds fields to the logger
func (l *Logger) WithFields(fields Fields) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newFields := make(Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &Logger{
		level:   l.level,
		output:  l.output,
		fields:  newFields,
		service: l.service,
	}
}

// WithField adds a single field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return l.WithFields(Fields{key: value})
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...Fields) {
	l.log(DEBUG, msg, fields...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...Fields) {
	l.log(INFO, msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...Fields) {
	l.log(WARN, msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...Fields) {
	l.log(ERROR, msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...Fields) {
	l.log(FATAL, msg, fields...)
	os.Exit(1)
}

// log writes a log entry
func (l *Logger) log(level Level, msg string, extraFields ...Fields) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Merge fields
	allFields := make(Fields)
	for k, v := range l.fields {
		allFields[k] = v
	}
	for _, f := range extraFields {
		for k, v := range f {
			allFields[k] = v
		}
	}

	// Add standard fields
	allFields["level"] = level.String()
	allFields["time"] = time.Now().Format(time.RFC3339)
	allFields["service"] = l.service
	allFields["message"] = msg

	// Add caller info
	if level >= WARN {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			parts := strings.Split(file, "/")
			if len(parts) > 2 {
				file = strings.Join(parts[len(parts)-2:], "/")
			}
			allFields["caller"] = fmt.Sprintf("%s:%d", file, line)
		}
	}

	// Marshal to JSON
	data, err := json.Marshal(allFields)
	if err != nil {
		fmt.Fprintf(l.output, "failed to marshal log entry: %v\n", err)
		return
	}

	fmt.Fprintln(l.output, string(data))

	if level == FATAL {
		os.Exit(1)
	}
}

// RequestIDKey is the key for request ID in context
type RequestIDKey struct{}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey{}, requestID)
}

// RequestIDFromContext extracts request ID from context
func RequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey{}).(string); ok {
		return requestID
	}
	return ""
}

// WithRequestID adds request ID to logger
func (l *Logger) WithRequestID(ctx context.Context) *Logger {
	if requestID := RequestIDFromContext(ctx); requestID != "" {
		return l.WithField("request_id", requestID)
	}
	return l
}
