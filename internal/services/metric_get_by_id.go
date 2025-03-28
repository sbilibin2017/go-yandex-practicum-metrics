package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
)

type MetricGetFindBatchRepository interface {
	Find(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metrics, bool)
}

type MetricGetService struct {
	findRepo MetricGetFindBatchRepository
}

func NewMetricGetService(
	findRepo MetricGetFindBatchRepository,
) *MetricGetService {
	return &MetricGetService{
		findRepo: findRepo,
	}
}

func (s *MetricGetService) GetByID(
	ctx context.Context, id *domain.MetricID,
) (*domain.Metrics, error) {
	existingMetrics, ok := s.findRepo.Find(ctx, []domain.MetricID{*id})
	if !ok {
		return nil, errors.ErrInternal
	}
	metrics, exists := existingMetrics[*id]
	if !exists {
		return nil, errors.ErrMetricNotFound
	}
	return metrics, nil
}
