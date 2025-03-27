package usecases

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
)

type MetricUpdateBodyService interface {
	UpdateBatch(ctx context.Context, metrics []*domain.Metrics) ([]*domain.Metrics, error)
}

type MetricUpdateBodyUsecase struct {
	svc MetricUpdateBodyService
}

func NewMetricUpdateBodyUsecase(svc MetricUpdateBodyService) *MetricUpdateBodyUsecase {
	return &MetricUpdateBodyUsecase{svc: svc}
}

func (uc MetricUpdateBodyUsecase) Execute(
	ctx context.Context, req *MetricUpdateBodyRequest,
) (*MetricUpdateBodyResponse, error) {
	metrics, err := MetricUpdateBodyRequestToDomain(req)
	if err != nil {
		return nil, err
	}
	metrics, err = uc.svc.UpdateBatch(ctx, metrics)
	if err != nil {
		return nil, err
	}
	resp := MetricUpdateBodyResponseFromDomain(metrics)
	return resp, nil
}

type MetricUpdateBodyRequest struct {
	*domain.Metrics
}

// Replace the method `ToDomain` with a function
func MetricUpdateBodyRequestToDomain(req *MetricUpdateBodyRequest) ([]*domain.Metrics, error) {
	var metricType string
	switch req.Type {
	case string(domain.Gauge):
		metricType = string(domain.Gauge)
	case string(domain.Counter):
		metricType = string(domain.Counter)
	default:
		return nil, ErrInvalidBodyMetricType
	}
	switch metricType {
	case string(domain.Gauge):
		if req.Value == nil {
			return nil, ErrInvalidBodyMetricValue
		}
	case string(domain.Counter):
		if req.Delta == nil {
			return nil, ErrInvalidBodyMetricDelta
		}
	}
	return []*domain.Metrics{req.Metrics}, nil
}

type MetricUpdateBodyResponse struct {
	*domain.Metrics
}

// Replace the method `FromDomain` with a function
func MetricUpdateBodyResponseFromDomain(metrics []*domain.Metrics) *MetricUpdateBodyResponse {
	return &MetricUpdateBodyResponse{
		Metrics: metrics[0],
	}
}

var (
	ErrInvalidBodyMetricType  = errors.New("invalid metric type")
	ErrInvalidBodyMetricValue = errors.New("invalid metric value")
	ErrInvalidBodyMetricDelta = errors.New("invalid metric delta")
	ErrMissingBodyMetricName  = errors.New("missing metric name")
)
