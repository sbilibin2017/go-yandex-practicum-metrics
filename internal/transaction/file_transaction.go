package transaction

import (
	"context"
)

type FileTransactionEngine struct{}

func (ft *FileTransactionEngine) WithTransaction(
	ctx context.Context, fn func(ctx context.Context) (any, error),
) (any, error) {
	return fn(ctx)
}
