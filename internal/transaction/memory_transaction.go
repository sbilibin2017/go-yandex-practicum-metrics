package transaction

import "context"

type MemoryTransactionEngine struct{}

func NewMemoryTransactionEngine() *MemoryTransactionEngine {
	return &MemoryTransactionEngine{}
}

func (dt *MemoryTransactionEngine) WithTransaction(
	ctx context.Context, fn func(ctx context.Context) (any, error),
) (any, error) {
	return fn(ctx)
}
