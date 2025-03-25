package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
)

type MetricUpdateSaveRepository interface {
	Save(ctx context.Context, metric *types.Metrics) error
}

type MetricUpdateFilterRepository interface {
	Filter(ctx context.Context, filter types.MetricID) (*types.Metrics, error)
}

type MetricUpdateService struct {
	save   MetricUpdateSaveRepository
	filter MetricUpdateFilterRepository
}

func NewMetricUpdateService(
	save MetricUpdateSaveRepository,
	filter MetricUpdateFilterRepository,
) *MetricUpdateService {
	return &MetricUpdateService{
		save:   save,
		filter: filter,
	}
}

// Метод обновления метрики
func (svc MetricUpdateService) Update(
	ctx context.Context, metric *types.Metrics,
) (*types.Metrics, error) {
	existingMetric, err := svc.filter.Filter(ctx, types.MetricID{
		ID:   metric.MetricID.ID,
		Type: metric.MetricID.Type,
	})
	if err != nil {
		return nil, errors.ErrMetricInternal
	}

	if existingMetric == nil {
		if err := svc.save.Save(ctx, metric); err != nil {
			return nil, errors.ErrMetricInternal
		}
		return metric, nil
	}

	switch metric.Type {
	case types.Gauge:
		existingMetric.Value = metric.Value
	case types.Counter:
		*existingMetric.Delta += *metric.Delta
	}

	if err := svc.save.Save(ctx, existingMetric); err != nil {
		return nil, errors.ErrMetricInternal
	}

	return existingMetric, nil
}
