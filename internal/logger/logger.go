package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel int

const (
	InfoLevel LogLevel = iota
	ErrorLevel
)

var logger *zap.SugaredLogger

func InitializeLogger(level LogLevel) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(convertLogLevel(level)) // Конвертируем кастомный уровень

	rawLogger, _ := cfg.Build()

	logger = rawLogger.Sugar()

}

func convertLogLevel(level LogLevel) zapcore.Level {
	switch level {
	case InfoLevel:
		return zap.InfoLevel
	case ErrorLevel:
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}
