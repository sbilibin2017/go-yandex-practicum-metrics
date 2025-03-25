package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel — кастомный enum для уровней логирования
type LogLevel int

const (
	InfoLevel LogLevel = iota
	ErrorLevel
)

// logger — глобальный экземпляр логгера
var logger *zap.SugaredLogger

// InitializeLogger инициализирует глобальный логгер с кастомным уровнем
func InitializeLogger(level LogLevel) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(convertLogLevel(level)) // Конвертируем кастомный уровень

	rawLogger, _ := cfg.Build()

	logger = rawLogger.Sugar()

}

// convertLogLevel конвертирует кастомный LogLevel в zapcore.Level
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

// Infow логирует информационные сообщения
func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

// Errorw логирует ошибки
func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}
