package usecases

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
	"strconv"
)

type MetricUpdatePathService interface {
	UpdateBatch(ctx context.Context, metrics []*domain.Metrics) ([]*domain.Metrics, error)
}

type MetricUpdatePathUsecase struct {
	svc MetricUpdatePathService
}

func NewMetricUpdatePathUsecase(svc MetricUpdatePathService) *MetricUpdatePathUsecase {
	return &MetricUpdatePathUsecase{svc: svc}
}

func (uc MetricUpdatePathUsecase) Execute(
	ctx context.Context, req *MetricUpdatePathRequest,
) (*MetricUpdatePathResponse, error) {
	metric, err := req.ToDomain()
	if err != nil {
		return nil, err
	}
	metricsToSave := []*domain.Metrics{metric}
	_, err = uc.svc.UpdateBatch(ctx, metricsToSave)
	if err != nil {
		return nil, err
	}
	successMessage := MetricUpdatePathResponse(MetricUpdateSuccessMessage)
	return &successMessage, nil
}

type MetricUpdatePathRequest struct {
	Type  string
	Name  string
	Value string
}

func (r *MetricUpdatePathRequest) ToDomain() (*domain.Metrics, error) {
	var metricType string
	switch r.Type {
	case string(domain.Gauge):
		metricType = string(domain.Gauge)
	case string(domain.Counter):
		metricType = string(domain.Counter)
	default:
		return nil, errors.New("invalid metric type")
	}
	if r.Name == "" {
		return nil, errors.New("missing metric name")
	}
	var delta *int64
	var value *float64
	switch metricType {
	case string(domain.Gauge):
		v, err := strconv.ParseFloat(r.Value, 64)
		if err != nil {
			return nil, errors.New("invalid metric value")
		}
		value = &v
	case string(domain.Counter):
		v, err := strconv.ParseInt(r.Value, 10, 64)
		if err != nil {
			return nil, errors.New("invalid metric value")
		}
		delta = &v
	}

	return &domain.Metrics{
		ID:    r.Name,
		Type:  metricType,
		Delta: delta,
		Value: value,
	}, nil
}

type MetricUpdatePathResponse []byte

var MetricUpdateSuccessMessage = []byte("Metric updated successfully")
