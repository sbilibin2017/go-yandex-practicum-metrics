package engines

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWithTransaction_Do_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTx := NewMockTransaction(ctrl)
	wt := NewWithTransaction(mockTx)
	mockTx.EXPECT().Begin(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit().Return(nil)
	f := func() error {
		return nil
	}
	success := wt.Do(context.Background(), f)
	assert.True(t, success)
}

func TestWithTransaction_Do_Panic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTx := NewMockTransaction(ctrl)
	wt := NewWithTransaction(mockTx)
	mockTx.EXPECT().Begin(gomock.Any()).Return(nil)
	mockTx.EXPECT().Rollback().Return(nil)
	mockTx.EXPECT().Commit().Times(0)
	f := func() error {
		panic("something went wrong")
	}
	assert.Panics(t, func() {
		wt.Do(context.Background(), f)
	})
}

func TestWithTransaction_Do_BeginError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTx := NewMockTransaction(ctrl)
	wt := NewWithTransaction(mockTx)
	mockTx.EXPECT().Begin(gomock.Any()).Return(errors.New("begin error"))
	f := func() error {
		return nil
	}
	success := wt.Do(context.Background(), f)
	assert.False(t, success)
}

func TestWithTransaction_Do_FuncError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTx := NewMockTransaction(ctrl)
	wt := NewWithTransaction(mockTx)
	mockTx.EXPECT().Begin(gomock.Any()).Return(nil)
	mockTx.EXPECT().Rollback().Return(nil)
	f := func() error {
		return errors.New("function error")
	}
	success := wt.Do(context.Background(), f)
	assert.False(t, success)
}

func TestWithTransaction_Do_RollbackError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTx := NewMockTransaction(ctrl)
	wt := NewWithTransaction(mockTx)
	mockTx.EXPECT().Begin(gomock.Any()).Return(nil)
	mockTx.EXPECT().Rollback().Return(errors.New("rollback error"))
	mockTx.EXPECT().Commit().Times(0)
	f := func() error {
		return errors.New("function error")
	}
	success := wt.Do(context.Background(), f)
	assert.False(t, success)
}

func TestWithTransaction_Do_CommitError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTx := NewMockTransaction(ctrl)
	wt := NewWithTransaction(mockTx)
	mockTx.EXPECT().Begin(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit().Return(errors.New("commit error"))
	f := func() error {
		return nil
	}
	success := wt.Do(context.Background(), f)
	assert.False(t, success)
}
