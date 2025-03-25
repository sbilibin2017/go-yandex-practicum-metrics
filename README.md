### Комментарии

пр с предыдущим инкрементом https://github.com/sbilibin2017/go-metrics-alerting/pull/14

приступил к рефакторингу в соответвии с комментариями


есть вопрос по поводу транзакции(как будет правильнее передать ее в сервисный слой?)
в текущей реализации проверка на транзакцию выглядит очень котсыльно



```
package services

import (
	"context"
	"database/sql"

	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
)

type Transaction interface {
	Begin() (*sql.Tx, error)
}

type MetricUpdatesSaveRepository interface {
	Save(ctx context.Context, metric *types.Metrics) error
}

type MetricUpdatesFilterRepository interface {
	Filter(ctx context.Context, filter types.MetricID) (*types.Metrics, error)
}

type MetricUpdatesService struct {
	save        MetricUpdatesSaveRepository
	filter      MetricUpdatesFilterRepository
	transaction Transaction
}

func NewMetricUpdatesService(
	save MetricUpdatesSaveRepository,
	filter MetricUpdatesFilterRepository,
	transaction Transaction,
) *MetricUpdatesService {
	return &MetricUpdatesService{
		save:        save,
		filter:      filter,
		transaction: transaction,
	}
}

func (svc *MetricUpdatesService) Updates(
	ctx context.Context, metrics []*types.Metrics,
) ([]*types.Metrics, error) {
	var updatedMetrics []*types.Metrics
	var tx *sql.Tx
	var err error

	// Если транзакция поддерживается — начинаем её
	if svc.transaction != nil {
		tx, err = svc.transaction.Begin()
		if err != nil {
			return nil, errors.ErrMetricInternal
		}
		defer func() {
			if err := tx.Rollback(); err != nil {
			}
		}()
	}

	for _, metric := range metrics {
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
			updatedMetrics = append(updatedMetrics, metric)
			continue
		}

		switch metric.MetricID.Type {
		case types.Gauge:
			existingMetric.Value = metric.Value
		case types.Counter:
			*existingMetric.Delta += *metric.Delta
		}

		if err := svc.save.Save(ctx, existingMetric); err != nil {
			return nil, errors.ErrMetricInternal
		}

		updatedMetrics = append(updatedMetrics, existingMetric)
	}

	if tx != nil {
		if err := tx.Commit(); err != nil {
			return nil, errors.ErrMetricInternal
		}
	}

	return updatedMetrics, nil
}
```