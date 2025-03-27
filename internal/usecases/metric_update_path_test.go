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

func TestMetricUpdatePathRequest_ToDomain_Gauge_Success(t *testing.T) {
	req := &MetricUpdatePathRequest{
		Type:  "gauge",
		Name:  "metric1",
		Value: "123.45", // valid float value for Gauge
	}

	metric, err := MetricUpdatePathRequestToDomain(req)

	require.NoError(t, err)
	assert.Equal(t, "metric1", metric.ID)
	assert.Equal(t, "gauge", metric.Type)
	assert.Equal(t, float64(123.45), *metric.Value)
	assert.Nil(t, metric.Delta) // Ensure Delta is nil for Gauge type
}

func TestMetricUpdatePathRequest_ToDomain_Gauge_InvalidValue(t *testing.T) {
	req := &MetricUpdatePathRequest{
		Type:  "gauge",
		Name:  "metric1",
		Value: "invalid", // invalid float value for Gauge
	}

	metric, err := MetricUpdatePathRequestToDomain(req)

	require.Error(t, err)
	assert.Nil(t, metric)
	assert.Equal(t, ErrInvalidPathMetricValue, err)
}

func TestMetricUpdatePathUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricUpdatePathService(ctrl)
	usecase := NewMetricUpdatePathUsecase(mockSvc)

	req := &MetricUpdatePathRequest{
		Type:  "counter",
		Name:  "metric1",
		Value: "10",
	}
	metric := &domain.Metrics{
		ID:    "metric1",
		Type:  "counter",
		Delta: new(int64),
		Value: new(float64),
	}
	mockSvc.EXPECT().UpdateBatch(gomock.Any(), gomock.Any()).Return([]*domain.Metrics{metric}, nil)

	result, err := usecase.Execute(context.Background(), req)

	require.NoError(t, err)

	expected := MetricUpdateSuccessMessage

	// Dereference the result to compare []byte directly
	assert.Equal(t, expected, []byte(*result))
}

func TestMetricUpdatePathUsecase_Execute_InvalidMetricType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricUpdatePathService(ctrl)
	usecase := NewMetricUpdatePathUsecase(mockSvc)

	req := &MetricUpdatePathRequest{
		Type:  "invalidType",
		Name:  "metric1",
		Value: "10",
	}

	result, err := usecase.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, err, ErrInvalidPathMetricType)
}

func TestMetricUpdatePathUsecase_Execute_InvalidMetricValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricUpdatePathService(ctrl)
	usecase := NewMetricUpdatePathUsecase(mockSvc)

	req := &MetricUpdatePathRequest{
		Type:  "counter",
		Name:  "metric1",
		Value: "invalidValue",
	}

	result, err := usecase.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, err, ErrInvalidPathMetricValue)
}

func TestMetricUpdatePathUsecase_Execute_UpdateBatchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockMetricUpdatePathService(ctrl)
	usecase := NewMetricUpdatePathUsecase(mockSvc)

	req := &MetricUpdatePathRequest{
		Type:  "counter",
		Name:  "metric1",
		Value: "10",
	}
	mockSvc.EXPECT().UpdateBatch(gomock.Any(), gomock.Any()).Return(nil, errors.New("update error"))

	result, err := usecase.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, err.Error(), "update error")
}
