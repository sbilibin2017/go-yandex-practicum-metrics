package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"go-yandex-practicum-metrics/internal/types"
	"io"
)

// Интерфейс для работы с файлом: чтение и перемещение
type FileListReader interface {
	io.Reader
	io.Seeker
}

// Интерфейс для декодирования данных
type FileListDecoder interface {
	Decode(v interface{}) error // Для декодирования данных
}

// Репозиторий для работы с метриками в файле
type MetricFileListRepository struct {
	file    FileListReader  // Интерфейс для работы с файлом
	decoder FileListDecoder // Интерфейс для декодирования данных
}

// Конструктор для создания нового репозитория с зависимостями через DI
func NewMetricFileListRepository(file FileListReader, decoder FileListDecoder) *MetricFileListRepository {
	return &MetricFileListRepository{file: file, decoder: decoder}
}

// ListAll метод для получения всех метрик из файла
func (m *MetricFileListRepository) List(ctx context.Context) ([]*types.Metrics, error) {
	if m.file == nil {
		return nil, errors.New("file is not initialized") // Возвращаем ошибку, если файл не инициализирован
	}

	// Перемещаем указатель в начало файла
	_, err := m.file.Seek(0, 0)
	if err != nil {
		return nil, err // Возвращаем ошибку, если не удалось переместиться в начало
	}

	// Декодируем данные из файла
	decoder := json.NewDecoder(m.file)
	result := make(map[types.MetricID]*types.Metrics)

	for {
		var metric types.Metrics
		// Декодируем метрику
		if err := decoder.Decode(&metric); err != nil {
			if err.Error() == "EOF" {
				break // Конец файла, выходим
			}
			return nil, err // Возвращаем ошибку, если не удалось декодировать
		}
		result[metric.MetricID] = &metric
	}

	// Преобразуем результаты в срез
	var metricsSlice []*types.Metrics
	for _, metric := range result {
		metricsSlice = append(metricsSlice, metric)
	}

	return metricsSlice, nil // Возвращаем срез метрик и nil для ошибки
}
