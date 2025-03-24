package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerSetAddress(t *testing.T) {
	cfg := NewServerConfig()
	assert.True(t, cfg.SetAddress("localhost:8080"))
	assert.Equal(t, "localhost:8080", cfg.GetAddress())
}

func TestServerSetDatabaseDSN(t *testing.T) {
	cfg := NewServerConfig()
	assert.True(t, cfg.SetDatabaseDSN("user:password@/dbname"))
	assert.Equal(t, "user:password@/dbname", cfg.GetDatabaseDSN())
}

func TestServerSetStoreInterval(t *testing.T) {
	cfg := NewServerConfig()
	assert.True(t, cfg.SetStoreInterval("10"))
	assert.Equal(t, 10, cfg.GetStoreInterval())
}

func TestServerSetStoreIntervalInvalid(t *testing.T) {
	cfg := NewServerConfig()
	assert.False(t, cfg.SetStoreInterval("invalid"))
}

func TestServerSetFileStoragePath(t *testing.T) {
	cfg := NewServerConfig()
	assert.True(t, cfg.SetFileStoragePath("/path/to/storage"))
	assert.Equal(t, "/path/to/storage", cfg.GetFileStoragePath())
}

func TestServerSetRestore(t *testing.T) {
	cfg := NewServerConfig()
	assert.True(t, cfg.SetRestore("true"))
	assert.True(t, cfg.GetRestore())

	assert.True(t, cfg.SetRestore("false"))
	assert.False(t, cfg.GetRestore())
}

func TestServerSetRestoreInvalid(t *testing.T) {
	cfg := NewServerConfig()
	assert.False(t, cfg.SetRestore("invalid"))
}
