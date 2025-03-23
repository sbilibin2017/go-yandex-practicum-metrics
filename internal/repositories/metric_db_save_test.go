package repositories

import (
	"context"
	"database/sql"
	"testing"

	"go-yandex-practicum-metrics/internal/types"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Тест успешного сохранения метрики
func TestMetricDBSaveRepository_Save_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMetricDBSaveRepository(db)

	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "test_metric", Type: "gauge"},
		Delta:    nil,
		Value:    new(float64),
	}
	*metric.Value = 42.0

	mock.ExpectExec(`INSERT INTO metrics`).
		WithArgs(metric.ID, metric.Type, metric.Delta, metric.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	success := repo.Save(context.Background(), metric)
	assert.True(t, success)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест обновления существующей метрики (ON CONFLICT DO UPDATE)
func TestMetricDBSaveRepository_Save_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMetricDBSaveRepository(db)

	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "existing_metric", Type: "counter"},
		Delta:    new(int64),
		Value:    nil,
	}
	*metric.Delta = 10

	mock.ExpectExec(`INSERT INTO metrics`).
		WithArgs(metric.ID, metric.Type, metric.Delta, metric.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	success := repo.Save(context.Background(), metric)
	assert.True(t, success)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест обработки ошибки при выполнении запроса
func TestMetricDBSaveRepository_Save_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMetricDBSaveRepository(db)

	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "fail_metric", Type: "gauge"},
		Delta:    nil,
		Value:    new(float64),
	}
	*metric.Value = 100.0

	mock.ExpectExec(`INSERT INTO metrics`).
		WithArgs(metric.ID, metric.Type, metric.Delta, metric.Value).
		WillReturnError(sql.ErrConnDone)

	success := repo.Save(context.Background(), metric)
	assert.False(t, success)
	assert.NoError(t, mock.ExpectationsWereMet())
}
