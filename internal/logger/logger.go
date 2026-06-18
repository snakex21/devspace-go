package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger wraps zerolog for structured logging.
type Logger struct {
	zl zerolog.Logger
}

// New creates a new Logger with the given configuration.
func New(level string, format string, w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}

	zl := zerolog.New(w).With().Timestamp().Logger()

	// Set log level
	switch level {
	case "silent":
		zl = zl.Level(zerolog.Disabled)
	case "error":
		zl = zl.Level(zerolog.ErrorLevel)
	case "warn":
		zl = zl.Level(zerolog.WarnLevel)
	case "info":
		zl = zl.Level(zerolog.InfoLevel)
	case "debug":
		zl = zl.Level(zerolog.DebugLevel)
	default:
		zl = zl.Level(zerolog.InfoLevel)
	}

	// Set format
	if format == "text" {
		zl = zl.Output(zerolog.ConsoleWriter{
			Out:        w,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		})
	}

	return &Logger{zl: zl}
}

// Info logs at info level.
func (l *Logger) Info() *zerolog.Event {
	return l.zl.Info()
}

// Debug logs at debug level.
func (l *Logger) Debug() *zerolog.Event {
	return l.zl.Debug()
}

// Warn logs at warn level.
func (l *Logger) Warn() *zerolog.Event {
	return l.zl.Warn()
}

// Error logs at error level.
func (l *Logger) Error() *zerolog.Event {
	return l.zl.Error()
}

// WithContext returns a logger with context values.
func (l *Logger) WithContext(ctx context.Context) *zerolog.Logger {
	logger := l.zl.With().Logger()
	return &logger
}

// Default returns a default logger (info level, JSON format).
func Default() *Logger {
	return New("info", "json", os.Stderr)
}

// Init initializes the global logger. Call once at startup.
func Init(level, format string) {
	zl := zerolog.New(os.Stderr).With().Timestamp().Logger()

	switch level {
	case "silent":
		zl = zl.Level(zerolog.Disabled)
	case "error":
		zl = zl.Level(zerolog.ErrorLevel)
	case "warn":
		zl = zl.Level(zerolog.WarnLevel)
	case "info":
		zl = zl.Level(zerolog.InfoLevel)
	case "debug":
		zl = zl.Level(zerolog.DebugLevel)
	default:
		zl = zl.Level(zerolog.InfoLevel)
	}

	if format == "text" {
		zl = zl.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		})
	}

	log.Logger = zl
}
