package usecases

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
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
				return nil, errors.New("invalid metric value")
			}
		case string(domain.Counter):
			if r.Delta == nil {
				return nil, errors.New("invalid metric delta")
			}
		default:
			return nil, errors.New("invalid metric type")
		}
		if r.ID == "" {
			return nil, errors.New("missing metric id")
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
