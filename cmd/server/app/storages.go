package app

import (
	"go-yandex-practicum-metrics/internal/types"
)

type DBConfig interface {
	GetDatabaseDSN() string
}

type FileConfig interface {
	GetFileStoragePath() string
}

type Config interface {
	DBConfig
	FileConfig
}

type DB interface {
	Open() error
}

type File interface {
	Open() error
}

type Storages struct {
	DB     DB
	File   File
	Memory map[types.MetricID]*types.Metrics
}

func NewStorages(config Config, db DB, file File) (
	Storages,
	error,
) {
	var storages Storages

	if config.GetDatabaseDSN() != "" {
		if err := db.Open(); err != nil {
			return storages, err
		}
		storages.DB = db
	}

	if config.GetFileStoragePath() != "" {
		if err := file.Open(); err != nil {
			return storages, err
		}
		storages.File = file
	}

	if config.GetDatabaseDSN() == "" && config.GetFileStoragePath() == "" {
		storages.Memory = make(map[types.MetricID]*types.Metrics)
	}

	return storages, nil
}
