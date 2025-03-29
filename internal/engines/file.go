package engines

import (
	"go-yandex-practicum-metrics/internal/logger"
	"os"
)

type FileEngine struct {
	*os.File
}

type FileStoragePathGetter interface {
	GetFileStoragePath() string
}

func NewFileEngine(g FileStoragePathGetter) (*FileEngine, error) {
	file, err := os.Open(g.GetFileStoragePath())
	if err != nil {
		logger.Error("Error opening file", "error", err)
		return nil, err
	}
	return &FileEngine{File: file}, nil
}
