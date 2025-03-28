package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBuildFindQuery(t *testing.T) {
	tests := []struct {
		name          string
		filters       []domain.MetricID
		expectedQuery string
		expectedArgs  []any
	}{
		{
			name: "single filter with Gauge type",
			filters: []domain.MetricID{
				{ID: "1", Type: domain.Gauge},
			},
			expectedQuery: "SELECT id, type, delta, value FROM metrics WHERE (id = $1 AND type = $2)",
			expectedArgs:  []any{"1", domain.Gauge},
		},
		{
			name: "single filter with Counter type",
			filters: []domain.MetricID{
				{ID: "2", Type: domain.Counter},
			},
			expectedQuery: "SELECT id, type, delta, value FROM metrics WHERE (id = $1 AND type = $2)",
			expectedArgs:  []any{"2", domain.Counter},
		},
		{
			name: "multiple filters with mixed types",
			filters: []domain.MetricID{
				{ID: "1", Type: domain.Gauge},
				{ID: "2", Type: domain.Counter},
			},
			expectedQuery: "SELECT id, type, delta, value FROM metrics WHERE (id = $1 AND type = $2) OR (id = $3 AND type = $4)",
			expectedArgs:  []any{"1", domain.Gauge, "2", domain.Counter},
		},
		{
			name:          "empty filters",
			filters:       []domain.MetricID{},
			expectedQuery: "SELECT id, type, delta, value FROM metrics WHERE ",
			expectedArgs:  []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, args := buildFindQuery(tt.filters)
			assert.Equal(t, tt.expectedQuery, query)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestFindBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ptrToInt64 := func(i int64) *int64 {
		return &i
	}
	ptrToFloat64 := func(f float64) *float64 {
		return &f
	}
	mockDBQuerier := NewMockDBQuerierEngine(ctrl)
	mockDBScanner := NewMockDBScannerEngine(ctrl)
	repo := NewMetricDBFindBatchRepository(mockDBQuerier, mockDBScanner)
	tests := []struct {
		name          string
		filters       []domain.MetricID
		mockQueryResp DBScannerEngine
		mockScanResp  map[any]any
		mockQueryOk   bool
		mockScanOk    bool
		expectedRes   map[domain.MetricID]*domain.Metrics
		expectedOk    bool
	}{
		{
			name: "successful query and scan",
			filters: []domain.MetricID{
				{ID: "1", Type: domain.Gauge},
			},
			mockQueryResp: mockDBScanner,
			mockScanResp: map[any]any{
				"row1": map[string]any{
					"id":    "1",
					"type":  "gauge",
					"delta": int64(10),
					"value": float64(5.5),
				},
			},
			mockQueryOk: true,
			mockScanOk:  true,
			expectedRes: map[domain.MetricID]*domain.Metrics{
				{"1", domain.Gauge}: {
					ID:    "1",
					Type:  domain.Gauge,
					Delta: ptrToInt64(10),
					Value: ptrToFloat64(5.5),
				},
			},
			expectedOk: true,
		},
		{
			name:          "failed query",
			filters:       []domain.MetricID{{ID: "2", Type: domain.Counter}},
			mockQueryResp: nil,
			mockQueryOk:   false,
			expectedRes:   nil,
			expectedOk:    false,
		},
		{
			name: "successful query but failed scan",
			filters: []domain.MetricID{
				{ID: "3", Type: domain.Counter},
			},
			mockQueryResp: mockDBScanner,
			mockScanResp:  nil,
			mockQueryOk:   true,
			mockScanOk:    false,
			expectedRes:   nil,
			expectedOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDBQuerier.EXPECT().
				Query(context.Background(), gomock.Any(), gomock.Any()).
				Return(tt.mockQueryResp, tt.mockQueryOk)
			if tt.mockQueryOk {
				mockDBScanner.EXPECT().
					Scan(context.Background(), gomock.Any(), gomock.Any()).
					Return(tt.mockScanResp, tt.mockScanOk)
			}
			result, ok := repo.FindBatch(context.Background(), tt.filters)
			assert.Equal(t, tt.expectedRes, result)
			assert.Equal(t, tt.expectedOk, ok)
		})
	}
}
