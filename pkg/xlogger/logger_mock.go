package xlogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"testing"
)

type call struct {
	level  string
	msg    string
	fields []zap.Field
}
type LoggerMock struct {
	calls []call
}

func (l *LoggerMock) Log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	l.calls = append(l.calls, call{
		level:  lvl.String(),
		msg:    msg,
		fields: fields,
	})
}

func (l *LoggerMock) Info(msg string, fields ...zap.Field) {
	l.Log(zapcore.InfoLevel, msg, fields...)
}

func (l *LoggerMock) Warn(msg string, fields ...zap.Field) {
	l.Log(zapcore.WarnLevel, msg, fields...)
}

func (l *LoggerMock) Error(msg string, fields ...zap.Field) {
	l.Log(zapcore.ErrorLevel, msg, fields...)
}

func (l *LoggerMock) Fatal(msg string, fields ...zap.Field) {
	l.Log(zapcore.FatalLevel, msg, fields...)
}

// AssertLog asserts that the logger was called with the given log level, message and fields.
func (l *LoggerMock) AssertLog(t *testing.T, level zapcore.Level, msg string) {
	for _, call := range l.calls {
		if strings.Contains(call.msg, msg) {
			if call.level != level.String() {
				t.Errorf("expected log level %s with message %s, got %s", level, msg, call.level)
			}

			return
		}
	}

	switch len(l.calls) {
	case 0:
		t.Errorf("expected log level %s with message %s, got no logs", level, msg)
	case 1:
		t.Errorf("expected log level %s with message %s, got %s with message %s", level, msg, l.calls[0].level, l.calls[0].msg)
	default:
		t.Errorf("expected log with level %s and message %s, but no such log was found", level, msg)
	}
}
