package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/logger"
)

type MetricGetFindBatchRepository interface {
	Find(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metrics, error)
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
	logger.Info("Fetching metric by ID", "id", id)
	existingMetrics, err := s.findRepo.Find(ctx, []domain.MetricID{*id})
	if err != nil {
		logger.Error("Error finding metrics for ID", "id", id)
		return nil, errors.ErrInternal
	}
	metrics, exists := existingMetrics[*id]
	if !exists {
		logger.Warn("Metric not found", "id", id)
		return nil, errors.ErrMetricNotFound
	}
	logger.Info("Metric found", "id", id)
	return metrics, nil
}
