package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	globalLogger *slog.Logger
	once         sync.Once
	LogLevel     = "error"
)

func Init(handler slog.Handler) {
	once.Do(func() {
		globalLogger = slog.New(handler)
	})
}

func logger() *slog.Logger {
	once.Do(func() {
		var level slog.Level
		switch LogLevel {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelError
		}
		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
		globalLogger = slog.New(handler)
	})
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
