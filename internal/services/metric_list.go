package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/logger"
)

type MetricListRepository interface {
	Find(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metrics, error)
}

type MetricListService struct {
	listRepo MetricListRepository
}

func NewMetricListService(
	listRepo MetricListRepository,
) *MetricListService {
	return &MetricListService{
		listRepo: listRepo,
	}
}

func (s *MetricListService) List(
	ctx context.Context,
) ([]*domain.Metrics, error) {
	logger.Info("Starting to list metrics")
	existingMetrics, err := s.listRepo.Find(ctx, []domain.MetricID{})
	if err != nil {
		logger.Error("Error finding metrics")
		return nil, errors.ErrInternal
	}
	var metrics []*domain.Metrics
	for _, metric := range existingMetrics {
		metrics = append(metrics, metric)
	}
	logger.Info("Metrics listed successfully", "metrics_count", len(metrics))
	return metrics, nil
}
