package engines

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DBConfig is a generic interface for database configuration
type DBConfig interface {
	GetDatabaseDSN() string
}

type DBEngine struct {
	config DBConfig
	*sql.DB
}

func NewDBEngine(config DBConfig) *DBEngine {
	return &DBEngine{config: config}
}

// Open method accepts a DBConfig interface to support any configuration structure
func (db *DBEngine) Open() error {
	if db.DB != nil {
		return nil
	}
	conn, err := sql.Open("pgx", db.config.GetDatabaseDSN())
	if err != nil {
		return err
	}
	db.DB = conn
	return nil
}
