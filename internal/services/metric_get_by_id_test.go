package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricGetService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	findRepo := NewMockMetricGetFindBatchRepository(ctrl)
	service := NewMetricGetService(findRepo)
	metrics := &domain.Metrics{
		ID: "metric1", Type: domain.Gauge,
		Value: new(float64),
	}
	metricID := domain.MetricID{ID: metrics.ID, Type: metrics.Type}
	findRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metrics{
		metricID: metrics,
	}, true)
	result, err := service.GetByID(context.Background(), &metricID)
	assert.NoError(t, err)
	assert.Equal(t, metrics, result)
}

func TestMetricGetService_GetByID_MetricNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	findRepo := NewMockMetricGetFindBatchRepository(ctrl)
	service := NewMetricGetService(findRepo)
	metricID := &domain.MetricID{ID: "metric1", Type: domain.Gauge}
	findRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metrics{}, true)
	result, err := service.GetByID(context.Background(), metricID)
	assert.Equal(t, errors.ErrMetricNotFound, err)
	assert.Nil(t, result)
}

func TestMetricGetService_GetByID_ErrorFinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	findRepo := NewMockMetricGetFindBatchRepository(ctrl)
	service := NewMetricGetService(findRepo)
	metricID := &domain.MetricID{ID: "metric1", Type: domain.Gauge}
	findRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, false)
	result, err := service.GetByID(context.Background(), metricID)
	assert.Equal(t, errors.ErrInternal, err)
	assert.Nil(t, result)
}
