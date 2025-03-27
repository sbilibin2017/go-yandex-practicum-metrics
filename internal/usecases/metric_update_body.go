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
	metrics, err := req.ToDomain()
	if err != nil {
		return nil, err
	}
	metrics, err = uc.svc.UpdateBatch(ctx, metrics)
	if err != nil {
		return nil, err
	}
	var resp *MetricUpdateBodyResponse
	resp = resp.FromDomain(metrics)
	return resp, nil
}

type MetricUpdateBodyRequest struct {
	*domain.Metrics
}

func (req *MetricUpdateBodyRequest) ToDomain() ([]*domain.Metrics, error) {
	var metricType string
	switch req.Type {
	case string(domain.Gauge):
		metricType = string(domain.Gauge)
	case string(domain.Counter):
		metricType = string(domain.Counter)
	default:
		return nil, errors.New("invalid metric type")
	}
	if req.ID == "" {
		return nil, errors.New("missing metric id")
	}
	switch metricType {
	case string(domain.Gauge):
		if req.Value == nil {
			return nil, errors.New("invalid metric value")
		}
	case string(domain.Counter):
		if req.Delta == nil {
			return nil, errors.New("invalid metric delta")
		}
	}
	return []*domain.Metrics{req.Metrics}, nil
}

type MetricUpdateBodyResponse struct {
	*domain.Metrics
}

func (r *MetricUpdateBodyResponse) FromDomain(metrics []*domain.Metrics) *MetricUpdateBodyResponse {
	return &MetricUpdateBodyResponse{
		Metrics: metrics[0],
	}
}
