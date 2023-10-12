package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func init() {
	logger := NewLogger()

	Log = logger
	slog.SetDefault(logger)
}

func newHandler() slog.Handler {
	var handler slog.Handler

	if os.Getenv("GIN_MODE") != "release" {
		handler = slog.NewTextHandler(os.Stdout, nil)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}

	return handler
}

func NewLogger() *slog.Logger {
	handler := newHandler()
	logger := slog.New(handler)

	return logger
}

func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}
