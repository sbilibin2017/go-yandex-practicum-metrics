package requests

import (
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdatePathRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		request   *MetricUpdatePathRequest
		expectErr bool
	}{
		{
			name: "valid Gauge request",
			request: &MetricUpdatePathRequest{
				Type:  string(domain.Gauge),
				Name:  "metric1",
				Value: "12.34",
			},
			expectErr: false,
		},
		{
			name: "valid Counter request",
			request: &MetricUpdatePathRequest{
				Type:  string(domain.Counter),
				Name:  "metric2",
				Value: "100",
			},
			expectErr: false,
		},
		{
			name: "invalid empty Name",
			request: &MetricUpdatePathRequest{
				Type:  string(domain.Gauge),
				Name:  "",
				Value: "12.34",
			},
			expectErr: true,
		},
		{
			name: "invalid invalid Metric Type",
			request: &MetricUpdatePathRequest{
				Type:  "InvalidType",
				Name:  "metric3",
				Value: "12.34",
			},
			expectErr: true,
		},
		{
			name: "invalid invalid Gauge Value",
			request: &MetricUpdatePathRequest{
				Type:  string(domain.Gauge),
				Name:  "metric4",
				Value: "invalid_value",
			},
			expectErr: true,
		},
		{
			name: "invalid invalid Counter Value",
			request: &MetricUpdatePathRequest{
				Type:  string(domain.Counter),
				Name:  "metric5",
				Value: "invalid_value",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMetricUpdatePathRequest_ToDomain(t *testing.T) {
	tests := []struct {
		name    string
		request *MetricUpdatePathRequest
		want    *domain.Metrics
	}{
		{
			name: "valid Gauge request",
			request: &MetricUpdatePathRequest{
				Type:  string(domain.Gauge),
				Name:  "metric1",
				Value: "12.34",
			},
			want: &domain.Metrics{
				ID:    "metric1",
				Type:  domain.Gauge,
				Value: floatPointer(12.34),
			},
		},
		{
			name: "valid Counter request",
			request: &MetricUpdatePathRequest{
				Type:  string(domain.Counter),
				Name:  "metric2",
				Value: "100",
			},
			want: &domain.Metrics{
				ID:    "metric2",
				Type:  domain.Counter,
				Delta: int64Pointer(100),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.ToDomain()
			require.Equal(t, tt.want, got)
		})
	}
}

func floatPointer(v float64) *float64 {
	return &v
}

func int64Pointer(v int64) *int64 {
	return &v
}
