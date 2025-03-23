package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"go-yandex-practicum-metrics/internal/types"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тест успешного чтения всех метрик из файла
func TestListAll_Success(t *testing.T) {
	// Создаем временный файл
	file, err := os.CreateTemp("", "metrics_test_*.json")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	// Создаем тестовые метрики
	metric1 := &types.Metrics{MetricID: types.MetricID{ID: "metric_1", Type: "gauge"}, Value: new(float64)}
	metric2 := &types.Metrics{MetricID: types.MetricID{ID: "metric_2", Type: "counter"}, Delta: new(int64)}
	*metric1.Value = 42.0
	*metric2.Delta = 10

	// Записываем метрики в файл
	encoder := json.NewEncoder(file)
	assert.NoError(t, encoder.Encode(metric1))
	assert.NoError(t, encoder.Encode(metric2))

	// Перемещаем указатель в начало файла перед чтением
	_, err = file.Seek(0, io.SeekStart)
	assert.NoError(t, err)

	// Создаем репозиторий
	repo := NewMetricFileListAllRepository(file)

	// Вызываем ListAll
	metrics, success := repo.ListAll(context.Background())

	// Проверяем, что чтение прошло успешно
	assert.True(t, success, "Expected ListAll to return success")
	assert.Len(t, metrics, 2, "Expected to retrieve 2 metrics")

	// Проверяем содержимое
	assert.Equal(t, metric1.MetricID, metrics[0].MetricID)
	assert.Equal(t, *metric1.Value, *metrics[0].Value)
	assert.Equal(t, metric2.MetricID, metrics[1].MetricID)
	assert.Equal(t, *metric2.Delta, *metrics[1].Delta)
}

// Тест ошибки при отсутствии файла
func TestListAll_FileNil(t *testing.T) {
	repo := NewMetricFileListAllRepository(nil)

	metrics, success := repo.ListAll(context.Background())

	assert.False(t, success, "Expected ListAll to return false")
	assert.Nil(t, metrics, "Expected ListAll to return nil when file is nil")
}

// Тест ошибки при поврежденных данных
func TestListAll_CorruptData(t *testing.T) {
	// Создаем временный файл с некорректными данными
	file, err := os.CreateTemp("", "metrics_test_*.json")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	// Записываем некорректные данные
	_, err = file.WriteString("INVALID_JSON")
	assert.NoError(t, err)

	// Перемещаем указатель в начало
	_, err = file.Seek(0, io.SeekStart)
	assert.NoError(t, err)

	// Создаем репозиторий
	repo := NewMetricFileListAllRepository(file)

	// Вызываем ListAll
	metrics, success := repo.ListAll(context.Background())

	// Ожидаем ошибку
	assert.False(t, success, "Expected ListAll to return false due to corrupt data")
	assert.Nil(t, metrics, "Expected ListAll to return nil due to corrupt data")
}

// Mock для симуляции ошибки при вызове Seek()
type errorSeekReader struct{}

func (e *errorSeekReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF // Симулируем конец файла
}

func (e *errorSeekReader) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("seek error") // Всегда возвращает ошибку
}

// Тест обработки ошибки Seek()
func TestListAll_SeekError(t *testing.T) {
	// Используем кастомный reader, который всегда возвращает ошибку в Seek()
	errorReader := &errorSeekReader{}

	// Создаем репозиторий
	repo := NewMetricFileListAllRepository(errorReader)

	// Вызываем ListAll
	metrics, success := repo.ListAll(context.Background())

	// Проверяем, что вернулась ошибка
	assert.False(t, success, "Expected ListAll to return false due to Seek error")
	assert.Nil(t, metrics, "Expected ListAll to return nil due to Seek error")
}
