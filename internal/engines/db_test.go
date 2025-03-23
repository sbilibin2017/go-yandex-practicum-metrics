package engines

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Подключаем драйвер pgx для PostgreSQL
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

// Структура, реализующая интерфейс DatabaseDSNConfig
type DatabaseConfig struct {
	dsn string
}

func (d *DatabaseConfig) GetDatabaseDSN() string {
	return d.dsn
}

// TestNewDBStorage тестирует функцию NewDBStorage с использованием Testcontainers для PostgreSQL
func TestNewDBStorage(t *testing.T) {
	// Создаем контейнер с PostgreSQL с использованием testcontainers
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14", // Можно изменить версию PostgreSQL
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
	}

	// Запускаем контейнер
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer func() {
		// Завершаем контейнер после выполнения теста
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err)
	}()

	// Получаем хост и порт контейнера
	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Формируем DSN с использованием динамически полученного порта
	dsn := fmt.Sprintf("postgres://postgres:password@%s:%s/testdb?sslmode=disable", host, port.Port())

	// Создаем структуру DatabaseConfig, которая реализует интерфейс DatabaseDSNConfig
	config := &DatabaseConfig{dsn: dsn}

	// Подключаемся к базе данных с использованием DSN
	var db *sql.DB
	var errDB error
	for i := 0; i < 5; i++ { // Пробуем подключиться до 5 раз
		db, errDB = sql.Open("pgx", dsn) // Используем драйвер pgx
		if errDB == nil {
			// Пробуем пинговать базу данных, чтобы убедиться, что она доступна
			errDB = db.Ping()
			if errDB == nil {
				break
			}
		}
		time.Sleep(2 * time.Second) // Ждем перед повторной попыткой
	}

	// Проверяем, что подключение было успешным
	require.NoError(t, errDB)

	// Создаем новый экземпляр DBStorage с настоящим DSN
	storage, ok := NewDBStorage(config)
	require.True(t, ok)

	// Проверяем, что соединение с базой данных установлено
	assert.NotNil(t, storage.DB)
}
