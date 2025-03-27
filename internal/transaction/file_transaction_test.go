package transaction

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileTransactionEngine_WithTransaction(t *testing.T) {
	ft := &FileTransactionEngine{}
	mockFn := func(ctx context.Context) (any, error) {
		return "success", nil
	}
	result, err := ft.WithTransaction(context.Background(), mockFn)
	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	mockFnWithError := func(ctx context.Context) (any, error) {
		return nil, assert.AnError
	}
	result, err = ft.WithTransaction(context.Background(), mockFnWithError)
	assert.Error(t, err)
	assert.Nil(t, result)
}
