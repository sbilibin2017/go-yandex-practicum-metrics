package transaction

import (
	"context"
	"database/sql"
)

// MemoryTxBeginer реализует интерфейс TxBeginer для транзакций в памяти
type MemoryTxBeginer struct{}

// NewMemoryTxBeginer создает новый объект MemoryTxBeginer
func NewMemoryTxBeginer() *MemoryTxBeginer {
	return &MemoryTxBeginer{}
}

// BeginTx реализует метод интерфейса TxBeginer для начала транзакции
func (r *MemoryTxBeginer) BeginTx(ctx context.Context, opts *sql.TxOptions) (*MemoryTx, error) {
	return &MemoryTx{}, nil
}

// MemoryTx реализует интерфейс Tx для транзакций в памяти
type MemoryTx struct{}

// Commit реализует метод интерфейса Tx для фиксации транзакции
func (r *MemoryTx) Commit() error {
	return nil
}

// Rollback реализует метод интерфейса Tx для отката транзакции
func (r *MemoryTx) Rollback() error {
	return nil
}
