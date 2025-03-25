package usecases

import (
	"context"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
)

type MetricGetBodyService interface {
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

type MetricGetBodyUsecase struct {
	svc MetricGetBodyService
}

func NewMetricGetBodyUsecase(
	svc MetricGetBodyService,
) *MetricGetBodyUsecase {
	return &MetricGetBodyUsecase{svc: svc}
}

func (uc *MetricGetBodyUsecase) Execute(
	ctx context.Context, req *types.MetricGetByTypeAndIDBodyRequest,
) (*types.MetricGetByTypeAndIDBodyResponse, error) {
	if req.ID == "" {
		return nil, errors.ErrMetricIDRequired
	}
	if req.Type == "" {
		return nil, errors.ErrMetricTypeRequired
	}
	if req.Type != string(types.Counter) && req.Type != string(types.Gauge) {
		return nil, errors.ErrMetricInvalidType
	}

	var id types.MetricID
	id.ID = req.ID
	id.Type = types.MetricType(req.Type)

	metric, err := uc.svc.Get(ctx, id)
	if err != nil {
		return nil, errors.ErrMetricInternal
	}
	if metric == nil {
		return nil, errors.ErrMetricNotFound
	}

	response := types.MetricGetByTypeAndIDBodyResponse{
		ID:    metric.ID,
		Type:  string(metric.Type),
		Delta: metric.Delta,
		Value: metric.Value,
	}

	return &response, nil
}
