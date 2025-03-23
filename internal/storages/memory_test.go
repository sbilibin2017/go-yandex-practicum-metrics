package storages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryStorage(t *testing.T) {
	// Создаем новый объект MemoryStorage
	storage := NewMemoryStorage[string, int]()

	// Проверяем, что storage не nil
	assert.NotNil(t, storage)

	// Проверяем, что начальная карта данных пуста
	assert.Empty(t, storage.data)
}
