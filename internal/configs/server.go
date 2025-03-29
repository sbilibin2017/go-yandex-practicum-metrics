package configs

import (
	"errors"
	"regexp"
	"strconv"
)

type ServerConfig struct {
	address         string
	databaseDSN     string
	storeInterval   string
	fileStoragePath string
	restore         string
}

func NewServerConfig(
	address string,
	databaseDSN string,
	storeInterval string,
	fileStoragePath string,
	restore string,
) *ServerConfig {
	return &ServerConfig{
		address:         address,
		databaseDSN:     databaseDSN,
		storeInterval:   storeInterval,
		fileStoragePath: fileStoragePath,
		restore:         restore,
	}
}

func (c *ServerConfig) SetAddress(address string) error {
	if match, _ := regexp.MatchString(`^\S+:\d+$`, address); !match {
		return errors.New("invalid address format, expected host:port")
	}
	c.address = address
	return nil
}

func (c *ServerConfig) SetDatabaseDSN(dsn string) error {
	pattern := `^(postgresql?|postgres):\/\/[^:]+:[^@]+@[^:]+:\d+\/[^\s]+$`
	if match, _ := regexp.MatchString(pattern, dsn); !match {
		return errors.New("invalid database DSN, must match 'postgres://user:password@host:port/dbname'")
	}
	c.databaseDSN = dsn
	return nil
}

func (c *ServerConfig) SetStoreInterval(interval string) error {
	if interval == "" {
		return nil
	}
	value, err := strconv.Atoi(interval)
	if err != nil || value <= 0 {
		return errors.New("store interval must be a positive integer")
	}
	c.storeInterval = interval
	return nil
}

func (c *ServerConfig) SetFileStoragePath(path string) error {
	if match, _ := regexp.MatchString(`^data/(.+/)*[^/]+\.json$`, path); !match {
		return errors.New("invalid file storage path, must start with 'data/' and end with '.json'")
	}
	c.fileStoragePath = path
	return nil
}

func (c *ServerConfig) SetRestore(restore string) error {
	if restore == "" {
		c.restore = "false"
		return nil
	}
	if restore != "true" && restore != "false" {
		return errors.New("restore must be 'true' or 'false'")
	}
	c.restore = restore
	return nil
}

func (c *ServerConfig) GetAddress() string {
	return c.address
}

func (c *ServerConfig) GetDatabaseDSN() string {
	return c.databaseDSN
}

func (c *ServerConfig) GetStoreInterval() string {
	return c.storeInterval
}

func (c *ServerConfig) GetFileStoragePath() string {
	return c.fileStoragePath
}

func (c *ServerConfig) GetRestore() string {
	return c.restore
}
