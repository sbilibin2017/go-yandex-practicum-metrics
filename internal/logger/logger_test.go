package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitializeLogger(t *testing.T) {
	level := zap.NewAtomicLevelAt(zap.DebugLevel)
	err := InitializeLogger(level)
	assert.NoError(t, err, "Logger should initialize without error")
}

func TestGetLogger(t *testing.T) {
	level := zap.NewAtomicLevelAt(zap.DebugLevel)
	InitializeLogger(level)
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zap.DebugLevel,
	)

	logger = zap.New(core)

	Debug("debug message")
	Info("info message")
	Error("error message")

	logOutput := buf.String()
	assert.Contains(t, logOutput, "debug message", "Debug message should be logged")
	assert.Contains(t, logOutput, "info message", "Info message should be logged")
	assert.Contains(t, logOutput, "error message", "Error message should be logged")
}

func TestTraceIDPresence(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zap.DebugLevel,
	)

	logger = zap.New(core)

	Info("test message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err, "Log output should be valid JSON")
	_, exists := logEntry["trace_id"]
	assert.True(t, exists, "trace_id should be present in log entry")
}
