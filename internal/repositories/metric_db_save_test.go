package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBuildSaveBatchQuerySingleMetric(t *testing.T) {
	metrics := []*domain.Metrics{
		{ID: "metric1", Type: domain.Counter, Delta: nil, Value: ptrFloat64(123.45)},
	}
	expectedQuery := "INSERT INTO metrics (id, type, delta, value) VALUES ($1, $2, $3, $4) ON CONFLICT (id, type) DO UPDATE SET delta = EXCLUDED.delta, value = EXCLUDED.value"
	expectedArgs := []interface{}{"metric1", domain.Counter, (*int64)(nil), ptrFloat64(123.45)}

	query, args := buildSaveBatchQuery(metrics)

	// Clean the formatting of the query by removing extra whitespace
	expectedQueryClean := strings.Join(strings.Fields(expectedQuery), " ")
	actualQueryClean := strings.Join(strings.Fields(query), " ")

	assert.Equal(t, expectedQueryClean, actualQueryClean)
	assert.Equal(t, expectedArgs, args)
}

func TestBuildSaveBatchQueryMultipleMetrics(t *testing.T) {
	metrics := []*domain.Metrics{
		{ID: "metric1", Type: domain.Counter, Delta: ptrInt64(10), Value: nil},
		{ID: "metric2", Type: domain.Gauge, Delta: nil, Value: ptrFloat64(45.67)},
	}
	expectedQuery := "INSERT INTO metrics (id, type, delta, value) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8) ON CONFLICT (id, type) DO UPDATE SET delta = EXCLUDED.delta, value = EXCLUDED.value"
	expectedArgs := []interface{}{"metric1", domain.Counter, ptrInt64(10), (*float64)(nil), "metric2", domain.Gauge, (*int64)(nil), ptrFloat64(45.67)}

	query, args := buildSaveBatchQuery(metrics)

	// Clean the formatting of the query by removing extra whitespace
	expectedQueryClean := strings.Join(strings.Fields(expectedQuery), " ")
	actualQueryClean := strings.Join(strings.Fields(query), " ")

	assert.Equal(t, expectedQueryClean, actualQueryClean)
	assert.Equal(t, expectedArgs, args)
}

func TestBuildSaveBatchQueryEmptyMetrics(t *testing.T) {
	metrics := []*domain.Metrics{}
	expectedQuery := "INSERT INTO metrics (id, type, delta, value) VALUES  ON CONFLICT (id, type) DO UPDATE SET delta = EXCLUDED.delta, value = EXCLUDED.value"
	expectedArgs := []interface{}{} // Ожидаем пустой срез

	query, args := buildSaveBatchQuery(metrics)

	// Clean the formatting of the query by removing extra whitespace
	expectedQueryClean := strings.Join(strings.Fields(expectedQuery), " ")
	actualQueryClean := strings.Join(strings.Fields(query), " ")

	// Если результат пустой срез (или nil), то проверяем, что args действительно пустой срез
	if args == nil {
		args = []interface{}{}
	}

	// Сравниваем очищенные строки запросов и аргументы
	assert.Equal(t, expectedQueryClean, actualQueryClean)
	assert.Equal(t, expectedArgs, args)
}

func TestSaveBatch_Success(t *testing.T) {
	// Подготовка мока
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockExecutor(ctrl)

	// Создаем тестовые метрики
	metrics := []*domain.Metrics{
		{ID: "metric1", Type: domain.Counter, Delta: nil, Value: float64Pointer(123.45)},
		{ID: "metric2", Type: domain.Gauge, Delta: nil, Value: float64Pointer(45.67)},
	}

	// Ожидаем, что метод Execute будет вызван один раз с нужными параметрами, но без проверки самого запроса
	mockExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()). // Проверяем, что Execute вызывается с любым запросом и аргументами
		Return(nil).Times(1)

	// Создаем репозиторий
	repo := NewMetricDBSaveBatchRepository(mockExecutor)

	// Вызов метода SaveBatch
	result := repo.SaveBatch(context.Background(), metrics)

	// Проверка, что метод вернул true (успешное выполнение)
	assert.True(t, result)
}

func TestSaveBatch_Failure(t *testing.T) {
	// Подготовка мока
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockExecutor(ctrl)

	// Создаем тестовые метрики
	metrics := []*domain.Metrics{
		{ID: "metric1", Type: domain.Counter, Delta: nil, Value: float64Pointer(123.45)},
		{ID: "metric2", Type: domain.Gauge, Delta: nil, Value: float64Pointer(45.67)},
	}

	// Ожидаем, что метод Execute будет вызван один раз с нужными параметрами, но с ошибкой в ответе
	mockExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()). // Проверяем, что Execute вызывается с любым запросом и аргументами
		Return(assert.AnError).Times(1)

	// Создаем репозиторий
	repo := NewMetricDBSaveBatchRepository(mockExecutor)

	// Вызов метода SaveBatch
	result := repo.SaveBatch(context.Background(), metrics)

	// Проверка, что метод вернул false (ошибка при выполнении)
	assert.False(t, result)
}

func TestSaveBatch_EmptyMetrics(t *testing.T) {
	// Подготовка мока
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockExecutor(ctrl)

	// Пустой срез метрик
	metrics := []*domain.Metrics{}

	// Ожидаем, что метод Execute не будет вызван, так как нет метрик
	mockExecutor.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	// Создаем репозиторий
	repo := NewMetricDBSaveBatchRepository(mockExecutor)

	// Вызов метода SaveBatch с пустым срезом
	result := repo.SaveBatch(context.Background(), metrics)

	// Проверка, что метод вернул true (поскольку нечего сохранять, предполагаем, что операция проходит успешно)
	assert.True(t, result)
}

func ptrInt64(value int64) *int64 {
	return &value
}

// Утилита для создания указателя на float64
func float64Pointer(f float64) *float64 {
	return &f
}
