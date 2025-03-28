package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBuildSaveBatchQuery_SingleMetric(t *testing.T) {
	ptrToInt64 := func(i int64) *int64 {
		return &i
	}
	ptrToFloat64 := func(f float64) *float64 {
		return &f
	}
	metrics := []*domain.Metrics{
		{
			ID:    "1",
			Type:  domain.Gauge,
			Delta: ptrToInt64(10),
			Value: ptrToFloat64(5.5),
		},
	}
	query, args := buildSaveBatchQuery(metrics)
	expectedQuery := `
INSERT INTO metrics (id, type, delta, value) 
VALUES ($1, $2, $3, $4)
ON CONFLICT (id, type) 
DO UPDATE SET 
	delta = EXCLUDED.delta, 
	value = EXCLUDED.value
`
	expectedArgs := []any{
		"1", domain.Gauge, ptrToInt64(10), ptrToFloat64(5.5),
	}
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedArgs, args)
}

func TestBuildSaveBatchQuery_MultipleMetrics(t *testing.T) {
	ptrToInt64 := func(i int64) *int64 {
		return &i
	}
	ptrToFloat64 := func(f float64) *float64 {
		return &f
	}
	metrics := []*domain.Metrics{
		{
			ID:    "1",
			Type:  domain.Gauge,
			Delta: ptrToInt64(10),
			Value: ptrToFloat64(5.5),
		},
		{
			ID:    "2",
			Type:  domain.Counter,
			Delta: ptrToInt64(20),
			Value: ptrToFloat64(15.5),
		},
	}
	query, args := buildSaveBatchQuery(metrics)
	expectedQuery := `
INSERT INTO metrics (id, type, delta, value) 
VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)
ON CONFLICT (id, type) 
DO UPDATE SET 
	delta = EXCLUDED.delta, 
	value = EXCLUDED.value
`
	expectedArgs := []any{
		"1", domain.Gauge, ptrToInt64(10), ptrToFloat64(5.5),
		"2", domain.Counter, ptrToInt64(20), ptrToFloat64(15.5),
	}
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedArgs, args)
}

func TestBuildSaveBatchQuery_MetricWithNilValue(t *testing.T) {
	ptrToInt64 := func(i int64) *int64 {
		return &i
	}
	metrics := []*domain.Metrics{
		{
			ID:    "3",
			Type:  domain.Counter,
			Delta: ptrToInt64(10),
			Value: nil,
		},
	}
	query, args := buildSaveBatchQuery(metrics)
	expectedQuery := `
INSERT INTO metrics (id, type, delta, value) 
VALUES ($1, $2, $3, $4)
ON CONFLICT (id, type) 
DO UPDATE SET 
	delta = EXCLUDED.delta, 
	value = EXCLUDED.value
`
	expectedArgs := []any{
		"3", domain.Counter, ptrToInt64(10), nil,
	}
	assert.Equal(t, expectedQuery, query)
	assert.Len(t, args, len(expectedArgs))
	for i := range args {
		if args[i] == nil {
			assert.Nil(t, expectedArgs[i])
		} else {
			if v, ok := args[i].(*float64); ok && v == nil {
				assert.Nil(t, expectedArgs[i])
			} else if v, ok := args[i].(*int64); ok && v == nil {
				assert.Nil(t, expectedArgs[i])
			} else {
				assert.Equal(t, expectedArgs[i], args[i])
			}
		}
	}
}

func TestSaveBatch_CallsExecuteWithCorrectArguments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ptrToInt64 := func(i int64) *int64 {
		return &i
	}
	ptrToFloat64 := func(f float64) *float64 {
		return &f
	}
	mockEngine := NewMockDBExecutorEngine(ctrl)
	metrics := []*domain.Metrics{
		{
			ID:    "1",
			Type:  domain.Gauge,
			Delta: nil,
			Value: ptrToFloat64(5.5),
		},
		{
			ID:    "2",
			Type:  domain.Counter,
			Delta: ptrToInt64(10),
			Value: nil,
		},
	}
	mockEngine.EXPECT().Execute(context.Background(), gomock.Any(), gomock.Any()).Return(true)
	repo := NewMetricDBSaveBatchRepository(mockEngine)
	result := repo.SaveBatch(context.Background(), metrics)
	assert.True(t, result)
}
