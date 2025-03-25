package configs

type ServerConfig struct {
	Address         string
	DatabaseDSN     string
	StoreInterval   string
	FileStoragePath string
	Restore         string
}

// Конструктор для создания нового экземпляра ServerConfig
func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

// Геттер для Address
func (c *ServerConfig) GetAddress() string {
	return c.Address
}

// Геттер для DatabaseDSN
func (c *ServerConfig) GetDatabaseDSN() string {
	return c.DatabaseDSN
}

// Геттер для StoreInterval
func (c *ServerConfig) GetStoreInterval() string {
	return c.StoreInterval
}

// Геттер для FileStoragePath
func (c *ServerConfig) GetFileStoragePath() string {
	return c.FileStoragePath
}

// Геттер для Restore
func (c *ServerConfig) GetRestore() string {
	return c.Restore
}
