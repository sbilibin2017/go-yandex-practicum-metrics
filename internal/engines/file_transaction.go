package engines

import (
	"context"
	"os"
)

type FileTransactionEngine struct {
	*os.File
}

func NewFileTransactionEngine(file *os.File) (*FileTransactionEngine, error) {
	return &FileTransactionEngine{
		File: file,
	}, nil
}

func (e *FileTransactionEngine) Begin(ctx context.Context) bool { return true }
func (e *FileTransactionEngine) Commit() bool                   { return true }
func (e *FileTransactionEngine) Rollback() bool                 { return true }
