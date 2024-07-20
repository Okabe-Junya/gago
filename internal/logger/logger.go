package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(enabled bool) *Logger {
	if !enabled {
		return nil
	}
	return &Logger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (l *Logger) Log(msg string, key string, value interface{}) {
	if l != nil && l.logger != nil {
		l.logger.Info(msg, key, value)
	}
}
