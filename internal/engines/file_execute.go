package engines

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

type FileExecuteEngine struct {
	file *os.File
}

func NewFileExecuteEngine(file *os.File) (*FileExecuteEngine, error) {
	return &FileExecuteEngine{
		file: file,
	}, nil
}

func (e *FileExecuteEngine) Execute(ctx context.Context, query string, args ...any) bool {
	if e.file == nil {
		return false
	}
	line := fmt.Sprintf(query, args...)
	writer := bufio.NewWriter(e.file)
	_, err := writer.WriteString(line + "\n")
	if err != nil {
		return false
	}
	writer.Flush()
	return true
}
