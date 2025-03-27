package usecases

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdateBodyUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := NewMockMetricUpdateBodyService(ctrl)
	metric := &domain.Metrics{
		ID:    "metric1",
		Type:  string(domain.Counter),
		Delta: ptrInt64(100),
		Value: nil,
	}
	mockSvc.EXPECT().
		UpdateBatch(context.Background(), []*domain.Metrics{metric}).
		Return([]*domain.Metrics{metric}, nil).
		Times(1)
	usecase := NewMetricUpdateBodyUsecase(mockSvc)
	req := &MetricUpdateBodyRequest{
		Metrics: metric,
	}
	resp, err := usecase.Execute(context.Background(), req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, metric, resp.Metrics)
}

func TestMetricUpdateBodyRequestToDomain_InvalidMetricDelta(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := NewMockMetricUpdateBodyService(ctrl)
	usecase := NewMetricUpdateBodyUsecase(mockSvc)
	req := &MetricUpdateBodyRequest{
		Metrics: &domain.Metrics{
			ID:    "metric1",
			Type:  string(domain.Counter),
			Delta: nil,
			Value: nil,
		},
	}
	resp, err := usecase.Execute(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidBodyMetricDelta, err)
}

func TestMetricUpdateBodyUsecase_Execute_InvalidMetricType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := NewMockMetricUpdateBodyService(ctrl)
	usecase := NewMetricUpdateBodyUsecase(mockSvc)
	req := &MetricUpdateBodyRequest{
		Metrics: &domain.Metrics{
			ID:    "metric1",
			Type:  "invalidType",
			Delta: ptrInt64(100),
			Value: nil,
		},
	}
	resp, err := usecase.Execute(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidBodyMetricType, err)
}

func TestMetricUpdateBodyUsecase_Execute_InvalidMetricValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := NewMockMetricUpdateBodyService(ctrl)
	usecase := NewMetricUpdateBodyUsecase(mockSvc)
	req := &MetricUpdateBodyRequest{
		Metrics: &domain.Metrics{
			ID:    "metric1",
			Type:  string(domain.Gauge),
			Delta: nil,
			Value: nil,
		},
	}
	resp, err := usecase.Execute(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidBodyMetricValue, err)
}

func TestMetricUpdateBodyUsecase_Execute_UpdateBatchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := NewMockMetricUpdateBodyService(ctrl)
	metric := &domain.Metrics{
		ID:    "metric1",
		Type:  string(domain.Counter),
		Delta: ptrInt64(100),
		Value: nil,
	}
	mockSvc.EXPECT().
		UpdateBatch(context.Background(), []*domain.Metrics{metric}).
		Return(nil, errors.New("update error")).
		Times(1)
	usecase := NewMetricUpdateBodyUsecase(mockSvc)
	req := &MetricUpdateBodyRequest{
		Metrics: metric,
	}
	resp, err := usecase.Execute(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "update error", err.Error())
}

func ptrInt64(i int64) *int64 {
	return &i
}
