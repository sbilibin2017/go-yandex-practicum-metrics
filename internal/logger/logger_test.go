package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger_DefaultValues(t *testing.T) {
	logger := NewLogger(Debug, nil)
	assert.Equal(t, Debug, logger.level)
	assert.IsType(t, &DefaultStringBuilder{}, logger.sb)
}

func TestLogger_Debug(t *testing.T) {
	logger := NewLogger(Debug, nil)
	message := "Test Debug message"
	log := logger.Debug(message, "user", "john_doe")
	assert.Contains(t, log, "\"level\": \"debug\"")
	assert.Contains(t, log, "\"message\": \"Test Debug message\"")
	assert.Contains(t, log, "\"user\": \"john_doe\"")
}

func TestLogger_Info(t *testing.T) {
	logger := NewLogger(Info, nil)
	message := "Test Info message"
	log := logger.Info(message, "status", "success")
	assert.Contains(t, log, "\"level\": \"info\"")
	assert.Contains(t, log, "\"message\": \"Test Info message\"")
	assert.Contains(t, log, "\"status\": \"success\"")
}

func TestLogger_Error(t *testing.T) {
	logger := NewLogger(Error, nil)
	message := "Test Error message"
	log := logger.Error(message, "error_code", "500")
	assert.Contains(t, log, "\"level\": \"error\"")
	assert.Contains(t, log, "\"message\": \"Test Error message\"")
	assert.Contains(t, log, "\"error_code\": \"500\"")
}

func TestLogger_LogLevelFilter(t *testing.T) {
	logger := NewLogger(Info, nil)
	message := "This should not show"
	log := logger.Debug(message)
	assert.Empty(t, log)
	log = logger.Info("This is info")
	assert.Contains(t, log, "\"level\": \"info\"")
	log = logger.Error("This is an error")
	assert.Contains(t, log, "\"level\": \"error\"")
}

func TestGenerateTraceID(t *testing.T) {
	traceID := generateTraceID()
	assert.Len(t, traceID, 16)
}

func TestLogger_EmptyExtraFields(t *testing.T) {
	logger := NewLogger(Debug, nil)
	message := "Test with no extra fields"
	log := logger.Debug(message)
	assert.NotContains(t, log, "\"key\":")
}

func TestLogger_ExtraFieldsWithOddLength(t *testing.T) {
	logger := NewLogger(Debug, nil)
	message := "Test with odd number of extra fields"
	log := logger.Debug(message, "key1", "value1", "key2")
	assert.Contains(t, log, "\"key1\": \"value1\"")
	assert.NotContains(t, log, "\"key2\":")
}

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{Debug, "debug"},
		{Info, "info"},
		{Error, "error"},
		{LogLevel(999), "debug"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.level), func(t *testing.T) {
			actual := tt.level.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
