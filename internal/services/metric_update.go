package services

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
)

type MetricSaveBatchRepository interface {
	SaveBatch(ctx context.Context, metrics []*domain.Metrics) error
}

type MetricFindBatchRepository interface {
	FindBatch(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metrics, error)
}

type Transaction interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
}

type MetricUpdateService struct {
	saveRepo MetricSaveBatchRepository
	findRepo MetricFindBatchRepository
	tx       Transaction
}

func NewMetricUpdateService(
	saveRepo MetricSaveBatchRepository,
	findRepo MetricFindBatchRepository,
	tx Transaction,
) *MetricUpdateService {
	return &MetricUpdateService{
		saveRepo: saveRepo,
		findRepo: findRepo,
		tx:       tx,
	}
}

var (
	ErrMetricUpdateInternal = errors.New("internal error")
)

func (s *MetricUpdateService) UpdateBatch(
	ctx context.Context, metrics []*domain.Metrics,
) ([]*domain.Metrics, error) {
	result, err := s.tx.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		metricIDs := make([]domain.MetricID, len(metrics))
		for i, metric := range metrics {
			metricIDs[i] = domain.MetricID{ID: metric.ID, Type: metric.Type}
		}
		existingMetrics, err := s.findRepo.FindBatch(ctx, metricIDs)
		if err != nil {
			return nil, ErrMetricUpdateInternal
		}
		metricsToSave := make([]*domain.Metrics, 0, len(metrics))
		for _, metric := range metrics {
			switch metric.Type {
			case string(domain.Gauge):
				metricsToSave = append(metricsToSave, metric)
			case string(domain.Counter):
				if existingMetric, ok := existingMetrics[domain.MetricID{ID: metric.ID, Type: metric.Type}]; ok {
					*existingMetric.Delta += *metric.Delta
					metricsToSave = append(metricsToSave, existingMetric)
				} else {
					metricsToSave = append(metricsToSave, metric)
				}
			}
		}
		if err := s.saveRepo.SaveBatch(ctx, metricsToSave); err != nil {
			return nil, ErrMetricUpdateInternal
		}
		return metricsToSave, nil
	})
	if err != nil {
		return nil, ErrMetricUpdateInternal
	}
	metricsToReturn, ok := result.([]*domain.Metrics)
	if !ok {
		return nil, ErrMetricUpdateInternal
	}
	return metricsToReturn, nil
}
