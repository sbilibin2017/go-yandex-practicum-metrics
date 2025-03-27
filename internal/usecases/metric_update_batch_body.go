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
	ctx context.Context, req []*MetricUpdateBodyRequest,
) ([]*MetricUpdateBodyResponse, error) {
	metrics, err := MetricUpdateBatchBodyRequestToDomain(req)
	if err != nil {
		return nil, err
	}
	metrics, err = uc.svc.UpdateBatch(ctx, metrics)
	if err != nil {
		return nil, err
	}
	resp := MetricUpdateBatchBodyResponseFromDomain(metrics)
	return resp, nil
}

type MetricUpdateBatchBodyRequest []*MetricUpdateBodyRequest

func MetricUpdateBatchBodyRequestToDomain(req []*MetricUpdateBodyRequest) ([]*domain.Metrics, error) {
	var metrics []*domain.Metrics
	for _, r := range req {
		switch r.Type {
		case string(domain.Gauge):
			if r.Value == nil {
				return nil, ErrInvalidBodyBatchMetricValue
			}
		case string(domain.Counter):
			if r.Delta == nil {
				return nil, ErrInvalidBodyBatchMetricDelta
			}
		default:
			return nil, ErrInvalidBodyBatchMetricType
		}
		metrics = append(metrics, r.Metrics)
	}
	return metrics, nil
}

type MetricUpdateBatchBodyResponse []*MetricUpdateBodyResponse

func MetricUpdateBatchBodyResponseFromDomain(metrics []*domain.Metrics) []*MetricUpdateBodyResponse {
	var responses []*MetricUpdateBodyResponse
	for _, metric := range metrics {
		responses = append(responses, &MetricUpdateBodyResponse{Metrics: metric})
	}
	return responses
}

var (
	ErrInvalidBodyBatchMetricType  = errors.New("invalid metric type")
	ErrInvalidBodyBatchMetricValue = errors.New("invalid metric value")
	ErrInvalidBodyBatchMetricDelta = errors.New("invalid metric delta")
	ErrMissingBodyBatchMetricName  = errors.New("missing metric name")
)
