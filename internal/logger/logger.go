// Package logger provides logging functionality for the genetic algorithm.
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"
)

// LogLevel represents the level of logging.
type LogLevel int

const (
	// LevelDebug is the debug level.
	LevelDebug LogLevel = iota
	// LevelInfo is the info level.
	LevelInfo
	// LevelWarn is the warning level.
	LevelWarn
	// LevelError is the error level.
	LevelError
)

// Logger wraps slog.Logger to provide genetic algorithm-specific logging.
type Logger struct {
	logger *slog.Logger
	level  LogLevel
}

// LoggerOption is a function that configures a Logger.
type LoggerOption func(*Logger)

// NewLogger creates a new logger with the specified options.
func NewLogger(enabled bool, options ...LoggerOption) *Logger {
	if !enabled {
		return nil
	}

	// Create the default logger
	l := &Logger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
		level: LevelInfo,
	}

	// Apply options
	for _, option := range options {
		option(l)
	}

	return l
}

// WithLevel sets the logging level.
func WithLevel(level LogLevel) LoggerOption {
	return func(l *Logger) {
		l.level = level
		var slogLevel slog.Level
		switch level {
		case LevelDebug:
			slogLevel = slog.LevelDebug
		case LevelInfo:
			slogLevel = slog.LevelInfo
		case LevelWarn:
			slogLevel = slog.LevelWarn
		case LevelError:
			slogLevel = slog.LevelError
		}

		// Update the handler with the new level
		handlerOptions := &slog.HandlerOptions{
			Level: slogLevel,
		}

		// Create a new handler with the same output as the old one
		handler := l.logger.Handler()
		switch handler.(type) {
		case *slog.TextHandler:
			l.logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
		case *slog.JSONHandler:
			l.logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
		}
	}
}

// WithJSON sets the logger to use JSON format.
func WithJSON() LoggerOption {
	return func(l *Logger) {
		handlerOptions := &slog.HandlerOptions{
			Level: slogLevelFromLogLevel(l.level),
		}
		l.logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
	}
}

// WithWriter sets the writer for the logger.
func WithWriter(w io.Writer) LoggerOption {
	return func(l *Logger) {
		handlerOptions := &slog.HandlerOptions{
			Level: slogLevelFromLogLevel(l.level),
		}

		handler := l.logger.Handler()
		if _, ok := handler.(*slog.TextHandler); ok {
			l.logger = slog.New(slog.NewTextHandler(w, handlerOptions))
		} else {
			l.logger = slog.New(slog.NewJSONHandler(w, handlerOptions))
		}
	}
}

// Debug logs a message at debug level.
func (l *Logger) Debug(msg string, args ...any) {
	if l != nil && l.logger != nil {
		l.logger.Debug(msg, args...)
	}
}

// Info logs a message at info level.
func (l *Logger) Info(msg string, args ...any) {
	if l != nil && l.logger != nil {
		l.logger.Info(msg, args...)
	}
}

// Warn logs a message at warning level.
func (l *Logger) Warn(msg string, args ...any) {
	if l != nil && l.logger != nil {
		l.logger.Warn(msg, args...)
	}
}

// Error logs a message at error level.
func (l *Logger) Error(msg string, args ...any) {
	if l != nil && l.logger != nil {
		l.logger.Error(msg, args...)
	}
}

// WithContext returns a Logger that includes context information.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if l == nil || l.logger == nil {
		return nil
	}

	// slog.Logger does not have a WithContext method,
	// so add context information as attributes instead
	ctxLogger := l.logger
	if ctx != nil {
		// Example of extracting important information from the context
		if reqID, ok := ctx.Value("request_id").(string); ok {
			ctxLogger = ctxLogger.With("request_id", reqID)
		}
	}

	return &Logger{
		logger: ctxLogger,
		level:  l.level,
	}
}

// WithGroup returns a Logger that includes a group.
func (l *Logger) WithGroup(name string) *Logger {
	if l == nil || l.logger == nil {
		return nil
	}
	return &Logger{
		logger: l.logger.WithGroup(name),
		level:  l.level,
	}
}

// LogGenerationStats logs statistics about a generation.
func (l *Logger) LogGenerationStats(generation int, stats map[string]interface{}, elapsed time.Duration) {
	if l == nil || l.logger == nil {
		return
	}

	attrs := []any{
		slog.Int("generation", generation),
		slog.Duration("elapsed", elapsed),
	}

	for k, v := range stats {
		switch val := v.(type) {
		case float64:
			attrs = append(attrs, slog.Float64(k, val))
		case int:
			attrs = append(attrs, slog.Int(k, val))
		case string:
			attrs = append(attrs, slog.String(k, val))
		case bool:
			attrs = append(attrs, slog.Bool(k, val))
		}
	}

	l.logger.Info("Generation stats", attrs...)
}

// LogError logs an error with context.
func (l *Logger) LogError(err error, msg string, args ...any) {
	if l == nil || l.logger == nil || err == nil {
		return
	}

	allArgs := append([]any{slog.Any("error", err)}, args...)
	l.logger.Error(msg, allArgs...)
}

// slogLevelFromLogLevel converts LogLevel to slog.Level
func slogLevelFromLogLevel(level LogLevel) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
