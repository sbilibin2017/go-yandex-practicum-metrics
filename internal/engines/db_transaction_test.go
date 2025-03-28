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

func TestDBTransactionEngine_BeginCommitRollback(t *testing.T) {
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
	createTableQuery := `CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT)`
	_, err = db.ExecContext(ctx, createTableQuery)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	engine := &DBTransactionEngine{db: db}
	assert.True(t, engine.Begin(ctx))
	insertQuery := `INSERT INTO users (name) VALUES ($1)`
	_, err = engine.tx.ExecContext(ctx, insertQuery, "John Doe")
	assert.NoError(t, err)
	assert.True(t, engine.Commit())
	var name string
	row := db.QueryRowContext(ctx, "SELECT name FROM users WHERE name = $1", "John Doe")
	err = row.Scan(&name)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", name)
	assert.True(t, engine.Begin(ctx))
	insertQuery2 := `INSERT INTO users (name) VALUES ($1)`
	_, err = engine.tx.ExecContext(ctx, insertQuery2, "Jane Doe")
	assert.NoError(t, err)
	assert.True(t, engine.Rollback())
	row = db.QueryRowContext(ctx, "SELECT name FROM users WHERE name = $1", "Jane Doe")
	err = row.Scan(&name)
	assert.Error(t, err)
}
