package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricListService_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	listRepo := NewMockMetricListRepository(ctrl)
	service := NewMetricListService(listRepo)
	metric1 := &domain.Metrics{
		ID: "metric1", Type: domain.Gauge,
		Value: new(float64),
	}
	metric2 := &domain.Metrics{
		ID: "metric2", Type: domain.Counter,
		Value: new(float64),
	}
	metric1ID := domain.MetricID{ID: metric1.ID, Type: metric1.Type}
	metric2ID := domain.MetricID{ID: metric2.ID, Type: metric2.Type}
	listRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metrics{
		metric1ID: metric1,
		metric2ID: metric2,
	}, true)
	result, err := service.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, metric1)
	assert.Contains(t, result, metric2)
}

func TestMetricListService_List_ErrorFinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	listRepo := NewMockMetricListRepository(ctrl)
	service := NewMetricListService(listRepo)
	listRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, false)
	result, err := service.List(context.Background())
	assert.Equal(t, errors.ErrInternal, err)
	assert.Nil(t, result)
}
