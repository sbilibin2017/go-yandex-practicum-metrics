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

type PostgresContainer2 struct {
	Container testcontainers.Container
	DB        *sql.DB
}

func setupPostgresContainer2(t *testing.T) *PostgresContainer {
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

func teardownPostgresContainer2(t *testing.T, pc *PostgresContainer) {
	if err := pc.DB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
	if err := pc.Container.Terminate(context.Background()); err != nil {
		log.Printf("Error stopping container: %v", err)
	}
}

func TestSaveBatch_Success(t *testing.T) {
	pc := setupPostgresContainer2(t)
	defer teardownPostgresContainer2(t, pc)
	repo := NewMetricDBSaveBatchRepository(pc.DB)
	metrics := []*domain.Metrics{
		{ID: "metric1", Type: "counter", Delta: new(int64), Value: nil},
		{ID: "metric2", Type: "gauge", Delta: nil, Value: new(float64)},
	}
	*metrics[0].Delta = 100
	*metrics[1].Value = 42.5
	success := repo.SaveBatch(context.Background(), metrics)
	assert.True(t, success, "SaveBatch should return true on successful insert")
	var count int
	err := pc.DB.QueryRow("SELECT COUNT(*) FROM metrics WHERE id = $1", "metric1").Scan(&count)
	if err != nil {
		t.Fatalf("Query execution error: %v", err)
	}
	assert.Equal(t, 1, count, "There should be 1 metric with ID metric1")
	err = pc.DB.QueryRow("SELECT COUNT(*) FROM metrics WHERE id = $1", "metric2").Scan(&count)
	if err != nil {
		t.Fatalf("Query execution error: %v", err)
	}
	assert.Equal(t, 1, count, "There should be 1 metric with ID metric2")
}
