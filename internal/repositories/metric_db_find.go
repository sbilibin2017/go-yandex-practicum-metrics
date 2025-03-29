package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/logger"
	"strings"
)

type DBFind interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type MetricDBFindRepository struct {
	db DBFind
}

func NewMetricDBFindRepository(db DBFind) *MetricDBFindRepository {
	return &MetricDBFindRepository{db: db}
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

func (repo *MetricDBFindRepository) Find(
	ctx context.Context, filters []domain.MetricID,
) (map[domain.MetricID]*domain.Metrics, error) {
	query, args := buildFindQuery(filters)
	logger.Info("Executing batch find query", "query", query, "args_count", len(args))
	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("Error executing find batch query", "error", err)
		return nil, err
	}
	defer rows.Close()
	results := make(map[domain.MetricID]*domain.Metrics)
	for rows.Next() {
		var m domain.Metrics
		if err := rows.Scan(&m.ID, &m.Type, &m.Delta, &m.Value); err != nil {
			logger.Error("Error scanning row", "error", err)
			return nil, err
		}
		results[domain.MetricID{ID: m.ID, Type: m.Type}] = &m
	}
	if err := rows.Err(); err != nil {
		logger.Error("Error during rows iteration", "error", err)
		return nil, err
	}
	logger.Info("Batch find completed successfully", "results_count", len(results))
	return results, nil
}
