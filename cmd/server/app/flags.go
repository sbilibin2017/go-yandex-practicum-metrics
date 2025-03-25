package app

import (
	"flag"
	"go-yandex-practicum-metrics/internal/configs"

	"os"
)

const (
	// Значения по умолчанию
	DefaultAddress         = ":8080"
	DefaultDatabaseDSN     = ""
	DefaultStoreInterval   = ""
	DefaultFileStoragePath = ""
	DefaultRestore         = ""

	// Переменные окружения
	EnvAddress         = "ADDRESS"
	EnvDatabaseDSN     = "DATABASE_DSN"
	EnvStoreInterval   = "STORE_INTERVAL"
	EnvFileStoragePath = "FILE_STORAGE_PATH"
	EnvRestore         = "RESTORE"

	// Флаги командной строки
	FlagServerAddress   = "a"
	FlagDatabaseDSN     = "d"
	FlagStoreInterval   = "s"
	FlagFileStoragePath = "f"
	FlagRestore         = "r"

	// Описания флагов
	DescriptionServerAddress   = "Server address"
	DescriptionDatabaseDSN     = "Database DSN"
	DescriptionStoreInterval   = "Interval in seconds for data store"
	DescriptionFileStoragePath = "Path to file storage"
	DescriptionRestore         = "Restore backup (true/false)"
)

func ParseFlags() *configs.ServerConfig {
	const emptyString = ""
	var config configs.ServerConfig

	flag.StringVar(&config.Address, FlagServerAddress, DefaultAddress, DescriptionServerAddress)
	flag.StringVar(&config.DatabaseDSN, FlagDatabaseDSN, DefaultDatabaseDSN, DescriptionDatabaseDSN)
	flag.StringVar(&config.StoreInterval, FlagStoreInterval, DefaultStoreInterval, DescriptionStoreInterval)
	flag.StringVar(&config.FileStoragePath, FlagFileStoragePath, DefaultFileStoragePath, DescriptionFileStoragePath)
	flag.StringVar(&config.Restore, FlagRestore, DefaultRestore, DescriptionRestore)

	flag.Parse()

	if envAddress := os.Getenv(EnvAddress); envAddress != emptyString {
		config.Address = envAddress
	}
	if envDatabaseDSN := os.Getenv(EnvDatabaseDSN); envDatabaseDSN != emptyString {
		config.DatabaseDSN = envDatabaseDSN
	}
	if envStoreInterval := os.Getenv(EnvStoreInterval); envStoreInterval != emptyString {
		config.StoreInterval = envStoreInterval
	}
	if envFileStoragePath := os.Getenv(EnvFileStoragePath); envFileStoragePath != emptyString {
		config.FileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv(EnvRestore); envRestore != emptyString {
		config.Restore = envRestore
	}

	return &config

}
