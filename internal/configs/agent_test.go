package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentSetAddress(t *testing.T) {
	cfg := NewAgentConfig()
	assert.True(t, cfg.SetAddress("localhost:8080"))
	assert.Equal(t, "localhost:8080", cfg.GetAddress())
}

func TestAgentSetReportInterval(t *testing.T) {
	cfg := NewAgentConfig()
	assert.True(t, cfg.SetReportInterval("10"))
	assert.Equal(t, 10, cfg.GetReportInterval())
}

func TestAgentSetReportIntervalInvalid(t *testing.T) {
	cfg := NewAgentConfig()
	assert.False(t, cfg.SetReportInterval("invalid"))
}

func TestAgentSetPollInterval(t *testing.T) {
	cfg := NewAgentConfig()
	assert.True(t, cfg.SetPollInterval("5"))
	assert.Equal(t, 5, cfg.GetPollInterval())
}

func TestAgentSetPollIntervalInvalid(t *testing.T) {
	cfg := NewAgentConfig()
	assert.False(t, cfg.SetPollInterval("invalid"))
}
