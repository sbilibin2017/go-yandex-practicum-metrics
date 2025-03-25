package usecases

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
	"strconv"
)

type MetricUpdatePathRequest struct {
	ID    string
	Type  string
	Value string
}

type MetricUpdatePathResponse string

type MetricUpdatePathService interface {
	Update(ctx context.Context, metric *domain.Metrics) (*domain.Metrics, error)
}

type MetricUpdatePathUsecase struct {
	svc MetricUpdatePathService
}

var (
	ErrMetricUpdatePathInternal            = errors.New("internal error")
	ErrMetricUpdatePathTypeRequired        = errors.New("metric type required")
	ErrMetricUpdatePathInvalidType         = errors.New("metric counter or gauge required")
	ErrMetricUpdatePathIDRequired          = errors.New("metric id required")
	ErrMetricUpdatePathValueRequired       = errors.New("metric value required")
	ErrMetricUpdatePathInvalidCounterValue = errors.New("metric invalid counter value")
	ErrMetricUpdatePathInvalidGaugeValue   = errors.New("metric invalid gauge value")
)

func NewMetricUpdatePathUsecase(svc MetricUpdatePathService) *MetricUpdatePathUsecase {
	return &MetricUpdatePathUsecase{svc: svc}
}

// Execute processes the update request for a metric path.
func (uc *MetricUpdatePathUsecase) Execute(
	ctx context.Context, req *MetricUpdatePathRequest,
) (*MetricUpdatePathResponse, error) {
	var emptyString = ""

	if req.ID == emptyString {
		return nil, ErrMetricUpdatePathIDRequired
	}
	if req.Type == emptyString {
		return nil, ErrMetricUpdatePathTypeRequired
	}
	if req.Type != string(domain.Counter) && req.Type != string(domain.Gauge) {
		return nil, ErrMetricUpdatePathInvalidType
	}
	if req.Value == emptyString {
		return nil, ErrMetricUpdatePathValueRequired
	}

	var metric domain.Metrics
	metric.ID = req.ID
	metric.Type = domain.MetricType(req.Type)

	switch metric.Type {
	case domain.Counter:
		value, err := strconv.ParseInt(req.Value, 10, 64)
		if err != nil {
			return nil, ErrMetricUpdatePathInvalidCounterValue
		}
		metric.Delta = &value
	case domain.Gauge:
		value, err := strconv.ParseFloat(req.Value, 64)
		if err != nil {
			return nil, ErrMetricUpdatePathInvalidGaugeValue
		}
		metric.Value = &value
	}

	_, err := uc.svc.Update(ctx, &metric)
	if err != nil {
		return nil, ErrMetricUpdatePathInternal
	}

	successMessage := MetricUpdatePathResponse("Metric updated successfully")
	return &successMessage, nil
}
