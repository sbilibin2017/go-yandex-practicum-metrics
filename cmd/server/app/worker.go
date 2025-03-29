package app

import (
	"context"
	"fmt"
	"go-yandex-practicum-metrics/internal/configs"
	"go-yandex-practicum-metrics/internal/domain"
	"strconv"
	"time"
)

type Worker struct {
	config    *configs.ServerConfig
	container *Container
}

func NewWorker(config *configs.ServerConfig, container *Container) (*Worker, error) {
	return &Worker{
		config:    config,
		container: container,
	}, nil
}

func (w *Worker) Start(ctx context.Context) error {
	if w.config.GetRestore() == "true" {
		err := w.loadMetricsFromFile(ctx)
		if err != nil {
			return fmt.Errorf("failed to restore data from file: %v", err)
		}
	}
	if w.config.GetStoreInterval() == "0" {
		<-ctx.Done()
		return w.saveMetricsToFile(ctx)
	}
	storeInterval, _ := strconv.Atoi(w.config.GetStoreInterval())
	ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return w.saveMetricsToFile(ctx)
		case <-ticker.C:
			err := w.saveMetricsToFile(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (w *Worker) saveMetricsToFile(ctx context.Context) error {
	var metrics map[domain.MetricID]*domain.Metrics
	var err error
	if w.config.GetDatabaseDSN() != "" {
		metrics, err = w.container.DBFindRepo.Find(ctx, []domain.MetricID{})
		if err != nil {
			return fmt.Errorf("failed to load data from db: %v", err)
		}
	} else {
		metrics, err = w.container.FileFindRepo.Find(ctx, []domain.MetricID{})
		if err != nil {
			return fmt.Errorf("failed to load data from file: %v", err)
		}

	}
	var data []*domain.Metrics
	for _, m := range metrics {
		data = append(data, m)
	}
	err = w.container.FileSaveRepo.Save(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to save data to file: %v", err)
	}
	return nil
}

func (w *Worker) loadMetricsFromFile(ctx context.Context) error {
	metrics, err := w.container.FileFindRepo.Find(ctx, []domain.MetricID{})
	if err != nil {
		return fmt.Errorf("failed to save data to file: %v", err)
	}
	var data []*domain.Metrics
	for _, m := range metrics {
		data = append(data, m)
	}
	if w.config.GetDatabaseDSN() != "" {
		err = w.container.DBSaveRepo.Save(ctx, data)
		if err != nil {
			return fmt.Errorf("failed to save data to file: %v", err)
		}
	} else {
		err = w.container.FileSaveRepo.Save(ctx, data)
		if err != nil {
			return fmt.Errorf("failed to save data to file: %v", err)
		}

	}
	return nil
}
