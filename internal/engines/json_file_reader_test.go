package engines

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonFileReaderEngine_ReadRow(t *testing.T) {
	// Создаем временный файл
	tempFile, err := os.CreateTemp("", "test_storage.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Данные, которые будем записывать в файл
	data := map[string]string{"key1": "value1"}

	// Записываем данные в файл
	encoder := json.NewEncoder(tempFile)
	encoder.SetIndent("", "  ")
	require.NoError(t, encoder.Encode(data))

	// Закрываем файл перед чтением
	require.NoError(t, tempFile.Close())

	// Открываем файл для чтения
	tempFile, err = os.Open(tempFile.Name())
	require.NoError(t, err)
	defer tempFile.Close()

	// Инициализируем хранилище и движок для чтения
	storage := &JsonFileStorage[map[string]string]{File: tempFile}
	sut := &JsonFileReaderEngine[map[string]string]{storage: storage}

	// Читаем данные с помощью ReadRow
	result, ok := sut.ReadRow()
	assert.True(t, ok)
	assert.Equal(t, data, *result)

	// Чтение следующей строки (должен быть EOF)
	result, ok = sut.ReadRow()
	assert.False(t, ok)
	assert.Nil(t, result)
}
