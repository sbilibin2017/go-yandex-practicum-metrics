package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitializeLoggerWithInfoLevel(t *testing.T) {
	assert.NotPanics(t, func() {
		InitializeLogger(InfoLevel)
	})
}

func TestInitializeLoggerWithErrorLevel(t *testing.T) {
	assert.NotPanics(t, func() {
		InitializeLogger(ErrorLevel)
	})
}

func TestInfowLogging(t *testing.T) {
	// Initializing logger
	InitializeLogger(InfoLevel)

	assert.NotPanics(t, func() {
		Infow("Test info message", "key", "value")
	})
}

func TestErrorwLogging(t *testing.T) {
	// Initializing logger
	InitializeLogger(InfoLevel)

	assert.NotPanics(t, func() {
		Errorw("Test error message", "key", "value")
	})
}

func TestConvertLogLevel(t *testing.T) {
	testCases := []struct {
		input    LogLevel
		expected zapcore.Level
	}{
		{InfoLevel, zap.InfoLevel},
		{ErrorLevel, zap.ErrorLevel},
		{LogLevel(99), zap.InfoLevel}, // Проверка значения по умолчанию
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, convertLogLevel(tc.input))
	}
}
