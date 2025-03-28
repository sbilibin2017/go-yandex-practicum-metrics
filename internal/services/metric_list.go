package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
)

type MetricListRepository interface {
	Find(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metrics, bool)
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
	existingMetrics, ok := s.listRepo.Find(ctx, []domain.MetricID{})
	if !ok {
		return nil, errors.ErrInternal
	}
	var metrics []*domain.Metrics
	for _, metric := range existingMetrics {
		metrics = append(metrics, metric)
	}
	return metrics, nil
}
