package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
)

type MetricListRepository interface {
	List(ctx context.Context) ([]*types.Metrics, error)
}

type MetricListService struct {
	list MetricListRepository
}

func NewMetricListService(
	list MetricListRepository,
) *MetricListService {
	return &MetricListService{
		list: list,
	}
}

// Метод для получения метрики по типу и ID
func (svc *MetricListService) List(
	ctx context.Context,
) ([]*types.Metrics, error) {
	metric, err := svc.list.List(ctx)
	if err != nil {
		return nil, errors.ErrMetricInternal
	}
	return metric, nil
}
