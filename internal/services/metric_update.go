package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
)

type MetricUpdateSaveBatchRepository interface {
	Save(ctx context.Context, metrics []*domain.Metrics) bool
}

type MetricUpdateFindBatchRepository interface {
	Find(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metrics, bool)
}

type WithTransaction interface {
	Do(ctx context.Context, f func() error) error
}

type MetricUpdateService struct {
	saveRepo MetricUpdateSaveBatchRepository
	findRepo MetricUpdateFindBatchRepository
	withTx   WithTransaction
}

func NewMetricUpdateService(
	saveRepo MetricUpdateSaveBatchRepository,
	findRepo MetricUpdateFindBatchRepository,
	withTx WithTransaction,
) *MetricUpdateService {
	return &MetricUpdateService{
		saveRepo: saveRepo,
		findRepo: findRepo,
		withTx:   withTx,
	}
}

func (s *MetricUpdateService) Update(
	ctx context.Context, metrics []*domain.Metrics,
) ([]*domain.Metrics, error) {
	err := s.withTx.Do(ctx, func() error {
		metricIDs := make([]domain.MetricID, len(metrics))
		for i, metric := range metrics {
			metricIDs[i] = domain.MetricID{ID: metric.ID, Type: metric.Type}
		}
		existingMetrics, ok := s.findRepo.Find(ctx, metricIDs)
		if !ok {
			return errors.ErrInternal
		}
		for i, metric := range metrics {
			switch metric.Type {
			case domain.Counter:
				if existingMetric, ok := existingMetrics[domain.MetricID{ID: metric.ID, Type: metric.Type}]; ok {
					*existingMetric.Delta += *metric.Delta
				}
				metrics[i] = metric
			case domain.Gauge:
				metrics[i] = metric
			}
		}
		if ok := s.saveRepo.Save(ctx, metrics); !ok {
			return errors.ErrInternal
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
