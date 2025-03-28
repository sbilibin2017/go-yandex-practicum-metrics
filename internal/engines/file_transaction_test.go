package engines

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileTransactionEngine_BeginCommitRollback(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile_*.txt")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	engine, err := NewFileTransactionEngine(tmpFile)
	if err != nil {
		t.Fatalf("failed to create FileTransactionEngine: %v", err)
	}
	success := engine.Begin(context.Background())
	assert.True(t, success)
	success = engine.Commit()
	assert.True(t, success)
	success = engine.Rollback()
	assert.True(t, success)
}

func TestFileTransactionEngine_NilFile(t *testing.T) {
	engine, err := NewFileTransactionEngine(nil)
	if err != nil {
		t.Fatalf("failed to create FileTransactionEngine with nil file: %v", err)
	}
	success := engine.Begin(context.Background())
	assert.True(t, success)
	success = engine.Commit()
	assert.True(t, success)
	success = engine.Rollback()
	assert.True(t, success)
}
