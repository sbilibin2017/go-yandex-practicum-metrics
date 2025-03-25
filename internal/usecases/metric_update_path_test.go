package usecases_test

import (
	"context"
	"errors"
	"testing"

	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/usecases"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupMetricUpdatePath(t *testing.T) (*gomock.Controller, *usecases.MockMetricUpdatePathService, *usecases.MetricUpdatePathUsecase) {
	ctrl := gomock.NewController(t)
	mockService := usecases.NewMockMetricUpdatePathService(ctrl)
	usecase := usecases.NewMetricUpdatePathUsecase(mockService)
	return ctrl, mockService, usecase
}

func TestSuccessfulCounterUpdate(t *testing.T) {
	ctrl, mockService, usecase := setupMetricUpdatePath(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &usecases.MetricUpdatePathRequest{
		ID:    "test_counter",
		Type:  string(domain.Counter),
		Value: "42",
	}

	metric := &domain.Metrics{
		MetricID: domain.MetricID{
			ID:   "test_counter",
			Type: domain.Counter,
		},
		Delta: func(v int64) *int64 { return &v }(42),
	}

	mockService.EXPECT().Update(ctx, metric).Return(metric, nil)

	resp, err := usecase.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Metric updated successfully", string(*resp))
}

func TestSuccessfulGaugeUpdate(t *testing.T) {
	ctrl, mockService, usecase := setupMetricUpdatePath(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &usecases.MetricUpdatePathRequest{
		ID:    "test_gauge",
		Type:  string(domain.Gauge),
		Value: "3.14",
	}

	metric := &domain.Metrics{
		MetricID: domain.MetricID{
			ID:   "test_gauge",
			Type: domain.Gauge,
		},
		Value: func(v float64) *float64 { return &v }(3.14),
	}

	mockService.EXPECT().Update(ctx, metric).Return(metric, nil)

	resp, err := usecase.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Metric updated successfully", string(*resp))
}

func TestErrorOnInvalidCounterValue(t *testing.T) {
	_, _, usecase := setupMetricUpdatePath(t)

	ctx := context.Background()
	req := &usecases.MetricUpdatePathRequest{
		ID:    "test_counter",
		Type:  string(domain.Counter),
		Value: "invalid",
	}

	resp, err := usecase.Execute(ctx, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, usecases.ErrMetricUpdatePathInvalidCounterValue)
}

func TestErrorOnInvalidGaugeValue(t *testing.T) {
	_, _, usecase := setupMetricUpdatePath(t)

	ctx := context.Background()
	req := &usecases.MetricUpdatePathRequest{
		ID:    "test_gauge",
		Type:  string(domain.Gauge),
		Value: "invalid",
	}

	resp, err := usecase.Execute(ctx, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, usecases.ErrMetricUpdatePathInvalidGaugeValue)
}

func TestErrorFromService(t *testing.T) {
	ctrl, mockService, usecase := setupMetricUpdatePath(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &usecases.MetricUpdatePathRequest{
		ID:    "test_counter",
		Type:  string(domain.Counter),
		Value: "42",
	}

	metric := &domain.Metrics{
		MetricID: domain.MetricID{
			ID:   "test_counter",
			Type: domain.Counter,
		},
		Delta: func(v int64) *int64 { return &v }(42),
	}

	mockService.EXPECT().Update(ctx, metric).Return(nil, errors.New("service error"))

	resp, err := usecase.Execute(ctx, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, usecases.ErrMetricUpdatePathInternal)
}

func TestErrorOnMissingID(t *testing.T) {
	_, _, usecase := setupMetricUpdatePath(t)

	ctx := context.Background()
	req := &usecases.MetricUpdatePathRequest{
		ID:    "",
		Type:  string(domain.Counter),
		Value: "42",
	}

	resp, err := usecase.Execute(ctx, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, usecases.ErrMetricUpdatePathIDRequired)
}

func TestErrorOnMissingType(t *testing.T) {
	_, _, usecase := setupMetricUpdatePath(t)

	ctx := context.Background()
	req := &usecases.MetricUpdatePathRequest{
		ID:    "test_metric",
		Type:  "",
		Value: "42",
	}

	resp, err := usecase.Execute(ctx, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, usecases.ErrMetricUpdatePathTypeRequired)
}

func TestErrorOnInvalidType(t *testing.T) {
	_, _, usecase := setupMetricUpdatePath(t)

	ctx := context.Background()
	req := &usecases.MetricUpdatePathRequest{
		ID:    "test_metric",
		Type:  "invalid_type",
		Value: "42",
	}

	resp, err := usecase.Execute(ctx, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, usecases.ErrMetricUpdatePathInvalidType)
}

func TestErrorOnMissingValue(t *testing.T) {
	_, _, usecase := setupMetricUpdatePath(t)

	ctx := context.Background()
	req := &usecases.MetricUpdatePathRequest{
		ID:    "test_metric",
		Type:  string(domain.Counter),
		Value: "",
	}

	resp, err := usecase.Execute(ctx, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, usecases.ErrMetricUpdatePathValueRequired)
}
