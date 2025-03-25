package app

import (
	"context"
	"go-yandex-practicum-metrics/internal/logger"
	"go-yandex-practicum-metrics/internal/types"
)

type LoadMetricListRepository interface {
	List(ctx context.Context) ([]*types.Metrics, error)
}

type LoadMetricUpdateSaveRepository interface {
	Save(ctx context.Context, metric *types.Metrics) error
}

type LoadMetricsFromFileWorker struct {
	list   LoadMetricListRepository
	update LoadMetricUpdateSaveRepository
}

func NewLoadMetricsFromFileWorker(list LoadMetricListRepository, update LoadMetricUpdateSaveRepository) *LoadMetricsFromFileWorker {
	return &LoadMetricsFromFileWorker{
		list:   list,
		update: update,
	}
}

func (w *LoadMetricsFromFileWorker) Start(ctx context.Context) {
	logger.Logger.Info("Starting LoadMetricsFromFileWorker...")
	metrics, err := w.list.List(ctx)
	if err != nil {
		logger.Logger.Error("Error loading metrics: ", err)
		return
	}

	for _, metric := range metrics {
		if err := w.update.Save(ctx, metric); err != nil {
			logger.Logger.Error("Error saving metric: ", err)
		}
	}
}
