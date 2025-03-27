package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestConvertLogLevel(t *testing.T) {
	testCases := []struct {
		name     string
		input    LogLevel
		expected zapcore.Level
	}{
		{
			name:     "DebugLevel",
			input:    DebugLevel,
			expected: zapcore.DebugLevel,
		},
		{
			name:     "InfoLevel",
			input:    InfoLevel,
			expected: zapcore.InfoLevel,
		},
		{
			name:     "ErrorLevel",
			input:    ErrorLevel,
			expected: zapcore.ErrorLevel,
		},
		{
			name:     "DefaultCase",
			input:    LogLevel(100),     // Неподдерживаемое значение
			expected: zapcore.InfoLevel, // Ожидаем значение по умолчанию
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := convertLogLevel(tc.input)
			assert.Equal(t, tc.expected, actual, "Unexpected log level conversion result")
		})
	}
}
