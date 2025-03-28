package transaction

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTx реализует интерфейс Tx
type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

// MockDB реализует интерфейс DB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin(ctx context.Context) (Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(Tx), args.Error(1)
}

func TestWithTx_Do_Success(t *testing.T) {
	mockDB := new(MockDB)
	mockTx := new(MockTx)
	ctx := context.Background()

	mockDB.On("Begin", ctx).Return(mockTx, nil)
	mockTx.On("Commit").Return(nil)

	withTx := NewWithTx(mockDB)

	err := withTx.Do(ctx, func(tx Tx) error {
		assert.NotNil(t, tx)
		return nil
	})

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestWithTx_Do_FuncError(t *testing.T) {
	mockDB := new(MockDB)
	mockTx := new(MockTx)
	ctx := context.Background()

	mockDB.On("Begin", ctx).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)

	withTx := NewWithTx(mockDB)

	expectedErr := errors.New("function error")
	err := withTx.Do(ctx, func(tx Tx) error {
		assert.NotNil(t, tx)
		return expectedErr
	})

	assert.ErrorIs(t, err, expectedErr)
	mockDB.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestWithTx_Do_CommitError(t *testing.T) {
	mockDB := new(MockDB)
	mockTx := new(MockTx)
	ctx := context.Background()

	mockDB.On("Begin", ctx).Return(mockTx, nil)
	mockTx.On("Commit").Return(errors.New("commit error"))

	withTx := NewWithTx(mockDB)

	err := withTx.Do(ctx, func(tx Tx) error {
		assert.NotNil(t, tx)
		return nil
	})

	assert.EqualError(t, err, "commit error")
	mockDB.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestWithTx_Do_BeginError(t *testing.T) {
	mockDB := new(MockDB)
	ctx := context.Background()

	mockDB.On("Begin", ctx).Return((*MockTx)(nil), errors.New("begin error"))

	withTx := NewWithTx(mockDB)

	err := withTx.Do(ctx, func(tx Tx) error {
		return nil
	})

	assert.EqualError(t, err, "begin error")
	mockDB.AssertExpectations(t)
}

func TestWithTx_Do_NoDB(t *testing.T) {
	withTx := NewWithTx(nil)
	ctx := context.Background()

	err := withTx.Do(ctx, func(tx Tx) error {
		assert.Nil(t, tx)
		return nil
	})

	assert.NoError(t, err)
}
