package configs

import (
	"strconv"
)

type ServerConfig struct {
	address         string
	databaseDSN     string
	storeInterval   string
	fileStoragePath string
	restore         string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

// Getters
func (s *ServerConfig) GetAddress() string {
	return s.address
}

func (s *ServerConfig) GetDatabaseDSN() string {
	return s.databaseDSN
}

func (s *ServerConfig) GetStoreInterval() int {
	interval, _ := strconv.Atoi(s.storeInterval)
	return interval
}

func (s *ServerConfig) GetFileStoragePath() string {
	return s.fileStoragePath
}

func (s *ServerConfig) GetRestore() bool {
	return s.restore == "true"
}

// Setters
func (s *ServerConfig) SetAddress(address string) bool {
	s.address = address
	return true
}

func (s *ServerConfig) SetDatabaseDSN(databaseDSN string) bool {
	s.databaseDSN = databaseDSN
	return true
}

func (s *ServerConfig) SetStoreInterval(storeInterval string) bool {
	if _, err := strconv.Atoi(storeInterval); err != nil {
		return false
	}
	s.storeInterval = storeInterval
	return true
}

func (s *ServerConfig) SetFileStoragePath(fileStoragePath string) bool {
	s.fileStoragePath = fileStoragePath
	return true
}

func (s *ServerConfig) SetRestore(restore string) bool {
	if restore != "true" && restore != "false" {
		return false
	}
	s.restore = restore
	return true
}
