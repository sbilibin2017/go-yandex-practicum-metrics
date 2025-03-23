package engines

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Структура для теста, которая будет использоваться для декодирования данных
type TestRow struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestGenerateRow(t *testing.T) {
	// Подготовим данные для теста
	testData := []TestRow{
		{"Test1", 1},
		{"Test2", 2},
		{"Test3", 3},
	}

	// Создадим временный файл для тестирования
	tempFile, err := os.CreateTemp("", "testfile_*.json")
	require.NoError(t, err, "Failed to create temp file")
	defer os.Remove(tempFile.Name()) // Удаляем файл после завершения теста

	// Кодируем тестовые данные в JSON и записываем в файл
	enc := json.NewEncoder(tempFile)
	for _, row := range testData {
		require.NoError(t, enc.Encode(row), "Failed to encode row to file")
	}

	// Подготовим FileStorageConfig для создания JsonFileStorage
	config := &MockFilesStorageConfig{
		filePath: tempFile.Name(),
	}

	// Создаем экземпляр JsonFileStorage
	jsonFileStorage, ok := NewJsonFileStorage[TestRow](config)
	require.True(t, ok, "Failed to create JsonFileStorage")

	// Создаем FileGenerator для теста
	fg := &FileGenerator[TestRow]{JsonFileStorage: jsonFileStorage}

	// Генерация канала с данными
	ch, ok := fg.GenerateRow()
	assert.True(t, ok, "Expected GenerateRow to return true")

	// Проверяем данные в канале
	for i, expectedRow := range testData {
		// Получаем следующее значение из канала
		actualRow := <-ch
		// Сравниваем значение из канала с ожидаемым
		assert.Equal(t, expectedRow, actualRow, "Row %d did not match expected data", i)
	}

	// Проверяем, что канал закрывается после завершения чтения
	_, ok = <-ch
	assert.False(t, ok, "Expected channel to be closed after reading all rows")
}

// Мок для интерфейса FileStorageConfig
type MockFilesStorageConfig struct {
	filePath string
}

func (m *MockFilesStorageConfig) GetFileStoragePath() string {
	return m.filePath
}
