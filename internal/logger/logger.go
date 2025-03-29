package logger

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func Init(level LogLevel) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(convertLogLevel(level))
	cfg.DisableCaller = true
	rawLogger, err := cfg.Build()
	if err != nil {
		panic("unable to initialize logger: " + err.Error())
	}
	logger = rawLogger.Sugar()
}

func Debug(msg string, args ...any) {
	logger.Debugw(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Infow(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warnw(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Errorw(msg, args...)
}

func Fatal(msg string, args ...any) {
	logger.Fatalw(msg, args...)
}

func Sync() {
	if err := logger.Sync(); err != nil {
		logger.Errorw("Error during logger sync", "error", err)
	}
}

func init() {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(convertLogLevel(DEBUG))
	cfg.DisableCaller = true
	rawLogger, err := cfg.Build()
	if err != nil {
		panic("unable to initialize logger: " + err.Error())
	}
	logger = rawLogger.Sugar()
}
