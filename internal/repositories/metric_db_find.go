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
	if len(filters) > 0 {
		for i, filter := range filters {
			conditions = append(conditions, fmt.Sprintf("(id = $%d AND type = $%d)", i*2+1, i*2+2))
			args = append(args, filter.ID, string(filter.Type))
		}
	} else {
		args = []any{}
	}
	query := fmt.Sprintf(findBatchQueryTemplate, strings.Join(conditions, " OR "))
	return query, args
}

func (repo *MetricDBFindBatchRepository) FindBatch(
	ctx context.Context, filters []domain.MetricID,
) (map[domain.MetricID]*domain.Metrics, bool) {
	query, args := buildFindQuery(filters)
	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, false
	}
	defer rows.Close()
	results := make(map[domain.MetricID]*domain.Metrics)
	for rows.Next() {
		var m domain.Metrics
		if err := rows.Scan(&m.ID, &m.Type, &m.Delta, &m.Value); err != nil {
			return nil, false
		}
		results[domain.MetricID{ID: m.ID, Type: m.Type}] = &m
	}
	if err := rows.Err(); err != nil {
		return nil, false
	}
	return results, true
}
