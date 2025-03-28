package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSave_SingleMetricSaved(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEngine := NewMockMemoryExecutorEngine(ctrl)
	mockEngine.EXPECT().Execute(context.Background(), "", domain.MetricID{ID: "1", Type: domain.Gauge}, &domain.Metrics{ID: "1", Type: domain.Gauge}).Return(true)
	repo := NewMetricMemorySaveBatchRepository(mockEngine)
	metrics := []*domain.Metrics{
		{ID: "1", Type: domain.Gauge},
	}
	ok := repo.Save(context.Background(), metrics)
	assert.True(t, ok)
}

func TestSave_MultipleMetricsSaved(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEngine := NewMockMemoryExecutorEngine(ctrl)
	mockEngine.EXPECT().Execute(context.Background(), "", domain.MetricID{ID: "1", Type: domain.Gauge}, &domain.Metrics{ID: "1", Type: domain.Gauge}).Return(true)
	mockEngine.EXPECT().Execute(context.Background(), "", domain.MetricID{ID: "2", Type: domain.Counter}, &domain.Metrics{ID: "2", Type: domain.Counter}).Return(true)
	repo := NewMetricMemorySaveBatchRepository(mockEngine)
	metrics := []*domain.Metrics{
		{ID: "1", Type: domain.Gauge},
		{ID: "2", Type: domain.Counter},
	}
	ok := repo.Save(context.Background(), metrics)
	assert.True(t, ok)
}

func TestSave_FailOnExecute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEngine := NewMockMemoryExecutorEngine(ctrl)
	mockEngine.EXPECT().Execute(context.Background(), "", domain.MetricID{ID: "1", Type: domain.Gauge}, &domain.Metrics{ID: "1", Type: domain.Gauge}).Return(false)
	repo := NewMetricMemorySaveBatchRepository(mockEngine)
	metrics := []*domain.Metrics{
		{ID: "1", Type: domain.Gauge},
	}
	ok := repo.Save(context.Background(), metrics)
	assert.False(t, ok)
}

func TestSave_EmptyMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEngine := NewMockMemoryExecutorEngine(ctrl)
	repo := NewMetricMemorySaveBatchRepository(mockEngine)
	metrics := []*domain.Metrics{}
	ok := repo.Save(context.Background(), metrics)
	assert.True(t, ok)
}
