package engines

import (
	"context"
)

type MemoryTransactionEngine struct{}

func NewMemoryTransactionEngine() *MemoryTransactionEngine {
	return &MemoryTransactionEngine{}
}

func (e *MemoryTransactionEngine) Begin(ctx context.Context) bool { return true }
func (e *MemoryTransactionEngine) Commit() bool                   { return true }
func (e *MemoryTransactionEngine) Rollback() bool                 { return true }
