package configs

import (
	"strconv"
)

type AgentConfig struct {
	address        string
	reportInterval string
	pollInterval   string
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{}
}

// Getters
func (a *AgentConfig) GetAddress() string {
	return a.address
}

func (a *AgentConfig) GetReportInterval() int {
	interval, _ := strconv.Atoi(a.reportInterval)
	return interval
}

func (a *AgentConfig) GetPollInterval() int {
	interval, _ := strconv.Atoi(a.pollInterval)
	return interval
}

// Setters
func (a *AgentConfig) SetAddress(address string) bool {
	a.address = address
	return true
}

func (a *AgentConfig) SetReportInterval(reportInterval string) bool {
	if _, err := strconv.Atoi(reportInterval); err != nil {
		return false
	}
	a.reportInterval = reportInterval
	return true
}

func (a *AgentConfig) SetPollInterval(pollInterval string) bool {
	if _, err := strconv.Atoi(pollInterval); err != nil {
		return false
	}
	a.pollInterval = pollInterval
	return true
}
