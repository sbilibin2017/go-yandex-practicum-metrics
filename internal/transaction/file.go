package transaction

import (
	"context"
	"database/sql"
)

// FileTxBeginer реализует интерфейс TxBeginer для транзакций с файлами
type FileTxBeginer struct{}

// NewFileTxBeginer создает новый объект FileTxBeginer
func NewFileTxBeginer() *FileTxBeginer {
	return &FileTxBeginer{}
}

// BeginTx реализует метод интерфейса TxBeginer для начала транзакции с файлом
func (r *FileTxBeginer) BeginTx(ctx context.Context, opts *sql.TxOptions) (*FileTx, error) {
	return &FileTx{}, nil
}

// FileTx реализует интерфейс Tx для транзакций с файлами
type FileTx struct{}

// Commit реализует метод интерфейса Tx для фиксации транзакции с файлом
func (r *FileTx) Commit() error {
	// Здесь можно будет добавить логику сохранения данных в файл, если нужно
	return nil
}

// Rollback реализует метод интерфейса Tx для отката транзакции с файлом
func (r *FileTx) Rollback() error {
	// Здесь можно будет добавить логику отката изменений в файле, если нужно
	return nil
}
