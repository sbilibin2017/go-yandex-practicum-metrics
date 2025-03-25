package app

import (
	"context"
	"go-yandex-practicum-metrics/internal/configs"
	"go-yandex-practicum-metrics/internal/logger"
	"go-yandex-practicum-metrics/internal/types"
	"strconv"
	"time"
)

type DumpMetricListRepository interface {
	List(ctx context.Context) ([]*types.Metrics, error)
}

type DumpUpdateSaveRepository interface {
	Save(ctx context.Context, metric *types.Metrics) error
}

type DumpMetricsToFileWorker struct {
	config *configs.ServerConfig
	list   DumpMetricListRepository
	update DumpUpdateSaveRepository
}

func NewDumpMetricsToFileWorker(
	config *configs.ServerConfig,
	list DumpMetricListRepository,
	update DumpUpdateSaveRepository,
) *DumpMetricsToFileWorker {
	return &DumpMetricsToFileWorker{
		config: config,
		list:   list,
		update: update,
	}
}

func (w *DumpMetricsToFileWorker) Start(ctx context.Context) {
	storeIntervalSec, err := strconv.Atoi(w.config.StoreInterval)
	if err != nil {
		logger.Logger.Error("Invalid StoreInterval value: ", err)
		return
	}

	storeInterval := time.Duration(storeIntervalSec) * time.Second

	logger.Logger.Info("Starting DumpMetricsToFileWorker with interval: ", storeInterval)

	ticker := time.NewTicker(storeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Logger.Info("Context canceled, stopping worker...")
			return
		case <-ticker.C:
			metrics, err := w.list.List(ctx)
			if err != nil {
				logger.Logger.Error("Error loading metrics: ", err)
				continue
			}

			for _, metric := range metrics {
				if err := w.update.Save(ctx, metric); err != nil {
					logger.Logger.Error("Error saving metric: ", err)
				}
			}
		}
	}
}
