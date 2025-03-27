package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"go-yandex-practicum-metrics/internal/domain"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBatch_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Мокируем репозитории
	mockSaveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	mockFindRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	mockTxBeginer := NewMockTxBeginer(ctrl) // Мокируем интерфейс для транзакции
	mockTx := NewMockTx(ctrl)

	// Создаем сервис с моками
	service := NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockTxBeginer)

	// Подготовка тестовых данных
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  string(domain.Gauge),
			Value: new(float64),
		},
		{
			ID:    "metric2",
			Type:  string(domain.Counter),
			Value: nil,
			Delta: new(int64),
		},
	}

	// Мокируем поведение FindBatch (возвращаем пустую карту, т.к. это новые метрики)
	mockFindRepo.EXPECT().FindBatch(gomock.Any(), gomock.Any()).Return(make(map[domain.MetricID]*domain.Metrics), nil)

	// Мокируем поведение SaveBatch (не ожидаем ошибок)
	mockSaveRepo.EXPECT().SaveBatch(gomock.Any(), gomock.Any()).Return(nil)

	// Мокируем поведение BeginTx и Commit для транзакции
	mockTxBeginer.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(mockTx, nil)
	mockTx.EXPECT().Commit().Return(nil)

	// Выполняем обновление
	updatedMetrics, err := service.UpdateBatch(context.Background(), metrics)

	// Проверяем, что ошибки нет
	assert.NoError(t, err)

	// Проверяем, что метрики обновились
	assert.Len(t, updatedMetrics, 2)
	assert.Equal(t, updatedMetrics[0].ID, "metric1")
	assert.Equal(t, updatedMetrics[1].ID, "metric2")
}

func TestUpdateBatch_FailToBeginTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Мокируем репозитории
	mockSaveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	mockFindRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	mockTxBeginer := NewMockTxBeginer(ctrl)

	// Создаем сервис с моками
	service := NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockTxBeginer)

	// Подготовка тестовых данных
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  string(domain.Gauge),
			Value: new(float64),
		},
	}

	// Мокируем поведение BeginTx, которое вернет ошибку
	mockTxBeginer.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(nil, errors.New("failed to start transaction"))

	// Выполняем обновление
	updatedMetrics, err := service.UpdateBatch(context.Background(), metrics)

	// Проверяем, что произошла ошибка
	assert.Error(t, err)
	assert.Nil(t, updatedMetrics)
	assert.Contains(t, err.Error(), "failed to start transaction")
}

func TestUpdateBatch_FailToSaveBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Мокируем репозитории
	mockSaveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	mockFindRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	mockTxBeginer := NewMockTxBeginer(ctrl)
	mockTx := NewMockTx(ctrl)

	// Создаем сервис с моками
	service := NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockTxBeginer)

	// Подготовка тестовых данных
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  string(domain.Gauge),
			Value: new(float64),
		},
	}

	// Мокируем поведение FindBatch (возвращаем пустую карту, т.к. это новые метрики)
	mockFindRepo.EXPECT().FindBatch(gomock.Any(), gomock.Any()).Return(make(map[domain.MetricID]*domain.Metrics), nil)

	// Мокируем поведение SaveBatch (возвращаем ошибку)
	mockSaveRepo.EXPECT().SaveBatch(gomock.Any(), gomock.Any()).Return(errors.New("failed to save batch"))

	// Мокируем поведение BeginTx
	mockTxBeginer.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(mockTx, nil)

	// Мокируем поведение Rollback, если произойдет ошибка
	mockTx.EXPECT().Rollback().Return(nil)

	// Выполняем обновление
	updatedMetrics, err := service.UpdateBatch(context.Background(), metrics)

	// Проверяем, что произошла ошибка
	assert.Error(t, err)
	assert.Nil(t, updatedMetrics)
	assert.Contains(t, err.Error(), "failed to save batch")
}

func TestUpdateBatch_FailToCommitTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Мокируем репозитории
	mockSaveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	mockFindRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	mockTxBeginer := NewMockTxBeginer(ctrl)
	mockTx := NewMockTx(ctrl)

	// Создаем сервис с моками
	service := NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockTxBeginer)

	// Подготовка тестовых данных
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  string(domain.Gauge),
			Value: new(float64),
		},
	}

	// Мокируем поведение FindBatch (возвращаем пустую карту, т.к. это новые метрики)
	mockFindRepo.EXPECT().FindBatch(gomock.Any(), gomock.Any()).Return(make(map[domain.MetricID]*domain.Metrics), nil)

	// Мокируем поведение SaveBatch (не ожидаем ошибок)
	mockSaveRepo.EXPECT().SaveBatch(gomock.Any(), gomock.Any()).Return(nil)

	// Мокируем поведение BeginTx
	mockTxBeginer.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(mockTx, nil)

	// Мокируем поведение Commit, которое возвращает ошибку
	mockTx.EXPECT().Commit().Return(errors.New("failed to commit transaction"))

	// Мокируем поведение Rollback, если произойдет ошибка
	mockTx.EXPECT().Rollback().Return(nil)

	// Выполняем обновление
	updatedMetrics, err := service.UpdateBatch(context.Background(), metrics)

	// Проверяем, что произошла ошибка
	assert.Error(t, err)
	assert.Nil(t, updatedMetrics)
	assert.Contains(t, err.Error(), "failed to commit transaction")
}

func TestUpdateBatch_FailToFindBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Мокируем репозитории
	mockSaveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	mockFindRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	mockTxBeginer := NewMockTxBeginer(ctrl)
	mockTx := NewMockTx(ctrl)

	// Создаем сервис с моками
	service := NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockTxBeginer)

	// Подготовка тестовых данных
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  string(domain.Gauge),
			Value: new(float64),
		},
	}

	// Мокируем поведение FindBatch (возвращаем ошибку)
	mockFindRepo.EXPECT().FindBatch(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some database error"))

	// Мокируем поведение BeginTx
	mockTxBeginer.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(mockTx, nil)

	// Мокируем поведение Rollback, которое будет вызвано в случае ошибки
	mockTx.EXPECT().Rollback().Return(nil)

	// Выполняем обновление
	updatedMetrics, err := service.UpdateBatch(context.Background(), metrics)

	// Проверяем, что произошла ошибка
	assert.Error(t, err)
	assert.Nil(t, updatedMetrics)
	assert.Contains(t, err.Error(), "failed to find batch")
}

func TestUpdateBatch_UpdateDeltaForExistingMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Мокируем репозитории
	mockSaveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	mockFindRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	mockTxBeginer := NewMockTxBeginer(ctrl)
	mockTx := NewMockTx(ctrl)

	// Создаем сервис с моками
	service := NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockTxBeginer)

	// Подготовка тестовых данных
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  string(domain.Counter), // Тип метрики - Counter
			Delta: new(int64),             // Delta, который мы хотим обновить
		},
	}

	// Мокируем поведение FindBatch (существующая метрика)
	existingMetric := &domain.Metrics{
		ID:    "metric1",
		Type:  string(domain.Counter),
		Delta: new(int64), // Существующая метрика без Delta или с нулевым значением
	}

	// Мы ожидаем, что FindBatch вернет существующую метрику
	mockFindRepo.EXPECT().FindBatch(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metrics{
		{ID: "metric1", Type: string(domain.Counter)}: existingMetric,
	}, nil)

	// Мокируем поведение BeginTx
	mockTxBeginer.EXPECT().BeginTx(gomock.Any(), gomock.Nil()).Return(mockTx, nil)

	// Мокируем поведение SaveBatch
	mockSaveRepo.EXPECT().SaveBatch(gomock.Any(), gomock.Any()).Return(nil)

	// Мокируем поведение Commit
	mockTx.EXPECT().Commit().Return(nil)

	// Выполняем обновление
	updatedMetrics, err := service.UpdateBatch(context.Background(), metrics)

	// Проверяем, что ошибки нет
	assert.NoError(t, err)

	// Проверяем, что метрика обновлена
	assert.Len(t, updatedMetrics, 1)
	assert.Equal(t, updatedMetrics[0].ID, "metric1")
	assert.NotNil(t, updatedMetrics[0].Delta)
	assert.Equal(t, *updatedMetrics[0].Delta, *existingMetric.Delta+*metrics[0].Delta)
}
