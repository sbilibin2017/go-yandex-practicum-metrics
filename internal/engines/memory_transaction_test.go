package engines

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryTransactionEngine_BeginCommitRollback(t *testing.T) {
	engine := NewMemoryTransactionEngine()
	success := engine.Begin(context.Background())
	assert.True(t, success)
	success = engine.Commit()
	assert.True(t, success)
	success = engine.Rollback()
	assert.True(t, success)
}

func TestMemoryTransactionEngine_NilEngine(t *testing.T) {
	engine := NewMemoryTransactionEngine()
	success := engine.Begin(context.Background())
	assert.True(t, success)
	success = engine.Commit()
	assert.True(t, success)
	success = engine.Rollback()
	assert.True(t, success)
}
