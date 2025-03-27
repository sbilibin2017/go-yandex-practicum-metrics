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
	metric, err := MetricUpdatePathRequestToDomain(req)
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

func MetricUpdatePathRequestToDomain(req *MetricUpdatePathRequest) (*domain.Metrics, error) {
	var metricType string
	switch req.Type {
	case string(domain.Gauge):
		metricType = string(domain.Gauge)
	case string(domain.Counter):
		metricType = string(domain.Counter)
	default:
		return nil, ErrInvalidPathMetricType
	}

	var delta *int64
	var value *float64
	switch metricType {
	case string(domain.Gauge):
		v, err := strconv.ParseFloat(req.Value, 64)
		if err != nil {
			return nil, ErrInvalidPathMetricValue
		}
		value = &v
	case string(domain.Counter):
		v, err := strconv.ParseInt(req.Value, 10, 64)
		if err != nil {
			return nil, ErrInvalidPathMetricValue
		}
		delta = &v
	}

	return &domain.Metrics{
		ID:    req.Name,
		Type:  metricType,
		Delta: delta,
		Value: value,
	}, nil
}

type MetricUpdatePathResponse []byte

var (
	ErrInvalidPathMetricType  = errors.New("invalid metric type")
	ErrInvalidPathMetricValue = errors.New("invalid metric value")
	ErrMissingPathMetricName  = errors.New("missing metric name")
)

var MetricUpdateSuccessMessage = []byte("Metric updated successfully")
