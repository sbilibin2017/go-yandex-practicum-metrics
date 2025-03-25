package usecases

import (
	"context"

	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
)

type MetricUpdatesBodyService interface {
	Updates(ctx context.Context, metrics []*types.Metrics) ([]*types.Metrics, error)
}

type MetricUpdatesBodyUsecase struct {
	svc MetricUpdatesBodyService
}

func NewMetricUpdatesBodyUsecase(svc MetricUpdatesBodyService) *MetricUpdatesBodyUsecase {
	return &MetricUpdatesBodyUsecase{svc: svc}
}

func (uc *MetricUpdatesBodyUsecase) Execute(
	ctx context.Context, req *types.MetricUpdatesBodyRequest,
) (*types.MetricUpdatesBodyResponse, error) {
	if len(*req) == 0 {
		return nil, errors.ErrMetricBodyNotProvider
	}

	var metrics []*types.Metrics
	for _, r := range *req {
		if r.ID == "" {
			return nil, errors.ErrMetricIDRequired
		}
		if r.Type == "" {
			return nil, errors.ErrMetricTypeRequired
		}
		if r.Type != string(types.Counter) && r.Type != string(types.Gauge) {
			return nil, errors.ErrMetricInvalidType
		}

		var metric types.Metrics
		metric.ID = r.ID
		metric.Type = types.MetricType(r.Type)

		switch metric.Type {
		case types.Counter:
			if r.Delta == nil {
				return nil, errors.ErrMetricInvalidDelta
			}
			metric.Delta = r.Delta
		case types.Gauge:
			if r.Value == nil {
				return nil, errors.ErrMetricInvalidValue
			}
			metric.Value = r.Value
		}

		metrics = append(metrics, &metric)
	}

	updatedMetrics, err := uc.svc.Updates(ctx, metrics)
	if err != nil {
		return nil, errors.ErrMetricInternal
	}

	var responseMetricResponses types.MetricUpdatesBodyResponse
	for _, updated := range updatedMetrics {
		responseMetricResponses = append(responseMetricResponses, &types.MetricUpdateBodyRequest{
			ID:    updated.ID,
			Type:  string(updated.Type),
			Delta: updated.Delta,
			Value: updated.Value,
		})
	}

	return &responseMetricResponses, nil
}
