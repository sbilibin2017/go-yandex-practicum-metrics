package usecases

import (
	"context"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
	"strconv"
)

type MetricUpdatePathService interface {
	Update(ctx context.Context, metric *types.Metrics) (*types.Metrics, error)
}

type MetricUpdatePathUsecase struct {
	svc MetricUpdatePathService
}

func NewMetricUpdatePathUsecase(svc MetricUpdatePathService) *MetricUpdatePathUsecase {
	return &MetricUpdatePathUsecase{svc: svc}
}

func (uc *MetricUpdatePathUsecase) Execute(
	ctx context.Context, req *types.MetricUpdatePathRequest,
) (*types.MetricUpdatePathResponse, error) {
	var emptyString = ""

	if req.Name == emptyString {
		return nil, errors.ErrMetricIDRequired
	}
	if req.Type == emptyString {
		return nil, errors.ErrMetricTypeRequired
	}
	if req.Type != string(types.Counter) && req.Type != string(types.Gauge) {
		return nil, errors.ErrMetricInvalidType
	}
	if req.Value == emptyString {
		return nil, errors.ErrMetricValueRequired
	}

	var metric types.Metrics
	metric.ID = req.Name
	metric.Type = types.MetricType(req.Type)

	switch metric.Type {
	case types.Counter:
		value, err := strconv.ParseInt(req.Value, 10, 64)
		if err != nil {
			return nil, errors.ErrMetricInvalidDelta
		}
		metric.Delta = &value
	case types.Gauge:
		value, err := strconv.ParseFloat(req.Value, 64)
		if err != nil {
			return nil, errors.ErrMetricInvalidValue
		}
		metric.Value = &value
	}

	_, err := uc.svc.Update(ctx, &metric)
	if err != nil {
		return nil, errors.ErrMetricInternal
	}

	response := types.MetricUpdatePathResponse("Metric updated successfully")

	return &response, nil
}
