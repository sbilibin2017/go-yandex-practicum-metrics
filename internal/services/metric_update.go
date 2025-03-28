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

type Tx interface {
	Commit() error
	Rollback() error
}

type WithTx interface {
	Begin(ctx context.Context) (Tx, error)
	Do(ctx context.Context, fn func(tx Tx) error) error
}

type MetricUpdateService struct {
	saveRepo MetricUpdateSaveBatchRepository
	findRepo MetricUpdateFindBatchRepository
	withTx   WithTx
}

func NewMetricUpdateService(
	saveRepo MetricUpdateSaveBatchRepository,
	findRepo MetricUpdateFindBatchRepository,
	withTx WithTx,
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
	var updatedMetrics []*domain.Metrics
	err := s.withTx.Do(ctx, func(tx Tx) error {
		metricIDs := make([]domain.MetricID, len(metrics))
		for i, metric := range metrics {
			metricIDs[i] = domain.MetricID{ID: metric.ID, Type: metric.Type}
		}
		existingMetrics, ok := s.findRepo.Find(ctx, metricIDs)
		if !ok {
			return errors.ErrInternal
		}
		updatedMetrics = make([]*domain.Metrics, len(metrics))
		for i, metric := range metrics {
			switch metric.Type {
			case domain.Counter:
				if existingMetric, exists := existingMetrics[domain.MetricID{
					ID:   metric.ID,
					Type: metric.Type,
				}]; exists {
					*metric.Delta += *existingMetric.Delta
				}
			}
			updatedMetrics[i] = metric
		}
		if ok := s.saveRepo.Save(ctx, updatedMetrics); !ok {
			return errors.ErrInternal
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updatedMetrics, nil
}
