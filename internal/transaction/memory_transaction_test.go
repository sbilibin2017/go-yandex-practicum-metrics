package transaction

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryTransactionEngine_WithTransaction(t *testing.T) {
	testCases := []struct {
		name    string
		fn      func(ctx context.Context) (any, error)
		expects any
		wantErr bool
	}{
		{
			name: "Test WithTransaction - Successful Execution",
			fn: func(ctx context.Context) (any, error) {
				return "success", nil
			},
			expects: "success",
			wantErr: false,
		},
		{
			name: "Test WithTransaction - Function Returns Error",
			fn: func(ctx context.Context) (any, error) {
				return nil, errors.New("transaction failed")
			},
			expects: nil,
			wantErr: true,
		},
		{
			name: "Test WithTransaction - Returns Integer",
			fn: func(ctx context.Context) (any, error) {
				return 42, nil
			},
			expects: 42,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			engine := NewMemoryTransactionEngine()
			ctx := context.Background()

			result, err := engine.WithTransaction(ctx, tc.fn)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expects, result)
			}
		})
	}
}
