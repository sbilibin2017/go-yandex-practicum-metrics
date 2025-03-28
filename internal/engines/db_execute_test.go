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

func TestDBExecuteEngine_Execute(t *testing.T) {
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
	engine, err := NewDBExecuteEngine(db)
	if err != nil {
		t.Fatalf("failed to create DBExecuteEngine: %v", err)
	}
	createTableQuery := `CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT)`
	assert.True(t, engine.Execute(ctx, createTableQuery))
	insertQuery := `INSERT INTO users (name) VALUES ($1)`
	assert.True(t, engine.Execute(ctx, insertQuery, "John Doe"))
	var name string
	row := db.QueryRowContext(ctx, "SELECT name FROM users WHERE name = $1", "John Doe")
	err = row.Scan(&name)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", name)
	dropTableQuery := `DROP TABLE IF EXISTS users`
	assert.True(t, engine.Execute(ctx, dropTableQuery))
}
