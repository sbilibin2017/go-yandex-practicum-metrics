package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-yandex-practicum-metrics/internal/domain"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type MetricDBFindBatchRepository struct {
	db *sql.DB
}

func NewMetricDBFindBatchRepository(db *sql.DB) *MetricDBFindBatchRepository {
	return &MetricDBFindBatchRepository{db: db}
}

var findBatchQueryTemplate = "SELECT id, type, delta, value FROM metrics WHERE %s"

func buildFindQuery(filters []domain.MetricID) (string, []any) {
	var args []any
	conditions := []string{}
	for i, filter := range filters {
		conditions = append(conditions, fmt.Sprintf("(id = $%d AND type = $%d)", i*2+1, i*2+2))
		args = append(args, filter.ID, filter.Type)
	}
	query := fmt.Sprintf(findBatchQueryTemplate, strings.Join(conditions, " OR "))
	return query, args
}

func (repo *MetricDBFindBatchRepository) FindBatch(
	ctx context.Context, filters []domain.MetricID,
) (map[domain.MetricID]*domain.Metrics, error) {
	if len(filters) == 0 {
		return nil, nil
	}
	query, args := buildFindQuery(filters)
	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make(map[domain.MetricID]*domain.Metrics)
	for rows.Next() {
		var metric domain.Metrics
		if err := rows.Scan(&metric.ID, &metric.Type, &metric.Delta, &metric.Value); err != nil {
			return nil, err
		}
		results[domain.MetricID{ID: metric.ID, Type: metric.Type}] = &metric
	}
	return results, nil
}
