package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func Init(level LogLevel) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(convertLogLevel(level))
	rawLogger, _ := cfg.Build()
	Logger = rawLogger.Sugar()
}
