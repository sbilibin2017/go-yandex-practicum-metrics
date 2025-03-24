package usecases

import (
	"context"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/formatters"
	"go-yandex-practicum-metrics/internal/types"
	"go-yandex-practicum-metrics/internal/validators"
)

type MetricUpdatePathService interface {
	Update(ctx context.Context, metric *types.Metrics) (*types.Metrics, error)
}

type MetricUpdatePathUsecase struct {
	svc MetricUpdatePathService
}

// Execute processes the update request for a metric path.
func (uc *MetricUpdatePathUsecase) Execute(
	ctx context.Context, req *types.MetricUpdatePathRequest,
) (*types.MetricUpdatePathResponse, error) {
	if err := validators.ValidateMetricType(req.Type); err != nil {
		return nil, err
	}
	if err := validators.ValidateMetricName(req.Name); err != nil {
		return nil, err
	}
	if err := validators.ValidateMetricValue(req.Value); err != nil {
		return nil, err
	}

	var metric types.Metrics
	metric.ID = req.Name
	metric.Type = types.MetricType(req.Type)

	switch metric.Type {
	case types.Counter:
		parsedValue, valid := formatters.ParseInt64(req.Value)
		if !valid {
			return nil, errors.ErrInvalidCounterValue
		}
		metric.Delta = &parsedValue
	case types.Gauge:
		parsedValue, valid := formatters.ParseFloat64(req.Value)
		if !valid {
			return nil, errors.ErrInvalidGaugeValue
		}
		metric.Value = &parsedValue
	}

	updatedMetric, err := uc.svc.Update(ctx, &metric)
	if err != nil {
		return nil, err
	}

	return &types.MetricUpdatePathResponse{
		MetricID: types.MetricID{
			Type: updatedMetric.Type,
			ID:   updatedMetric.ID,
		},
		Delta: updatedMetric.Delta,
		Value: updatedMetric.Value,
	}, nil
}
