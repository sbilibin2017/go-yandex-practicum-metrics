package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFind_SingleMetricFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEngine := NewMockMemoryQuerierEngine(ctrl)
	mockEngine.EXPECT().Query(context.Background(), "").Return([]any{
		&domain.Metrics{ID: "1", Type: domain.Gauge},
	}, true)

	repo := NewMetricMemoryFindBatchRepository(mockEngine)
	filters := []domain.MetricID{{ID: "1", Type: domain.Gauge}}

	result, ok := repo.Find(context.Background(), filters)

	assert.True(t, ok)
	assert.Len(t, result, 1)
	assert.Equal(t, "1", result[domain.MetricID{ID: "1", Type: domain.Gauge}].ID)
}

func TestFind_MultipleMetricsFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEngine := NewMockMemoryQuerierEngine(ctrl)
	mockEngine.EXPECT().Query(context.Background(), "").Return([]any{
		&domain.Metrics{ID: "1", Type: domain.Gauge},
		&domain.Metrics{ID: "2", Type: domain.Counter},
	}, true)

	repo := NewMetricMemoryFindBatchRepository(mockEngine)
	filters := []domain.MetricID{
		{ID: "1", Type: domain.Gauge},
		{ID: "2", Type: domain.Counter},
	}

	result, ok := repo.Find(context.Background(), filters)

	assert.True(t, ok)
	assert.Len(t, result, 2)
	assert.Equal(t, "1", result[domain.MetricID{ID: "1", Type: domain.Gauge}].ID)
	assert.Equal(t, "2", result[domain.MetricID{ID: "2", Type: domain.Counter}].ID)
}

func TestFind_QueryFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEngine := NewMockMemoryQuerierEngine(ctrl)
	mockEngine.EXPECT().Query(context.Background(), "").Return(nil, false)

	repo := NewMetricMemoryFindBatchRepository(mockEngine)
	filters := []domain.MetricID{{ID: "1", Type: domain.Gauge}}

	result, ok := repo.Find(context.Background(), filters)

	assert.False(t, ok)
	assert.Nil(t, result)
}
