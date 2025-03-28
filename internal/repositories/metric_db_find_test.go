package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"go-yandex-practicum-metrics/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	Container testcontainers.Container
	DB        *sql.DB
}

func setupPostgresContainer(t *testing.T) *PostgresContainer {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Error creating container: %v", err)
	}
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatalf("Database ping error: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE metrics (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			delta BIGINT NULL,
			value DOUBLE PRECISION NULL
		);
	`)
	if err != nil {
		t.Fatalf("Error creating table: %v", err)
	}
	_, err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_metrics_id_type ON metrics (id, type);
	`)
	if err != nil {
		t.Fatalf("Error creating index: %v", err)
	}
	return &PostgresContainer{
		Container: container,
		DB:        db,
	}
}

func teardownPostgresContainer(t *testing.T, pc *PostgresContainer) {
	if err := pc.DB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
	if err := pc.Container.Terminate(context.Background()); err != nil {
		log.Printf("Error stopping container: %v", err)
	}
}

func TestFindBatch_Success(t *testing.T) {
	pc := setupPostgresContainer(t)
	defer teardownPostgresContainer(t, pc)
	repo := NewMetricDBFindBatchRepository(pc.DB)
	_, err := pc.DB.Exec(`
		INSERT INTO metrics (id, type, delta, value) VALUES 
		('metric1', 'counter', 100, NULL),
		('metric2', 'gauge', NULL, 42.5),
		('metric3', 'counter', 200, NULL);
	`)
	assert.NoError(t, err)
	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
		{ID: "metric2", Type: "gauge"},
	}
	results, ok := repo.FindBatch(context.Background(), filters)
	assert.True(t, ok, "Should return true on success")
	assert.Len(t, results, 2, "Should find 2 metrics")
	metric1 := results[domain.MetricID{ID: "metric1", Type: "counter"}]
	assert.NotNil(t, metric1, "metric1 should not be nil")
	assert.Equal(t, int64(100), *metric1.Delta)
	metric2 := results[domain.MetricID{ID: "metric2", Type: "gauge"}]
	assert.NotNil(t, metric2, "metric2 should not be nil")
	assert.Equal(t, 42.5, *metric2.Value)
}

func TestFindBatch_QueryError(t *testing.T) {
	pc := setupPostgresContainer(t)
	defer teardownPostgresContainer(t, pc)
	_, err := pc.DB.Exec("DROP TABLE metrics;")
	assert.NoError(t, err)
	repo := NewMetricDBFindBatchRepository(pc.DB)
	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
	}
	results, ok := repo.FindBatch(context.Background(), filters)
	assert.False(t, ok, "Should return false on query error")
	assert.Nil(t, results, "Should return nil on query error")
}
