package engines

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileQueryEngine_Query(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile_*.txt")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.WriteString("Line 1\nLine 2\nLine 3\n")
	if err != nil {
		t.Fatalf("failed to write to the temporary file: %v", err)
	}
	tmpFile.Seek(0, 0) // Ensure the file pointer is at the beginning
	engine, err := NewFileQueryEngine(tmpFile)
	if err != nil {
		t.Fatalf("failed to create FileQueryEngine: %v", err)
	}
	results, success := engine.Query(context.Background(), "", nil...)
	assert.True(t, success)
	assert.Len(t, results, 3)
	assert.Equal(t, "Line 1", results[0])
	assert.Equal(t, "Line 2", results[1])
	assert.Equal(t, "Line 3", results[2])
}

func TestFileQueryEngine_NilFile(t *testing.T) {
	engine, err := NewFileQueryEngine(nil)
	if err != nil {
		t.Fatalf("failed to create FileQueryEngine with nil file: %v", err)
	}
	results, success := engine.Query(context.Background(), "", nil...)
	assert.False(t, success)
	assert.Nil(t, results)
}
