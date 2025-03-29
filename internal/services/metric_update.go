package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/logger"
)

type MetricUpdateSaveBatchRepository interface {
	Save(ctx context.Context, metrics []*domain.Metrics) error
}

type MetricUpdateFindBatchRepository interface {
	Find(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metrics, error)
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func() error) error
}

type MetricUpdateService struct {
	saveRepo MetricUpdateSaveBatchRepository
	findRepo MetricUpdateFindBatchRepository
	uow      UnitOfWork
}

func NewMetricUpdateService(
	saveRepo MetricUpdateSaveBatchRepository,
	findRepo MetricUpdateFindBatchRepository,
	uow UnitOfWork,
) *MetricUpdateService {
	return &MetricUpdateService{
		saveRepo: saveRepo,
		findRepo: findRepo,
		uow:      uow,
	}
}

func (s *MetricUpdateService) Update(
	ctx context.Context, metrics []*domain.Metrics,
) ([]*domain.Metrics, error) {
	var updatedMetrics []*domain.Metrics
	logger.Info("Starting metric update", "metrics_count", len(metrics))
	err := s.uow.Do(ctx, func() error {
		metricIDs := make([]domain.MetricID, len(metrics))
		for i, metric := range metrics {
			metricIDs[i] = domain.MetricID{ID: metric.ID, Type: metric.Type}
		}
		logger.Info("Finding existing metrics", "metric_ids", metricIDs)
		existingMetrics, err := s.findRepo.Find(ctx, metricIDs)
		if err != nil {
			logger.Error("Error finding existing metrics", "metric_ids", metricIDs)
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
					logger.Info("Updated counter metric delta", "metric_id", metric.ID, "new_delta", *metric.Delta)
				}
			}
			updatedMetrics[i] = metric
		}
		logger.Info("Saving updated metrics", "metrics_count", len(updatedMetrics))
		if err := s.saveRepo.Save(ctx, updatedMetrics); err != nil {
			logger.Error("Error saving updated metrics", "metrics_count", len(updatedMetrics))
			return errors.ErrInternal
		}
		return nil
	})
	if err != nil {
		logger.Error("Error during metric update", "error", err)
		return nil, err
	}
	logger.Info("Metric update completed", "metrics_count", len(updatedMetrics))
	return updatedMetrics, nil
}
