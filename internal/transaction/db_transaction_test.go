package transaction

import (
	"context"
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDBTransactionEngine_WithTransaction_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := NewMockDB(ctrl)
	mockTx := NewMockTx(ctrl)
	mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(mockTx, nil)
	mockTx.EXPECT().Commit().Return(nil)
	engine := &DBTransactionEngine{db: mockDB}
	err := engine.WithTransaction(context.Background(), func(ctx context.Context) error {
		tx := ctx.Value(dbTransactionKey).(Tx)
		assert.NotNil(t, tx)
		return nil
	})
	assert.NoError(t, err)
}

func TestDBTransactionEngine_WithTransaction_BeginTxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := NewMockDB(ctrl)
	mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(nil, fmt.Errorf("begin transaction error"))
	engine := &DBTransactionEngine{db: mockDB}
	err := engine.WithTransaction(context.Background(), func(ctx context.Context) error {
		return nil
	})
	assert.Error(t, err)
	assert.Equal(t, "begin transaction error", err.Error())
}

func TestDBTransactionEngine_WithTransaction_Panic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := NewMockDB(ctrl)
	mockTx := NewMockTx(ctrl)
	mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)
	engine := &DBTransactionEngine{db: mockDB}
	assert.Panics(t, func() {
		engine.WithTransaction(context.Background(), func(ctx context.Context) error {
			panic("panic in fn")
		})
	})
}

func TestDBTransactionEngine_WithTransaction_RollbackOnFnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := NewMockDB(ctrl)
	mockTx := NewMockTx(ctrl)
	mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback().Return(nil)
	engine := &DBTransactionEngine{db: mockDB}
	err := engine.WithTransaction(context.Background(), func(ctx context.Context) error {
		return fmt.Errorf("function error")
	})
	assert.Error(t, err)
	assert.Equal(t, "function error", err.Error())
}
