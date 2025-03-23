package engines

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib" // Correct import
)

// FileStorageConfig интерфейс для получения пути к файлу
type DatabaseDSNConfig interface {
	GetDatabaseDSN() string
}

// DBStorage - структура для работы с базой данных
type DBStorage struct {
	*sql.DB
}

// NewDBStorage - конструктор для создания новой структуры DBStorage
func NewDBStorage(c DatabaseDSNConfig) (*DBStorage, bool) {
	db, err := sql.Open("pgx", c.GetDatabaseDSN())
	if err != nil {
		return nil, false
	}

	if err := db.Ping(); err != nil {
		return nil, false
	}

	return &DBStorage{DB: db}, true
}
