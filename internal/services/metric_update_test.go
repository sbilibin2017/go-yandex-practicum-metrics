package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricUpdateService_UpdateBatch_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	saveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	findRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	withTx := NewMockWithTransaction(ctrl)
	service := NewMetricUpdateService(saveRepo, findRepo, withTx)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  domain.Gauge,
			Value: new(float64),
		},
		{
			ID:    "metric2",
			Type:  domain.Counter,
			Delta: new(int64),
			Value: new(float64),
		},
	}
	existingMetrics := map[domain.MetricID]*domain.Metrics{
		{ID: "metric2", Type: domain.Counter}: {
			ID:    "metric2",
			Type:  domain.Counter,
			Delta: new(int64),
			Value: new(float64),
		},
	}
	withTx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func() error) error {
		return f()
	})
	findRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(existingMetrics, true)
	saveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(true)
	result, err := service.Update(context.Background(), metrics)
	assert.NoError(t, err)
	assert.Equal(t, metrics, result)
}

func TestMetricUpdateService_UpdateBatch_ErrorSaving(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	saveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	findRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	withTx := NewMockWithTransaction(ctrl)
	service := NewMetricUpdateService(saveRepo, findRepo, withTx)
	metrics := []*domain.Metrics{
		{
			ID: "metric1", Type: domain.Gauge,
			Value: new(float64),
		},
	}
	withTx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func() error) error {
		return f()
	})
	findRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metrics{
		{ID: "metric1", Type: domain.Gauge}: {
			ID: "metric1", Type: domain.Gauge,
			Value: new(float64),
		},
	}, true)
	saveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(false)
	result, err := service.Update(context.Background(), metrics)
	assert.Equal(t, errors.ErrInternal, err)
	assert.Nil(t, result)
}

func TestMetricUpdateService_UpdateBatch_ErrorFinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	saveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	findRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	withTx := NewMockWithTransaction(ctrl)
	service := NewMetricUpdateService(saveRepo, findRepo, withTx)
	metrics := []*domain.Metrics{
		{
			ID: "metric1", Type: domain.Gauge,
			Value: new(float64),
		},
	}
	withTx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f func() error) error {
		return f()
	})
	findRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, false)
	result, err := service.Update(context.Background(), metrics)
	assert.Equal(t, errors.ErrInternal, err)
	assert.Nil(t, result)
}
