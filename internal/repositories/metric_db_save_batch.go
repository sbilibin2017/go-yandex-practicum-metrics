package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-yandex-practicum-metrics/internal/domain"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type MetricDBSaveBatchRepository struct {
	db *sql.DB
}

var saveBatchQueryTemplate = `
INSERT INTO metrics (id, type, delta, value) 
VALUES %s
ON CONFLICT (id, type) 
DO UPDATE SET 
	delta = EXCLUDED.delta, 
	value = EXCLUDED.value
`

func buildSaveBatchQuery(metrics []*domain.Metrics) (string, []any) {
	var args []any
	placeholders := []string{}
	for i, metric := range metrics {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		args = append(args, metric.ID, metric.Type, metric.Delta, metric.Value)
	}
	query := fmt.Sprintf(saveBatchQueryTemplate, strings.Join(placeholders, ", "))
	return query, args
}

func NewMetricDBSaveBatchRepository(db *sql.DB) *MetricDBSaveBatchRepository {
	return &MetricDBSaveBatchRepository{db: db}
}

func (repo *MetricDBSaveBatchRepository) SaveBatch(
	ctx context.Context, metrics []*domain.Metrics,
) error {
	if len(metrics) == 0 {
		return nil
	}
	query, args := buildSaveBatchQuery(metrics)
	_, err := repo.db.ExecContext(ctx, query, args...)
	return err
}
