package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-yandex-practicum-metrics/internal/domain"
)

type MetricUpdateSaveBatchRepository interface {
	SaveBatch(ctx context.Context, metrics []*domain.Metrics) error
}

type MetricUpdateFindBatchRepository interface {
	FindBatch(ctx context.Context, filters []domain.MetricID) (map[domain.MetricID]*domain.Metrics, error)
}

type TxBeginer interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
}

type Tx interface {
	Commit() error
	Rollback() error
}

type MetricUpdateService struct {
	saveRepo  MetricUpdateSaveBatchRepository
	findRepo  MetricUpdateFindBatchRepository
	txBeginer TxBeginer
}

func NewMetricUpdateService(
	saveRepo MetricUpdateSaveBatchRepository,
	findRepo MetricUpdateFindBatchRepository,
	txBeginer TxBeginer,
) *MetricUpdateService {
	return &MetricUpdateService{
		saveRepo:  saveRepo,
		findRepo:  findRepo,
		txBeginer: txBeginer,
	}
}

var (
	ErrMetricUpdateInternal = errors.New("internal error")
)

func (s *MetricUpdateService) UpdateBatch(
	ctx context.Context, metrics []*domain.Metrics,
) ([]*domain.Metrics, error) {
	tx, err := s.txBeginer.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	metricIDs := make([]domain.MetricID, len(metrics))
	for i, metric := range metrics {
		metricIDs[i] = domain.MetricID{ID: metric.ID, Type: metric.Type}
	}
	existingMetrics, err := s.findRepo.FindBatch(ctx, metricIDs)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to find batch: %w", err)
	}
	for i, metric := range metrics {
		if metric.Type == string(domain.Gauge) {
			metrics[i] = metric
		} else if metric.Type == string(domain.Counter) {
			if existingMetric, ok := existingMetrics[domain.MetricID{ID: metric.ID, Type: metric.Type}]; ok {
				if metric.Delta != nil {
					*existingMetric.Delta += *metric.Delta
				}
				metrics[i] = existingMetric
			} else {
				metrics[i] = metric
			}
		}
	}
	if err := s.saveRepo.SaveBatch(ctx, metrics); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to save batch: %w", err)
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return metrics, nil
}
