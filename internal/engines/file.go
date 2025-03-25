package engines

import (
	"fmt"
	"os"
)

type FileConfig interface {
	GetFileStoragePath() string
}

type FileEngine struct {
	config FileConfig
	*os.File
}

// NewFileEngine is a constructor that initializes and returns a new FileEngine.
func NewFileEngine(config FileConfig) *FileEngine {
	return &FileEngine{config: config}
}

func (fe *FileEngine) Open() error {
	if fe.File != nil {
		return nil
	}
	file, err := os.OpenFile(fe.config.GetFileStoragePath(), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	fe.File = file
	return nil
}
