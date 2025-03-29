package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

func convertLogLevel(level LogLevel) zapcore.Level {
	switch level {
	case DEBUG:
		return zap.DebugLevel
	case INFO:
		return zap.InfoLevel
	case WARN:
		return zap.WarnLevel
	case ERROR:
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
