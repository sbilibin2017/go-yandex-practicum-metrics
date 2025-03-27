package repositories_test

import (
	"context"
	"database/sql"
	"fmt"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/repositories"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupSaveBatchPostgresContainer(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second), // Waiting for the port to be open
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
		CONSTRAINT unique_id_type UNIQUE (id, type)  -- Unique constraint for the combination
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

func TestMetricDBSaveBatchRepository_SaveBatch(t *testing.T) {
	db, cleanup := setupSaveBatchPostgresContainer(t)
	defer cleanup()
	repo := repositories.NewMetricDBSaveBatchRepository(db)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Delta: ptrInt64(10),
		},
		{
			ID:    "metric2",
			Type:  "gauge",
			Value: ptrFloat64(20.5),
		},
	}
	err := repo.SaveBatch(context.Background(), metrics)
	require.NoError(t, err)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM metrics WHERE id = $1", "metric1").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	err = db.QueryRow("SELECT COUNT(*) FROM metrics WHERE id = $1", "metric2").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	var delta int64
	err = db.QueryRow("SELECT delta FROM metrics WHERE id = $1", "metric1").Scan(&delta)
	require.NoError(t, err)
	assert.Equal(t, int64(10), delta)
	var value float64
	err = db.QueryRow("SELECT value FROM metrics WHERE id = $1", "metric2").Scan(&value)
	require.NoError(t, err)
	assert.Equal(t, 20.5, value)
}
