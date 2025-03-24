package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
)

type MetricUpdatesSaveRepository interface {
	Save(ctx context.Context, metric *types.Metrics, tx Transaction) bool
}

type MetricUpdatesFilterRepository interface {
	Filter(ctx context.Context, filter types.MetricID, tx Transaction) (*types.Metrics, bool)
}

type Transaction interface {
	Commit() error
	Rollback() error
}

type MetricUpdatesService struct {
	save        MetricUpdatesSaveRepository
	filter      MetricUpdatesFilterRepository
	transaction Transaction
}

// Метод обновления метрик с использованием транзакции
func (svc MetricUpdatesService) Updates(
	ctx context.Context, metrics []*types.Metrics,
) ([]*types.Metrics, error) {
	var updatedMetrics []*types.Metrics

	tx := svc.transaction

	defer func() {
		if err := tx.Rollback(); err != nil {
		}
	}()

	for _, metric := range metrics {
		existingMetric, found := svc.filter.Filter(ctx, types.MetricID{
			ID:   metric.MetricID.ID,
			Type: metric.MetricID.Type,
		}, tx)

		if !found {
			if ok := svc.save.Save(ctx, metric, tx); !ok {
				return nil, errors.ErrMetricInternal
			}
			updatedMetrics = append(updatedMetrics, metric)
			continue
		}

		switch metric.MetricID.Type {
		case types.Gauge:
			existingMetric.Value = metric.Value
		case types.Counter:
			*existingMetric.Delta += *metric.Delta
		}

		svc.save.Save(ctx, existingMetric, tx)

		updatedMetrics = append(updatedMetrics, existingMetric)
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.ErrMetricInternal
	}

	return updatedMetrics, nil
}
