package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
)

type MetricGetByTypeAndIDFilterRepository interface {
	Filter(ctx context.Context, filter types.MetricID) (*types.Metrics, bool)
}

type MetricGetByTypeAndIDService struct {
	filter MetricGetByTypeAndIDFilterRepository
}

// Метод для получения метрики по типу и ID
func (svc *MetricGetByTypeAndIDService) GetByTypeAndID(
	ctx context.Context, id types.MetricID,
) (*types.Metrics, error) {
	metric, found := svc.filter.Filter(ctx, id)
	if !found {
		return nil, errors.ErrMetricNotFound
	}
	return metric, nil
}
