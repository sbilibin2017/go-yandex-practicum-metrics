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

func TestDBScannerEngine_Scan(t *testing.T) {
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
	engine, err := NewDBScannerEngine(db)
	if err != nil {
		t.Fatalf("failed to create DBScannerEngine: %v", err)
	}
	createTableQuery := `CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, age INT)`
	_, success := engine.Scan(ctx, createTableQuery)
	assert.True(t, success)
	insertQuery := `INSERT INTO users (name, age) VALUES ($1, $2)`
	_, success = engine.Scan(ctx, insertQuery, "John Doe", 30)
	assert.True(t, success)
	selectQuery := `SELECT id, name, age FROM users WHERE name = $1`
	results, success := engine.Scan(ctx, selectQuery, "John Doe")
	assert.True(t, success)
	assert.Len(t, results, 1)
	resultRow := results[0].(map[string]any)
	assert.Equal(t, int64(1), resultRow["id"])
	assert.Equal(t, "John Doe", resultRow["name"])
	assert.Equal(t, int64(30), resultRow["age"])
	dropTableQuery := `DROP TABLE IF EXISTS users`
	_, success = engine.Scan(ctx, dropTableQuery)
	assert.True(t, success)
}
