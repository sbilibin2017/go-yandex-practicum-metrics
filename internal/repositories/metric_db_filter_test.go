package repositories

import (
	"context"
	"database/sql"
	"testing"

	"go-yandex-practicum-metrics/internal/types"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Тест успешного получения метрики
func TestMetricDBFilterRepository_Filter_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMetricDBFilterRepository(db)

	filter := types.MetricID{ID: "test_metric", Type: "gauge"}
	value := 42.0

	rows := sqlmock.NewRows([]string{"id", "type", "delta", "value"}).
		AddRow(filter.ID, filter.Type, nil, value)

	mock.ExpectQuery(`SELECT id, type, delta, value FROM metrics`).
		WithArgs(filter.ID, filter.Type).
		WillReturnRows(rows)

	result, found := repo.Filter(context.Background(), filter)
	assert.True(t, found)
	assert.NotNil(t, result)
	assert.Equal(t, filter.ID, result.ID)
	assert.Equal(t, filter.Type, result.Type)
	assert.Nil(t, result.Delta)
	assert.NotNil(t, result.Value)
	assert.Equal(t, value, *result.Value)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест, когда метрика не найдена (sql.ErrNoRows)
func TestMetricDBFilterRepository_Filter_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMetricDBFilterRepository(db)

	filter := types.MetricID{ID: "missing_metric", Type: "counter"}

	mock.ExpectQuery(`SELECT id, type, delta, value FROM metrics`).
		WithArgs(filter.ID, filter.Type).
		WillReturnError(sql.ErrNoRows)

	result, found := repo.Filter(context.Background(), filter)
	assert.False(t, found)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест обработки ошибки запроса к БД
func TestMetricDBFilterRepository_Filter_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMetricDBFilterRepository(db)

	filter := types.MetricID{ID: "error_metric", Type: "gauge"}

	mock.ExpectQuery(`SELECT id, type, delta, value FROM metrics`).
		WithArgs(filter.ID, filter.Type).
		WillReturnError(sql.ErrConnDone)

	result, found := repo.Filter(context.Background(), filter)
	assert.False(t, found)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}
