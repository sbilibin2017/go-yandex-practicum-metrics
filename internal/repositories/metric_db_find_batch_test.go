package repositories_test

import (
	"context"
	"fmt"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/repositories"
	"testing"

	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupFindBatchPostgresContainer(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	var db *sql.DB
	var retries = 5
	var connectionError error
	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)
	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)
	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port.Port())
	time.Sleep(2 * time.Second)
	for i := 0; i < retries; i++ {
		db, connectionError = sql.Open("pgx", dsn)
		if connectionError == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	require.NoError(t, connectionError)
	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS metrics (
		id VARCHAR NOT NULL,
		type VARCHAR NOT NULL,
		delta BIGINT,
		value DOUBLE PRECISION,
		PRIMARY KEY (id, type)
	);
	`)
	require.NoError(t, err)
	cleanup := func() {
		_, err := db.ExecContext(ctx, "DROP TABLE IF EXISTS metrics;")
		require.NoError(t, err)
		_ = db.Close()
		_ = postgresContainer.Terminate(ctx)
	}
	return db, cleanup
}

func TestMetricDBFindBatchRepository_FindBatch(t *testing.T) {
	db, cleanup := setupFindBatchPostgresContainer(t)
	defer cleanup()
	repo := repositories.NewMetricDBFindBatchRepository(db)
	_, err := db.ExecContext(context.Background(), `
	INSERT INTO metrics (id, type, delta, value) VALUES
	('metric1', 'counter', 10, NULL),
	('metric2', 'gauge', NULL, 20.5),
	('metric3', 'counter', 5, NULL)
	`)
	require.NoError(t, err)
	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
		{ID: "metric2", Type: "gauge"},
	}
	result, err := repo.FindBatch(context.Background(), filters)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	metric1, exists := result[domain.MetricID{ID: "metric1", Type: "counter"}]
	assert.True(t, exists)
	assert.Equal(t, int64(10), *metric1.Delta)
	assert.Nil(t, metric1.Value)
	metric2, exists := result[domain.MetricID{ID: "metric2", Type: "gauge"}]
	assert.True(t, exists)
	assert.Equal(t, 20.5, *metric2.Value)
	assert.Nil(t, metric2.Delta)
	_, exists = result[domain.MetricID{ID: "metric3", Type: "counter"}]
	assert.False(t, exists)
}
