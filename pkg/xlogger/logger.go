package xlogger

import "go.uber.org/zap"

// Logger is a logger that can be used to log messages.
// Deprecated: use https://pkg.go.dev/go.uber.org/zap@v1.24.0/zap instead.
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}
