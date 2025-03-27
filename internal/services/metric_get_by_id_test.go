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

func TestMetricGetService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	findRepo := NewMockMetricGetFindBatchRepository(ctrl)
	service := NewMetricGetService(findRepo)
	metricID := domain.MetricID{ID: "metric1", Type: "counter"}
	metrics := &domain.Metrics{
		ID:    "metric1",
		Type:  "counter",
		Delta: new(int64),
	}
	findRepo.EXPECT().FindBatch(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metrics{
		metricID: metrics,
	}, nil)
	result, err := service.GetByID(context.Background(), metricID)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, result[0].ID, "metric1")
	assert.Equal(t, result[0].Type, "counter")
}

func TestMetricGetService_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	findRepo := NewMockMetricGetFindBatchRepository(ctrl)
	service := NewMetricGetService(findRepo)
	metricID := domain.MetricID{ID: "metric1", Type: "counter"}
	findRepo.EXPECT().FindBatch(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metrics{}, nil)
	result, err := service.GetByID(context.Background(), metricID)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, err, ErrMetricNotFound)
}

func TestMetricGetService_GetByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	findRepo := NewMockMetricGetFindBatchRepository(ctrl)
	service := NewMetricGetService(findRepo)
	metricID := domain.MetricID{ID: "metric1", Type: "counter"}
	findRepo.EXPECT().FindBatch(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))
	result, err := service.GetByID(context.Background(), metricID)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, err, ErrMetricGetInternal)
}
