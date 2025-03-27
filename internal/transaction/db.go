package transaction

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // Подключаем PostgreSQL драйвер
)

// RealTxBeginer реализует интерфейс TxBeginer
type DBTxBeginer struct {
	db *sql.DB // Это реальное соединение с БД
}

// NewRealTxBeginer создает новый объект RealTxBeginer с подключением к базе данных
func NewDBTxBeginer(db *sql.DB) *DBTxBeginer {
	return &DBTxBeginer{
		db: db,
	}
}

// BeginTx реализует метод интерфейса TxBeginer для начала транзакции
func (r *DBTxBeginer) BeginTx(ctx context.Context, opts *sql.TxOptions) (*DBTx, error) {
	// Начинаем транзакцию с указанными параметрами
	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &DBTx{tx: tx}, nil
}

// RealTx реализует интерфейс Tx
type DBTx struct {
	tx *sql.Tx // Это реальная транзакция
}

// Commit реализует метод интерфейса Tx для фиксации транзакции
func (r *DBTx) Commit() error {
	// Фиксируем транзакцию
	if err := r.tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Rollback реализует метод интерфейса Tx для отката транзакции
func (r *DBTx) Rollback() error {
	// Откатываем транзакцию
	if err := r.tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}
