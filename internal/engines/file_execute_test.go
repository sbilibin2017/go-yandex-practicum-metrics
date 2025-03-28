package engines

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExecuteEngine_Execute(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile_*.txt")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	engine, err := NewFileExecuteEngine(tmpFile)
	if err != nil {
		t.Fatalf("failed to create FileExecuteEngine: %v", err)
	}
	query := "Hello, %s!"
	args := []any{"World"}
	success := engine.Execute(context.Background(), query, args...)
	assert.True(t, success)
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read the temporary file: %v", err)
	}
	expectedContent := "Hello, World!\n"
	assert.Equal(t, expectedContent, string(content))
	success = engine.Execute(context.Background(), "Another test: %d", 123)
	assert.True(t, success)
	content, err = os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read the temporary file: %v", err)
	}
	expectedContent = "Hello, World!\nAnother test: 123\n"
	assert.Equal(t, expectedContent, string(content))
}

func TestFileExecuteEngine_NilFile(t *testing.T) {
	engine, err := NewFileExecuteEngine(nil)
	if err != nil {
		t.Fatalf("failed to create FileExecuteEngine with nil file: %v", err)
	}
	success := engine.Execute(context.Background(), "Query without file")
	assert.False(t, success)
}
