package usecases

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
)

type MetricUpdatesBodyRequest []*MetricUpdateBodyRequest

type MetricUpdatesBodyResponse []*MetricUpdateBodyRequest

type MetricUpdatesBodyService interface {
	Update(ctx context.Context, metrics []*domain.Metrics) ([]*domain.Metrics, error)
}

type MetricUpdatesBodyUsecase struct {
	svc MetricUpdatesBodyService
}

var (
	ErrMetricUpdatesBodyNotProvider  = errors.New("no metrics provided")
	ErrMetricUpdatesBodyInternal     = errors.New("internal error")
	ErrMetricUpdatesBodyTypeRequired = errors.New("metric type required")
	ErrMetricUpdatesBodyInvalidType  = errors.New("metric counter or gauge required")
	ErrMetricUpdatesBodyIDRequired   = errors.New("metric id required")
	ErrMetricUpdatesBodyInvalidValue = errors.New("metric value is invalid")
)

func NewMetricUpdatesBodyUsecase(svc MetricUpdatesBodyService) *MetricUpdatesBodyUsecase {
	return &MetricUpdatesBodyUsecase{svc: svc}
}

func (uc *MetricUpdatesBodyUsecase) Execute(
	ctx context.Context, req *MetricUpdatesBodyRequest,
) (*MetricUpdatesBodyResponse, error) {
	if len(*req) == 0 {
		return nil, ErrMetricUpdatesBodyNotProvider
	}

	var metrics []*domain.Metrics
	for _, r := range *req {
		if r.ID == "" {
			return nil, ErrMetricUpdatesBodyIDRequired
		}
		if r.Type == "" {
			return nil, ErrMetricUpdatesBodyTypeRequired
		}
		if r.Type != string(domain.Counter) && r.Type != string(domain.Gauge) {
			return nil, ErrMetricUpdatesBodyInvalidType
		}

		var metric domain.Metrics
		metric.ID = r.ID
		metric.Type = domain.MetricType(r.Type)

		switch metric.Type {
		case domain.Counter:
			if r.Delta == nil {
				return nil, ErrMetricUpdatesBodyInvalidValue
			}
			metric.Delta = r.Delta
		case domain.Gauge:
			if r.Value == nil {
				return nil, ErrMetricUpdatesBodyInvalidValue
			}
			metric.Value = r.Value
		}

		metrics = append(metrics, &metric)
	}

	updatedMetrics, err := uc.svc.Update(ctx, metrics)
	if err != nil {
		return nil, ErrMetricUpdatesBodyInternal
	}

	var responseMetricResponses MetricUpdatesBodyResponse
	for _, updated := range updatedMetrics {
		responseMetricResponses = append(responseMetricResponses, &MetricUpdateBodyRequest{
			ID:    updated.ID,
			Type:  string(updated.Type),
			Delta: updated.Delta,
			Value: updated.Value,
		})
	}

	return &responseMetricResponses, nil
}
