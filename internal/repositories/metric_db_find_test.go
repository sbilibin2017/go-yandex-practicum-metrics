package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestBuildFindQuery(t *testing.T) {
	tests := []struct {
		name          string
		filters       []domain.MetricID
		expectedQuery string
		expectedArgs  []any
	}{
		{
			name:          "No filters",
			filters:       []domain.MetricID{},
			expectedQuery: "SELECT id, type, delta, value FROM metrics WHERE ",
			expectedArgs:  []any{},
		},
		{
			name: "Single filter",
			filters: []domain.MetricID{
				{ID: "metric1", Type: "counter"},
			},
			expectedQuery: "SELECT id, type, delta, value FROM metrics WHERE (id = $1 AND type = $2)",
			expectedArgs:  []any{"metric1", "counter"},
		},
		{
			name: "Multiple filters",
			filters: []domain.MetricID{
				{ID: "metric1", Type: "counter"},
				{ID: "metric2", Type: "gauge"},
			},
			expectedQuery: "SELECT id, type, delta, value FROM metrics WHERE (id = $1 AND type = $2) OR (id = $3 AND type = $4)",
			expectedArgs:  []any{"metric1", "counter", "metric2", "gauge"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, args := buildFindQuery(tt.filters)
			assert.Equal(t, tt.expectedQuery, query)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestFindBatch_Success(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать mock для базы данных: %v", err)
	}
	defer db.Close()
	repo := NewMetricDBFindBatchRepository(db)
	filters := []domain.MetricID{
		{ID: "metric1", Type: domain.Counter},
		{ID: "metric2", Type: domain.Gauge},
	}
	rows := sqlmock.NewRows([]string{"id", "type", "delta", "value"}).
		AddRow("metric1", "counter", 100, 1.23).
		AddRow("metric2", "gauge", 200, 2.34)
	mock.ExpectQuery("SELECT id, type, delta, value FROM metrics WHERE").
		WithArgs("metric1", "counter", "metric2", "gauge").
		WillReturnRows(rows)
	result, ok := repo.FindBatch(context.Background(), filters)
	assert.True(t, ok)
	assert.Len(t, result, 2)
	assert.Equal(t, &domain.Metrics{
		ID:    "metric1",
		Type:  domain.Counter,
		Delta: int64Pointer(100),
		Value: float64Pointer(1.23),
	}, result[domain.MetricID{ID: "metric1", Type: domain.Counter}])

	assert.Equal(t, &domain.Metrics{
		ID:    "metric2",
		Type:  domain.Gauge,
		Delta: int64Pointer(200),
		Value: float64Pointer(2.34),
	}, result[domain.MetricID{ID: "metric2", Type: domain.Gauge}])
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Невыполненные ожидания: %v", err)
	}
}

func int64Pointer(i int64) *int64 {
	return &i
}

func TestFindBatch_Error_QueryContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать mock для базы данных: %v", err)
	}
	defer db.Close()
	repo := NewMetricDBFindBatchRepository(db)
	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
	}
	mock.ExpectQuery("SELECT id, type, delta, value FROM metrics WHERE").
		WithArgs("metric1", "counter").
		WillReturnError(sql.ErrConnDone)
	result, ok := repo.FindBatch(context.Background(), filters)
	assert.False(t, ok)
	assert.Nil(t, result)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Невыполненные ожидания: %v", err)
	}
}

func TestFindBatch_Error_Scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Не удалось создать mock для базы данных: %v", err)
	}
	defer db.Close()
	repo := NewMetricDBFindBatchRepository(db)
	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
	}
	mock.ExpectQuery("SELECT id, type, delta, value FROM metrics WHERE").
		WithArgs("metric1", "counter").
		WillReturnRows(sqlmock.NewRows([]string{"id", "type", "delta", "value"}).
			AddRow("metric1", "counter", nil, nil)).
		WillReturnError(fmt.Errorf("Scan error"))
	result, ok := repo.FindBatch(context.Background(), filters)
	assert.False(t, ok)
	assert.Nil(t, result)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Невыполненные ожидания: %v", err)
	}
}
