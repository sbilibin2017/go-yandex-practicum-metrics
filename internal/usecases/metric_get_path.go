package usecases

import (
	"context"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
	"strconv"
)

type MetricGetPathService interface {
	Get(ctx context.Context, id types.MetricID) (*types.Metrics, error)
}

type MetricGetPathUsecase struct {
	svc MetricGetPathService
}

func NewMetricGetPathUsecase(
	svc MetricGetPathService,
) *MetricGetPathUsecase {
	return &MetricGetPathUsecase{svc: svc}
}

func (uc *MetricGetPathUsecase) Execute(
	ctx context.Context, req *types.MetricGetByTypeAndIDPathRequest,
) (*types.MetricGetByTypeAndIDPathResponse, error) {
	if req.Name == "" {
		return nil, errors.ErrMetricIDRequired
	}
	if req.Type == "" {
		return nil, errors.ErrMetricTypeRequired
	}
	if req.Type != string(types.Counter) && req.Type != string(types.Gauge) {
		return nil, errors.ErrMetricInvalidType
	}

	var id types.MetricID
	id.ID = req.Name
	id.Type = types.MetricType(req.Type)

	metric, err := uc.svc.Get(ctx, id)
	if err != nil {
		return nil, errors.ErrMetricInternal
	}
	if metric == nil {
		return nil, errors.ErrMetricNotFound
	}

	var value string
	switch metric.Type {
	case types.Counter:
		if metric.Delta != nil {
			value = strconv.FormatInt(*metric.Delta, 10)
		}
	case types.Gauge:
		if metric.Value != nil {
			value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
		}
	}

	response := types.MetricGetByTypeAndIDPathResponse(value)
	return &response, nil
}
