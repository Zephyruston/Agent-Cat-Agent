package logger

import (
	"context"
	"log/slog"
	"os"
)

var (
	defaultLogger *slog.Logger
)

func init() {
	h := slog.NewJSONHandler(os.Stdout, nil)
	defaultLogger = slog.New(h)
}

func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

func WithContext(ctx context.Context) *slog.Logger {
	return defaultLogger.With(slog.Any("ctx", ctx))
}
