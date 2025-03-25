package usecases

import (
	"context"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
)

type MetricUpdateBodyService interface {
	Update(ctx context.Context, metric *types.Metrics) (*types.Metrics, error)
}

type MetricUpdateBodyUsecase struct {
	svc MetricUpdateBodyService
}

func NewMetricUpdateBodyUsecase(svc MetricUpdateBodyService) *MetricUpdateBodyUsecase {
	return &MetricUpdateBodyUsecase{svc: svc}
}

func (uc *MetricUpdateBodyUsecase) Execute(
	ctx context.Context, req *types.MetricUpdateBodyRequest,
) (*types.MetricUpdateBodyResponse, error) {
	if req.ID == "" {
		return nil, errors.ErrMetricIDRequired
	}
	if req.Type == "" {
		return nil, errors.ErrMetricTypeRequired
	}
	if req.Type != string(types.Counter) && req.Type != string(types.Gauge) {
		return nil, errors.ErrMetricInvalidType
	}

	var metric types.Metrics
	metric.ID = req.ID
	metric.Type = types.MetricType(req.Type)

	switch metric.Type {
	case types.Counter:
		if req.Delta == nil {
			return nil, errors.ErrMetricInvalidDelta
		}
		metric.Delta = req.Delta
	case types.Gauge:
		if req.Value == nil {
			return nil, errors.ErrMetricInvalidValue
		}
		metric.Value = req.Value
	}

	updated, err := uc.svc.Update(ctx, &metric)
	if err != nil {
		return nil, errors.ErrMetricInternal
	}

	response := types.MetricUpdateBodyResponse{
		MetricUpdateBodyRequest: types.MetricUpdateBodyRequest{
			ID:    updated.ID,
			Type:  string(updated.Type),
			Delta: updated.Delta,
			Value: updated.Value,
		},
	}

	return &response, nil
}
