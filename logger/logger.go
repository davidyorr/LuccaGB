package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	globalLogger *slog.Logger
	once         sync.Once
)

func Init(handler slog.Handler) {
	once.Do(func() {
		globalLogger = slog.New(handler)
	})
}

func logger() *slog.Logger {
	if globalLogger == nil {
		Init(slog.NewTextHandler(os.Stderr, nil))
	}
	return globalLogger
}

func Debug(msg string, args ...any) {
	logger().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logger().Error(msg, args...)
}
