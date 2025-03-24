package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
)

type MetricUpdateSaveRepository interface {
	Save(ctx context.Context, metric *types.Metrics) bool
}

type MetricUpdateFilterRepository interface {
	Filter(ctx context.Context, filter types.MetricID) (*types.Metrics, bool)
}

type MetricUpdateService struct {
	save   MetricUpdateSaveRepository
	filter MetricUpdateFilterRepository
}

// Метод обновления метрики
func (svc MetricUpdateService) Update(
	ctx context.Context, metric *types.Metrics,
) (*types.Metrics, error) {
	existingMetric, found := svc.filter.Filter(ctx, types.MetricID{
		ID:   metric.MetricID.ID,
		Type: metric.MetricID.Type,
	})
	if !found {
		if ok := svc.save.Save(ctx, metric); !ok {
			return nil, errors.ErrMetricInternal
		}
		return metric, nil
	}

	switch metric.MetricID.Type {
	case string(types.Gauge):
		existingMetric.Value = metric.Value
	case string(types.Counter):
		*existingMetric.Delta += *metric.Delta
	}

	if ok := svc.save.Save(ctx, existingMetric); !ok {
		return nil, errors.ErrMetricInternal
	}

	return existingMetric, nil
}
