package services

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
)

type MetricListRepository interface {
	List(ctx context.Context) (map[domain.MetricID]*domain.Metrics, error)
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

var (
	ErrMetricListInternal = errors.New("internal error")
)

func (s *MetricListService) List(
	ctx context.Context,
) ([]*domain.Metrics, error) {
	existingMetrics, err := s.listRepo.List(ctx)
	if err != nil {
		return nil, ErrMetricListInternal
	}
	var metrics []*domain.Metrics
	for _, metric := range existingMetrics {
		metrics = append(metrics, metric)
	}
	return metrics, nil
}
