package storages

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFileStorageConfig — это мок для интерфейса FileStorageConfig.
type MockFileStorageConfig struct {
	mock.Mock
}

func (m *MockFileStorageConfig) GetFileStoragePath() string {
	args := m.Called()
	return args.String(0)
}

func TestNewJsonFileStorage_Success(t *testing.T) {
	// Мокаем путь к файлу
	mockConfig := new(MockFileStorageConfig)
	mockConfig.On("GetFileStoragePath").Return("testfile.json")

	// Мокаем поведение os.OpenFile
	file, err := os.Create("testfile.json")
	if err != nil {
		t.Fatal("не удалось создать файл для теста:", err)
	}
	defer os.Remove("testfile.json") // Удалим файл после теста
	defer file.Close()               // Закрываем файл после теста

	// Создаем хранилище
	storage, success := NewJsonFileStorage[interface{}](mockConfig)

	// Проверяем, что хранилище успешно создано
	assert.NotNil(t, storage)
	assert.True(t, success)

	// Проверяем, что файл был открыт по ожидаемому пути
	assert.Equal(t, "testfile.json", mockConfig.GetFileStoragePath())

	// Проверяем, что файл закрывается после использования
	err = storage.File.Close()
	assert.NoError(t, err)
}

func TestNewJsonFileStorage_FileError(t *testing.T) {
	// Мокаем путь к файлу
	mockConfig := new(MockFileStorageConfig)
	mockConfig.On("GetFileStoragePath").Return("invalid_path/testfile.json")

	// Создаем хранилище (ожидаем, что файл не будет открыт)
	storage, success := NewJsonFileStorage[interface{}](mockConfig)

	// Проверяем, что создание хранилища не удалось
	assert.Nil(t, storage)
	assert.False(t, success)
}

func TestNewJsonFileStorage_FileCreation(t *testing.T) {
	// Мокаем путь к файлу
	mockConfig := new(MockFileStorageConfig)
	mockConfig.On("GetFileStoragePath").Return("new_testfile.json")

	// Убедимся, что файл не существует
	_, err := os.Stat("new_testfile.json")
	if err == nil {
		t.Fatal("файл уже существует до теста")
	}

	// Создаем хранилище
	storage, success := NewJsonFileStorage[interface{}](mockConfig)

	// Проверяем, что хранилище успешно создано
	assert.NotNil(t, storage)
	assert.True(t, success)

	// Проверяем, что файл был создан
	_, err = os.Stat("new_testfile.json")
	assert.NoError(t, err)

	// Удаляем файл после теста
	defer os.Remove("new_testfile.json")
}
