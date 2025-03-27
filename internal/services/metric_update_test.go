package services

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdateService_UpdateBatch_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	saveRepo := NewMockMetricSaveBatchRepository(ctrl)
	findRepo := NewMockMetricFindBatchRepository(ctrl)
	tx := NewMockTransaction(ctrl)
	service := NewMetricUpdateService(saveRepo, findRepo, tx)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Delta: new(int64),
		},
		{
			ID:    "metric2",
			Type:  "gauge",
			Value: new(float64),
		},
	}
	existingMetrics := map[domain.MetricID]*domain.Metrics{
		{ID: "metric1", Type: "counter"}: {ID: "metric1", Type: "counter", Delta: new(int64)},
	}
	findRepo.EXPECT().FindBatch(gomock.Any(), gomock.Any()).Return(existingMetrics, nil)
	saveRepo.EXPECT().SaveBatch(gomock.Any(), gomock.Any()).Return(nil)
	tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
		return fn(ctx)
	})
	updatedMetrics, err := service.UpdateBatch(context.Background(), metrics)
	require.NoError(t, err)
	assert.Len(t, updatedMetrics, 2)
	assert.Equal(t, *updatedMetrics[0].Delta, int64(0))
	assert.Equal(t, *updatedMetrics[1].Value, 0.0)
}

func TestMetricUpdateService_UpdateBatch_AddNewCounterMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	saveRepo := NewMockMetricSaveBatchRepository(ctrl)
	findRepo := NewMockMetricFindBatchRepository(ctrl)
	tx := NewMockTransaction(ctrl)
	service := NewMetricUpdateService(saveRepo, findRepo, tx)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Delta: new(int64),
		},
	}
	findRepo.EXPECT().FindBatch(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metrics{}, nil)
	saveRepo.EXPECT().SaveBatch(gomock.Any(), gomock.Any()).Return(nil)
	tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
			return fn(ctx)
		},
	).Times(1)
	updatedMetrics, err := service.UpdateBatch(context.Background(), metrics)
	require.NoError(t, err)
	assert.Len(t, updatedMetrics, 1)
	assert.Equal(t, updatedMetrics[0].ID, "metric1")
	assert.Equal(t, updatedMetrics[0].Type, "counter")
}

func TestMetricUpdateService_UpdateBatch_FindBatchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	saveRepo := NewMockMetricSaveBatchRepository(ctrl)
	findRepo := NewMockMetricFindBatchRepository(ctrl)
	tx := NewMockTransaction(ctrl)
	service := NewMetricUpdateService(saveRepo, findRepo, tx)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Delta: new(int64),
		},
	}
	findRepo.EXPECT().FindBatch(gomock.Any(), gomock.Any()).Return(nil, errors.New("find error"))
	tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
			return fn(ctx)
		},
	).Times(1)
	updatedMetrics, err := service.UpdateBatch(context.Background(), metrics)
	require.Error(t, err)
	assert.Nil(t, updatedMetrics)
	assert.Equal(t, ErrMetricUpdateInternal, err)
}

func TestMetricUpdateService_UpdateBatch_SaveBatchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	saveRepo := NewMockMetricSaveBatchRepository(ctrl)
	findRepo := NewMockMetricFindBatchRepository(ctrl)
	tx := NewMockTransaction(ctrl)
	service := NewMetricUpdateService(saveRepo, findRepo, tx)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Delta: new(int64),
		},
	}
	existingMetrics := map[domain.MetricID]*domain.Metrics{
		{ID: "metric1", Type: "counter"}: {ID: "metric1", Type: "counter", Delta: new(int64)},
	}
	findRepo.EXPECT().FindBatch(gomock.Any(), gomock.Any()).Return(existingMetrics, nil)
	saveRepo.EXPECT().SaveBatch(gomock.Any(), gomock.Any()).Return(errors.New("save error"))
	tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
			return fn(ctx)
		},
	).Times(1)
	updatedMetrics, err := service.UpdateBatch(context.Background(), metrics)
	require.Error(t, err)
	assert.Nil(t, updatedMetrics)
	assert.Equal(t, ErrMetricUpdateInternal, err)
}

func TestMetricUpdateService_UpdateBatch_TransactionFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	saveRepo := NewMockMetricSaveBatchRepository(ctrl)
	findRepo := NewMockMetricFindBatchRepository(ctrl)
	tx := NewMockTransaction(ctrl)
	service := NewMetricUpdateService(saveRepo, findRepo, tx)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Delta: new(int64),
		},
	}
	tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
			return nil, errors.New("transaction failed")
		},
	).Times(1)
	updatedMetrics, err := service.UpdateBatch(context.Background(), metrics)
	require.Error(t, err)
	assert.Nil(t, updatedMetrics)
	assert.Equal(t, ErrMetricUpdateInternal, err)
}

func TestMetricUpdateService_UpdateBatch_UnexpectedResultType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock dependencies
	saveRepo := NewMockMetricSaveBatchRepository(ctrl)
	findRepo := NewMockMetricFindBatchRepository(ctrl)
	tx := NewMockTransaction(ctrl)

	// Create MetricUpdateService with mocked dependencies
	service := NewMetricUpdateService(saveRepo, findRepo, tx)

	// Input test data
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Delta: new(int64),
		},
	}

	// Mocking the transaction to return an invalid type
	tx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) (any, error)) (any, error) {
			// Return an incorrect type (string instead of []*domain.Metrics)
			return "invalid result", nil
		},
	).Times(1)

	// Call the UpdateBatch method
	updatedMetrics, err := service.UpdateBatch(context.Background(), metrics)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, updatedMetrics)
	assert.Equal(t, ErrMetricUpdateInternal, err) // Ensure we check for the correct error
}
