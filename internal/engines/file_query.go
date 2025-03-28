package engines

import (
	"bufio"
	"context"
	"os"
)

type FileQueryEngine struct {
	file *os.File
}

func NewFileQueryEngine(file *os.File) (*FileQueryEngine, error) {
	return &FileQueryEngine{
		file: file,
	}, nil
}

func (e *FileQueryEngine) Query(ctx context.Context, query string, args ...any) ([]any, bool) {
	if e.file == nil {
		return nil, false
	}
	var results []any
	scanner := bufio.NewScanner(e.file)
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, false
	}
	return results, true
}
