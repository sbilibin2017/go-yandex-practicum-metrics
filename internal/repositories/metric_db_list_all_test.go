package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Тест успешного получения всех метрик
func TestMetricDBListAllRepository_ListAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMetricDBListAllRepository(db)

	rows := sqlmock.NewRows([]string{"id", "type", "delta", "value"}).
		AddRow("metric_1", "gauge", nil, 42.0).
		AddRow("metric_2", "counter", int64(10), nil)

	mock.ExpectQuery(`SELECT id, type, delta, value FROM metrics`).
		WillReturnRows(rows)

	result, success := repo.ListAll(context.Background())
	assert.True(t, success)
	assert.Len(t, result, 2)

	assert.Equal(t, "metric_1", result[0].ID)
	assert.Equal(t, "gauge", result[0].Type)
	assert.Nil(t, result[0].Delta)
	assert.NotNil(t, result[0].Value)
	assert.Equal(t, 42.0, *result[0].Value)

	assert.Equal(t, "metric_2", result[1].ID)
	assert.Equal(t, "counter", result[1].Type)
	assert.NotNil(t, result[1].Delta)
	assert.Nil(t, result[1].Value)
	assert.Equal(t, int64(10), *result[1].Delta)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест, если в БД нет метрик (пустой результат)
func TestMetricDBListAllRepository_ListAll_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMetricDBListAllRepository(db)

	rows := sqlmock.NewRows([]string{"id", "type", "delta", "value"})

	mock.ExpectQuery(`SELECT id, type, delta, value FROM metrics`).
		WillReturnRows(rows)

	result, success := repo.ListAll(context.Background())
	assert.True(t, success)
	assert.Empty(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест обработки ошибки при запросе к БД
func TestMetricDBListAllRepository_ListAll_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMetricDBListAllRepository(db)

	mock.ExpectQuery(`SELECT id, type, delta, value FROM metrics`).
		WillReturnError(sql.ErrConnDone)

	result, success := repo.ListAll(context.Background())
	assert.False(t, success)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест ошибки при rows.Scan
func TestMetricDBListAllRepository_ListAll_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMetricDBListAllRepository(db)

	rows := sqlmock.NewRows([]string{"id", "type", "delta", "value"}).
		AddRow("metric_1", "gauge", nil, "invalid_float") // Ошибка: значение value некорректное

	mock.ExpectQuery(`SELECT id, type, delta, value FROM metrics`).
		WillReturnRows(rows)

	result, success := repo.ListAll(context.Background())
	assert.False(t, success)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест для обработки ошибки в rows.Err()
func TestMetricDBListAllRepository_ListAll_RowsError(t *testing.T) {
	// Создаем новый mock для базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMetricDBListAllRepository(db)

	// Создаем mock-строки для выполнения запроса
	rows := sqlmock.NewRows([]string{"id", "type", "delta", "value"}).
		AddRow("metric_1", "gauge", nil, 42.0).
		AddRow("metric_2", "counter", nil, 24.0)

	// Ожидаем выполнения запроса с возвратом mock-строк
	mock.ExpectQuery(`SELECT id, type, delta, value FROM metrics`).
		WillReturnRows(rows)

	// Симулируем ошибку в методе `rows.Err()`
	rows.RowError(0, sql.ErrConnDone) // Симуляция ошибки для первой строки

	// Вызов ListAll должен завершиться с ошибкой из-за проблемы с соединением
	result, success := repo.ListAll(context.Background())

	// Проверяем, что результат - это ошибка, и возвращается nil
	assert.False(t, success)
	assert.Nil(t, result)

	// Проверяем, что все ожидания были выполнены
	assert.NoError(t, mock.ExpectationsWereMet())
}
