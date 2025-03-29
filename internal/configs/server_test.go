package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerConfig_SetAddress(t *testing.T) {
	cfg := &ServerConfig{}
	testCases := []struct {
		name    string
		input   string
		expects string
		wantErr bool
	}{
		{"valid address", "localhost:8080", "localhost:8080", false},
		{"invalid address", "localhost", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cfg.SetAddress(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expects, cfg.GetAddress())
			}
		})
	}
}

func TestServerConfig_SetDatabaseDSN(t *testing.T) {
	cfg := &ServerConfig{}
	testCases := []struct {
		name    string
		input   string
		expects string
		wantErr bool
	}{
		{"valid DSN", "postgres://user:password@localhost:5432/dbname", "postgres://user:password@localhost:5432/dbname", false},
		{"invalid DSN", "invalid_dsn", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cfg.SetDatabaseDSN(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expects, cfg.GetDatabaseDSN())
			}
		})
	}
}

func TestServerConfig_SetStoreInterval(t *testing.T) {
	cfg := &ServerConfig{}
	testCases := []struct {
		name    string
		input   string
		expects string
		wantErr bool
	}{
		{"valid interval", "10", "10", false},
		{"invalid interval", "invalid", "", true},
		{"negative interval", "-5", "", true},
		{"empty interval", "", "", false}, // новый тест для пустой строки
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Сначала очищаем значение для чистоты теста
			cfg.storeInterval = ""

			err := cfg.SetStoreInterval(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expects, cfg.GetStoreInterval())
			}
		})
	}
}

func TestServerConfig_SetFileStoragePath(t *testing.T) {
	cfg := &ServerConfig{}
	testCases := []struct {
		name    string
		input   string
		expects string
		wantErr bool
	}{
		{"valid path", "data/storage.json", "data/storage.json", false},
		{"invalid path", "invalid_path.json", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cfg.SetFileStoragePath(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expects, cfg.GetFileStoragePath())
			}
		})
	}
}

func TestServerConfig_SetRestore(t *testing.T) {
	cfg := &ServerConfig{}
	testCases := []struct {
		name    string
		input   string
		expects string
		wantErr bool
	}{
		{"empty restore", "", "false", false},
		{"valid restore true", "true", "true", false},
		{"valid restore false", "false", "false", false},
		{"invalid restore", "maybe", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cfg.SetRestore(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expects, cfg.GetRestore())
			}
		})
	}
}
