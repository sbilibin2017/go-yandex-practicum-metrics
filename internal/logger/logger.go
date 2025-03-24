package logger

import (
	"fmt"
	"math/rand"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// InitializeLogger initializes the global logger based on the provided log level
func InitializeLogger(level zap.AtomicLevel) error {
	config := zap.NewProductionConfig()

	config.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "@timestamp",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeCaller: zapcore.FullCallerEncoder,
	}

	config.Level = level

	var err error
	logger, err = config.Build(zap.AddCaller()) // AddCaller() for file/line information
	return err
}

// logMessage is a helper function to log messages with additional fields
func logMessage(level zapcore.Level, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("trace_id", fmt.Sprintf("%x", rand.Int63())))

	switch level {
	case zap.DebugLevel:
		logger.Debug(msg, fields...)
	case zap.InfoLevel:
		logger.Info(msg, fields...)
	case zap.ErrorLevel:
		logger.Error(msg, fields...)
	}

}

// Debug logs a message with debug level and additional fields
func Debug(msg string, fields ...zap.Field) {
	logMessage(zap.DebugLevel, msg, fields...)
}

// Info logs a message with info level and additional fields
func Info(msg string, fields ...zap.Field) {
	logMessage(zap.InfoLevel, msg, fields...)
}

// Error logs a message with error level and additional fields
func Error(msg string, fields ...zap.Field) {
	logMessage(zap.ErrorLevel, msg, fields...)
}
