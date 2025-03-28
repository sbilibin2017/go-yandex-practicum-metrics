package services

import (
	"context"
	"testing"

	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricUpdateService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSaveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	mockFindRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	mockWithTx := NewMockWithTx(ctrl)
	mockTx := NewMockTx(ctrl)
	service := NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockWithTx)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  domain.Counter,
			Delta: new(int64),
		},
	}
	metricID := domain.MetricID{ID: "metric1", Type: domain.Counter}
	existingMetrics := map[domain.MetricID]*domain.Metrics{
		metricID: {ID: "metric1", Type: domain.Counter, Delta: new(int64)},
	}
	mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(existingMetrics, true)
	mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(true)
	mockWithTx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx Tx) error) error {
		mockTx.EXPECT().Commit().Return(nil).Times(1)
		mockTx.EXPECT().Rollback().Return(nil).Times(0)
		err := fn(mockTx)
		mockTx.Commit()
		return err
	})
	updatedMetrics, err := service.Update(context.Background(), metrics)
	assert.NoError(t, err)
	assert.Len(t, updatedMetrics, 1)
	assert.Equal(t, *metrics[0].Delta, int64(0))
}

func TestMetricUpdateService_Update_FindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSaveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	mockFindRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	mockWithTx := NewMockWithTx(ctrl)
	mockTx := NewMockTx(ctrl)
	service := NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockWithTx)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  domain.Counter,
			Delta: new(int64),
		},
	}
	mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metrics, bool) {
		t.Log("Find was called with filters:", filters)
		return nil, false
	}).Times(1)
	mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
	mockWithTx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx Tx) error) error {
		t.Log("Inside Do function. Running the transaction function...")
		mockTx.EXPECT().Commit().Times(0)
		mockTx.EXPECT().Rollback().Times(0)
		err := fn(mockTx)
		t.Log("Transaction function executed. Error:", err)
		return err
	}).Times(1)
	updatedMetrics, err := service.Update(context.Background(), metrics)
	t.Log("Update result: Metrics:", updatedMetrics, "Error:", err)
	assert.Error(t, err)
	assert.Nil(t, updatedMetrics)
}

func TestMetricUpdateService_Update_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSaveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	mockFindRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	mockWithTx := NewMockWithTx(ctrl)
	mockTx := NewMockTx(ctrl)
	service := NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockWithTx)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  domain.Counter,
			Delta: new(int64),
		},
	}
	mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metrics{
		{
			ID:   "metric1",
			Type: domain.Counter,
		}: {
			ID:    "metric1",
			Type:  domain.Counter,
			Delta: new(int64),
		},
	}, true).Times(1)
	mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(false).Times(1)
	mockWithTx.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx Tx) error) error {
		t.Log("Inside Do function. Running the transaction function...")
		mockTx.EXPECT().Commit().Times(0)
		mockTx.EXPECT().Rollback().Times(0)
		err := fn(mockTx)
		t.Log("Transaction function executed. Error:", err)
		return err
	}).Times(1)
	updatedMetrics, err := service.Update(context.Background(), metrics)
	t.Log("Update result: Metrics:", updatedMetrics, "Error:", err)
	assert.Error(t, err)
	assert.Equal(t, errors.ErrInternal, err)
	assert.Nil(t, updatedMetrics)
}
