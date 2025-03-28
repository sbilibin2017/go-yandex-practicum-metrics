package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSave_SingleMetricSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEngine := NewMockFileExecutorEngine(ctrl)
	mockEngine.EXPECT().Execute(context.Background(), "", domain.MetricID{ID: "1", Type: domain.Gauge}, &domain.Metrics{ID: "1", Type: domain.Gauge}).Return(true)
	repo := NewMetricFileSaveBatchRepository(mockEngine)
	metrics := []*domain.Metrics{
		{ID: "1", Type: domain.Gauge},
	}
	result := repo.Save(context.Background(), metrics)
	assert.True(t, result)
}

func TestSave_MultipleMetricsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEngine := NewMockFileExecutorEngine(ctrl)
	mockEngine.EXPECT().Execute(context.Background(), "", domain.MetricID{ID: "1", Type: domain.Gauge}, &domain.Metrics{ID: "1", Type: domain.Gauge}).Return(true)
	mockEngine.EXPECT().Execute(context.Background(), "", domain.MetricID{ID: "2", Type: domain.Counter}, &domain.Metrics{ID: "2", Type: domain.Counter}).Return(true)
	repo := NewMetricFileSaveBatchRepository(mockEngine)
	metrics := []*domain.Metrics{
		{ID: "1", Type: domain.Gauge},
		{ID: "2", Type: domain.Counter},
	}
	result := repo.Save(context.Background(), metrics)
	assert.True(t, result)
}

func TestSave_MetricExecutionFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEngine := NewMockFileExecutorEngine(ctrl)
	mockEngine.EXPECT().Execute(context.Background(), "", domain.MetricID{ID: "1", Type: domain.Gauge}, &domain.Metrics{ID: "1", Type: domain.Gauge}).Return(false)
	repo := NewMetricFileSaveBatchRepository(mockEngine)
	metrics := []*domain.Metrics{
		{ID: "1", Type: domain.Gauge},
	}
	result := repo.Save(context.Background(), metrics)
	assert.False(t, result)
}

func TestSave_MultipleMetricsWithFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEngine := NewMockFileExecutorEngine(ctrl)
	mockEngine.EXPECT().Execute(context.Background(), "", domain.MetricID{ID: "1", Type: domain.Gauge}, &domain.Metrics{ID: "1", Type: domain.Gauge}).Return(true)
	mockEngine.EXPECT().Execute(context.Background(), "", domain.MetricID{ID: "2", Type: domain.Counter}, &domain.Metrics{ID: "2", Type: domain.Counter}).Return(false)
	repo := NewMetricFileSaveBatchRepository(mockEngine)
	metrics := []*domain.Metrics{
		{ID: "1", Type: domain.Gauge},
		{ID: "2", Type: domain.Counter},
	}
	result := repo.Save(context.Background(), metrics)
	assert.False(t, result)
}
