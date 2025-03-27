package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeLogger(t *testing.T) {
	Init(InfoLevel)
	assert.NotNil(t, Logger, "Logger should be initialized")
}
