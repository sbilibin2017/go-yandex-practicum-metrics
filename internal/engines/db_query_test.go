package engines

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestDBQueryEngine_Query(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}
	defer container.Terminate(ctx)
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}
	dsn := "user=user password=password dbname=testdb host=localhost port=" + port.Port() + " sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()
	engine, err := NewDBQueryEngine(db)
	if err != nil {
		t.Fatalf("failed to create DBQueryEngine: %v", err)
	}
	createTableQuery := `CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT)`
	_, success := engine.Query(ctx, createTableQuery)
	assert.True(t, success)
	insertQuery := `INSERT INTO users (name) VALUES ($1)`
	_, success = engine.Query(ctx, insertQuery, "John Doe")
	assert.True(t, success)
	selectQuery := `SELECT name FROM users WHERE name = $1`
	results, success := engine.Query(ctx, selectQuery, "John Doe")
	assert.True(t, success)
	assert.Len(t, results, 1)
	assert.Equal(t, "John Doe", results[0])
	dropTableQuery := `DROP TABLE IF EXISTS users`
	_, success = engine.Query(ctx, dropTableQuery)
	assert.True(t, success)
}
