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

func setupMetricUpdateBody(t *testing.T) (*gomock.Controller, *usecases.MockMetricUpdateBodyService, *usecases.MetricUpdateBodyUsecase) {
	ctrl := gomock.NewController(t)
	mockService := usecases.NewMockMetricUpdateBodyService(ctrl)
	usecase := usecases.NewMetricUpdateBodyUsecase(mockService)
	return ctrl, mockService, usecase
}

func TestSuccessCounterUpdate(t *testing.T) {
	ctrl, mockService, usecase := setupMetricUpdateBody(t)
	defer ctrl.Finish()

	ctx := context.Background()
	metricID := "test_metric"
	metricDelta := int64(10)
	req := &usecases.MetricUpdateBodyRequest{
		ID:    metricID,
		Type:  string(domain.Counter),
		Delta: &metricDelta,
	}
	updatedMetric := &domain.Metrics{
		MetricID: domain.MetricID{
			ID:   metricID,
			Type: domain.Counter,
		},
		Delta: &metricDelta,
	}

	mockService.EXPECT().Update(ctx, gomock.Any()).Return(updatedMetric, nil)

	resp, err := usecase.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, metricID, resp.ID)
	assert.Equal(t, string(domain.Counter), resp.Type)
	assert.Equal(t, metricDelta, *resp.Delta)
}

func TestSuccessGaugeUpdate(t *testing.T) {
	ctrl, mockService, usecase := setupMetricUpdateBody(t)
	defer ctrl.Finish()

	ctx := context.Background()
	metricID := "test_metric"
	metricDelta := int64(20)
	req := &usecases.MetricUpdateBodyRequest{
		ID:    metricID,
		Type:  string(domain.Counter),
		Delta: &metricDelta,
	}
	updatedMetric := &domain.Metrics{
		MetricID: domain.MetricID{
			ID:   metricID,
			Type: domain.Counter,
		},
		Delta: &metricDelta,
	}

	mockService.EXPECT().Update(ctx, gomock.Any()).Return(updatedMetric, nil)

	resp, err := usecase.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, metricID, resp.ID)
	assert.Equal(t, string(domain.Counter), resp.Type)
	assert.Equal(t, metricDelta, *resp.Delta)
}

func TestErrorMissingMetricID(t *testing.T) {
	_, _, usecase := setupMetricUpdateBody(t)

	ctx := context.Background()
	req := &usecases.MetricUpdateBodyRequest{Type: string(domain.Counter), Delta: new(int64)}

	resp, err := usecase.Execute(ctx, req)

	assert.ErrorIs(t, err, usecases.ErrMetricUpdateBodyIDRequired)
	assert.Nil(t, resp)
}

func TestErrorInvalidMetricType(t *testing.T) {
	_, _, usecase := setupMetricUpdateBody(t)

	ctx := context.Background()
	metricID := "test_metric"
	req := &usecases.MetricUpdateBodyRequest{ID: metricID, Type: "invalid", Delta: new(int64)}

	resp, err := usecase.Execute(ctx, req)

	assert.ErrorIs(t, err, usecases.ErrMetricUpdateBodyInvalidType)
	assert.Nil(t, resp)
}

func TestErrorBodyServiceFailure(t *testing.T) {
	ctrl, mockService, usecase := setupMetricUpdateBody(t)
	defer ctrl.Finish()

	ctx := context.Background()
	metricID := "test_metric"
	metricDelta := int64(10)
	req := &usecases.MetricUpdateBodyRequest{ID: metricID, Type: string(domain.Counter), Delta: &metricDelta}

	mockService.EXPECT().Update(ctx, gomock.Any()).Return(nil, errors.New("internal error"))

	resp, err := usecase.Execute(ctx, req)

	assert.ErrorIs(t, err, usecases.ErrMetricUpdateBodyInternal)
	assert.Nil(t, resp)
}

func TestErrorMissingCounterValue(t *testing.T) {
	_, _, usecase := setupMetricUpdateBody(t)

	ctx := context.Background()
	metricID := "test_metric"
	req := &usecases.MetricUpdateBodyRequest{ID: metricID, Type: string(domain.Counter)}

	resp, err := usecase.Execute(ctx, req)

	assert.ErrorIs(t, err, usecases.ErrMetricUpdateBodyInvalidValue)
	assert.Nil(t, resp)
}

func TestErrorMissingGaugeValue(t *testing.T) {
	_, _, usecase := setupMetricUpdateBody(t)

	ctx := context.Background()
	metricID := "test_metric"
	req := &usecases.MetricUpdateBodyRequest{ID: metricID, Type: string(domain.Gauge)}

	resp, err := usecase.Execute(ctx, req)

	assert.ErrorIs(t, err, usecases.ErrMetricUpdateBodyInvalidValue)
	assert.Nil(t, resp)
}

func TestErrorServiceFailure(t *testing.T) {
	ctrl, mockService, usecase := setupMetricUpdateBody(t)
	defer ctrl.Finish()

	ctx := context.Background()
	metricID := "test_metric"
	metricDelta := int64(10)
	req := &usecases.MetricUpdateBodyRequest{ID: metricID, Type: string(domain.Counter), Delta: &metricDelta}

	mockService.EXPECT().Update(ctx, gomock.Any()).Return(nil, errors.New("internal error"))

	resp, err := usecase.Execute(ctx, req)

	assert.ErrorIs(t, err, usecases.ErrMetricUpdateBodyInternal)
	assert.Nil(t, resp)
}

func TestErrorMissingType(t *testing.T) {
	_, _, usecase := setupMetricUpdateBody(t)

	ctx := context.Background()
	req := &usecases.MetricUpdateBodyRequest{ID: "test_metric"}

	resp, err := usecase.Execute(ctx, req)

	assert.ErrorIs(t, err, usecases.ErrMetricUpdateBodyTypeRequired)
	assert.Nil(t, resp)
}

func TestSuccessGaugeWithValueUpdate(t *testing.T) {
	ctrl, mockService, usecase := setupMetricUpdateBody(t)
	defer ctrl.Finish()

	ctx := context.Background()
	metricID := "test_metric"
	metricValue := float64(25.5)

	req := &usecases.MetricUpdateBodyRequest{
		ID:    metricID,
		Type:  string(domain.Gauge),
		Value: &metricValue,
	}

	updatedMetric := &domain.Metrics{
		MetricID: domain.MetricID{
			ID:   metricID,
			Type: domain.Gauge,
		},
		Value: &metricValue,
	}

	mockService.EXPECT().Update(ctx, gomock.Any()).Return(updatedMetric, nil)

	resp, err := usecase.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, metricID, resp.ID)
	assert.Equal(t, string(domain.Gauge), resp.Type)
	assert.Equal(t, metricValue, *resp.Value)
}
