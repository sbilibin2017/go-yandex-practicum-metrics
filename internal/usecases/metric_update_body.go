package usecases

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
)

type MetricUpdateBodyRequest struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type MetricUpdateBodyResponse struct {
	MetricUpdateBodyRequest
}

type MetricUpdateBodyService interface {
	Update(ctx context.Context, metric *domain.Metrics) (*domain.Metrics, error)
}

type MetricUpdateBodyUsecase struct {
	svc MetricUpdateBodyService
}

var (
	ErrMetricUpdateBodyInternal     = errors.New("internal error")
	ErrMetricUpdateBodyTypeRequired = errors.New("metric type required")
	ErrMetricUpdateBodyInvalidType  = errors.New("metric counter or gauge required")
	ErrMetricUpdateBodyIDRequired   = errors.New("metric id required")
	ErrMetricUpdateBodyInvalidValue = errors.New("metric value is invalid")
)

func NewMetricUpdateBodyUsecase(svc MetricUpdateBodyService) *MetricUpdateBodyUsecase {
	return &MetricUpdateBodyUsecase{svc: svc}
}

func (uc *MetricUpdateBodyUsecase) Execute(
	ctx context.Context, req *MetricUpdateBodyRequest,
) (*MetricUpdateBodyResponse, error) {
	if req.ID == "" {
		return nil, ErrMetricUpdateBodyIDRequired
	}
	if req.Type == "" {
		return nil, ErrMetricUpdateBodyTypeRequired
	}
	if req.Type != string(domain.Counter) && req.Type != string(domain.Gauge) {
		return nil, ErrMetricUpdateBodyInvalidType
	}

	var metric domain.Metrics
	metric.ID = req.ID
	metric.Type = domain.MetricType(req.Type)

	switch metric.Type {
	case domain.Counter:
		if req.Delta == nil {
			return nil, ErrMetricUpdateBodyInvalidValue
		}
		metric.Delta = req.Delta
	case domain.Gauge:
		if req.Value == nil {
			return nil, ErrMetricUpdateBodyInvalidValue
		}
		metric.Value = req.Value
	}

	updated, err := uc.svc.Update(ctx, &metric)
	if err != nil {
		return nil, ErrMetricUpdateBodyInternal
	}

	response := MetricUpdateBodyResponse{
		MetricUpdateBodyRequest{
			ID:    updated.ID,
			Type:  string(updated.Type),
			Delta: updated.Delta,
			Value: updated.Value,
		},
	}
	return &response, nil
}
