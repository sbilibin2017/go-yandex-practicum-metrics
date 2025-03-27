package usecases

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
)

type MetricUpdateBatchBodyService interface {
	UpdateBatch(ctx context.Context, metrics []*domain.Metrics) ([]*domain.Metrics, error)
}

type MetricUpdateBatchBodyUsecase struct {
	svc MetricUpdateBatchBodyService
}

func NewMetricUpdateBatchBodyUsecase(svc MetricUpdateBatchBodyService) *MetricUpdateBatchBodyUsecase {
	return &MetricUpdateBatchBodyUsecase{svc: svc}
}

func (uc MetricUpdateBatchBodyUsecase) Execute(
	ctx context.Context, req MetricUpdateBatchBodyRequest,
) ([]*MetricUpdateBodyResponse, error) {
	metrics, err := req.ToDomain()
	if err != nil {
		return nil, err
	}
	metrics, err = uc.svc.UpdateBatch(ctx, metrics)
	if err != nil {
		return nil, err
	}
	resp := MetricUpdateBatchBodyResponse{}
	resp.FromDomain(metrics)
	return resp, nil
}

type MetricUpdateBatchBodyRequest []*MetricUpdateBodyRequest

func (reqs MetricUpdateBatchBodyRequest) ToDomain() ([]*domain.Metrics, error) {
	var metrics []*domain.Metrics
	for _, r := range reqs {
		switch r.Type {
		case string(domain.Gauge):
			if r.Value == nil {
				return nil, errors.InvalidMetricValueError
			}
		case string(domain.Counter):
			if r.Delta == nil {
				return nil, errors.InvalidMetricValueError
			}
		default:
			return nil, errors.InvalidMetricTypeError
		}
		if r.ID == "" {
			return nil, errors.MissingMetricNameError
		}
		metrics = append(metrics, r.Metrics)
	}
	return metrics, nil
}

type MetricUpdateBatchBodyResponse []*MetricUpdateBodyResponse

func (resp *MetricUpdateBatchBodyResponse) FromDomain(metrics []*domain.Metrics) {
	for _, metric := range metrics {
		*resp = append(*resp, &MetricUpdateBodyResponse{Metrics: metric})
	}
}
