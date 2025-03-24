package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
)

type MetricListAllRepository interface {
	ListAll(ctx context.Context) ([]*types.Metrics, bool)
}

type MetricListAllService struct {
	list MetricListAllRepository
}

// Метод для получения метрики по типу и ID
func (svc *MetricListAllService) ListAll(
	ctx context.Context,
) ([]*types.Metrics, error) {
	metric, ok := svc.list.ListAll(ctx)
	if !ok {
		return nil, errors.ErrMetricInternal
	}
	return metric, nil
}
