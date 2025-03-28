package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFind_SuccessfulFind(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEngine := NewMockFileQuerierEngine(ctrl)
	mockEngine.EXPECT().Query(context.Background(), "").Return([]any{
		&domain.Metrics{
			ID:    "1",
			Type:  domain.Gauge,
			Delta: ptrToInt64(10),
			Value: ptrToFloat64(5.5),
		},
		&domain.Metrics{
			ID:    "2",
			Type:  domain.Counter,
			Delta: ptrToInt64(20),
			Value: ptrToFloat64(15.5),
		},
	}, true)
	repo := NewMetricFileFindBatchRepository(mockEngine)
	filters := []domain.MetricID{
		{ID: "1", Type: domain.Gauge},
	}
	result, ok := repo.Find(context.Background(), filters)
	assert.True(t, ok)
	assert.Len(t, result, 1)
	assert.Equal(t, result[domain.MetricID{ID: "1", Type: domain.Gauge}].ID, "1")
}

func TestFind_QueryFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEngine := NewMockFileQuerierEngine(ctrl)

	// Simulate the query failure
	mockEngine.EXPECT().Query(context.Background(), "").Return(nil, false)

	// Repository initialization
	repo := NewMetricFileFindBatchRepository(mockEngine)

	filters := []domain.MetricID{
		{ID: "1", Type: domain.Gauge},
	}

	// Call Find method
	result, ok := repo.Find(context.Background(), filters)

	// Assertions
	assert.False(t, ok)
	assert.Nil(t, result)
}

// Helper functions to create pointers to basic types
func ptrToInt64(i int64) *int64 {
	return &i
}

func ptrToFloat64(f float64) *float64 {
	return &f
}
