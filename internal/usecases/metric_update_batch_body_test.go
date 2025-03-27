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

// ptrFloat64 - вспомогательная функция для создания указателя на float64
func ptrFloat64(f float64) *float64 {
	return &f
}

func TestMetricUpdateBatchBodyUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricUpdateBatchBodyService(ctrl)

	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  string(domain.Counter),
			Delta: ptrInt64(100),
		},
		{
			ID:    "metric2",
			Type:  string(domain.Gauge),
			Value: ptrFloat64(42.5),
		},
	}

	mockSvc.EXPECT().
		UpdateBatch(context.Background(), metrics).
		Return(metrics, nil).
		Times(1)

	usecase := NewMetricUpdateBatchBodyUsecase(mockSvc)

	req := []*MetricUpdateBodyRequest{
		{Metrics: metrics[0]},
		{Metrics: metrics[1]},
	}

	resp, err := usecase.Execute(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp, 2)
	assert.Equal(t, metrics[0], resp[0].Metrics)
	assert.Equal(t, metrics[1], resp[1].Metrics)
}

func TestMetricUpdateBatchBodyUsecase_Execute_InvalidMetricType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricUpdateBatchBodyService(ctrl)
	usecase := NewMetricUpdateBatchBodyUsecase(mockSvc)

	req := []*MetricUpdateBodyRequest{
		{Metrics: &domain.Metrics{
			ID:   "metric1",
			Type: "invalidType",
		}},
	}

	resp, err := usecase.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidBodyBatchMetricType, err)
}

func TestMetricUpdateBatchBodyUsecase_Execute_InvalidMetricValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricUpdateBatchBodyService(ctrl)
	usecase := NewMetricUpdateBatchBodyUsecase(mockSvc)

	req := []*MetricUpdateBodyRequest{
		{Metrics: &domain.Metrics{
			ID:   "metric1",
			Type: string(domain.Gauge),
		}}, // Отсутствует Value для Gauge
	}

	resp, err := usecase.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidBodyBatchMetricValue, err)
}

func TestMetricUpdateBatchBodyUsecase_Execute_InvalidMetricDelta(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricUpdateBatchBodyService(ctrl)
	usecase := NewMetricUpdateBatchBodyUsecase(mockSvc)

	req := []*MetricUpdateBodyRequest{
		{Metrics: &domain.Metrics{
			ID:   "metric1",
			Type: string(domain.Counter),
		}}, // Отсутствует Delta для Counter
	}

	resp, err := usecase.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidBodyBatchMetricDelta, err)
}

func TestMetricUpdateBatchBodyUsecase_Execute_UpdateBatchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricUpdateBatchBodyService(ctrl)

	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  string(domain.Counter),
			Delta: ptrInt64(100),
		},
	}

	mockSvc.EXPECT().
		UpdateBatch(context.Background(), metrics).
		Return(nil, errors.New("update error")).
		Times(1)

	usecase := NewMetricUpdateBatchBodyUsecase(mockSvc)

	req := []*MetricUpdateBodyRequest{
		{Metrics: metrics[0]},
	}

	resp, err := usecase.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "update error", err.Error())
}
