package services

import (
	"context"

	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
)

type MetricGetFilterRepository interface {
	Filter(ctx context.Context, filter types.MetricID) (*types.Metrics, error)
}

type MetricGetService struct {
	filter MetricGetFilterRepository
}

func NewMetricGetService(
	filter MetricGetFilterRepository,
) *MetricGetService {
	return &MetricGetService{
		filter: filter,
	}
}

// Метод для получения метрики по типу и ID
func (svc *MetricGetService) Get(
	ctx context.Context, id types.MetricID,
) (*types.Metrics, error) {
	metric, err := svc.filter.Filter(ctx, id)
	if err != nil {
		return nil, errors.ErrMetricInternal
	}
	if metric == nil {
		return nil, errors.ErrMetricNotFound
	}
	return metric, nil
}
