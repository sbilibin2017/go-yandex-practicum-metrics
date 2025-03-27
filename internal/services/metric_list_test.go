package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricListService_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	listRepo := NewMockMetricListRepository(ctrl)
	service := NewMetricListService(listRepo)
	metrics := map[domain.MetricID]*domain.Metrics{
		{ID: "metric1", Type: "counter"}: {
			ID:    "metric1",
			Type:  "counter",
			Delta: new(int64),
		},
	}
	listRepo.EXPECT().List(gomock.Any()).Return(metrics, nil)
	result, err := service.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, result[0].ID, "metric1")
	assert.Equal(t, result[0].Type, "counter")
}

func TestMetricListService_List_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	listRepo := NewMockMetricListRepository(ctrl)
	service := NewMetricListService(listRepo)
	listRepo.EXPECT().List(gomock.Any()).Return(nil, ErrMetricListInternal)
	result, err := service.List(context.Background())
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, err, ErrMetricListInternal)
}
