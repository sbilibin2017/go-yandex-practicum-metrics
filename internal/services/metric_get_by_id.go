package services

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
)

type MetricGetFindBatchRepository interface {
	FindBatch(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metrics, error)
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

var (
	ErrMetricGetInternal = errors.New("internal error")
	ErrMetricNotFound    = errors.New("internal error")
)

func (s *MetricGetService) GetByID(
	ctx context.Context, id *domain.MetricID,
) (*domain.Metrics, error) {
	existingMetrics, err := s.findRepo.FindBatch(ctx, []domain.MetricID{*id})
	if err != nil {
		return nil, ErrMetricGetInternal
	}
	metrics, exists := existingMetrics[*id]
	if !exists {
		return nil, ErrMetricNotFound
	}
	return metrics, nil
}
