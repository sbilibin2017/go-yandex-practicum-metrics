package usecases

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
)

type MetricGetByIDBodyService interface {
	GetByID(ctx context.Context, id *domain.MetricID) (*domain.Metrics, error)
}

type MetricGetByIDBodyUsecase struct {
	svc MetricGetByIDBodyService
}

func NewMetricGetByIDBodyUsecase(svc MetricGetByIDBodyService) *MetricGetByIDBodyUsecase {
	return &MetricGetByIDBodyUsecase{svc: svc}
}

func (uc MetricGetByIDBodyUsecase) Execute(
	ctx context.Context, req *MetricGetByIDBodyRequest,
) (*MetricGetByIDBodyResponse, error) {
	id, err := req.ToDomain()
	if err != nil {
		return nil, err
	}
	metric, err := uc.svc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := MetricGetByIDBodyResponse{}
	err = resp.FromDomain(metric)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type MetricGetByIDBodyRequest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (req *MetricGetByIDBodyRequest) ToDomain() (*domain.MetricID, error) {
	if req.ID == "" {
		return nil, errors.New("missing metric id")
	}
	if req.Type != string(domain.Gauge) && req.Type != string(domain.Counter) {
		return nil, errors.New("invalid metric type")
	}
	return &domain.MetricID{
		ID:   req.ID,
		Type: req.Type,
	}, nil
}

type MetricGetByIDBodyResponse struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (resp *MetricGetByIDBodyResponse) FromDomain(metric *domain.Metrics) error {
	resp.ID = metric.ID
	resp.Type = metric.Type
	resp.Delta = metric.Delta
	resp.Value = metric.Value
	return nil
}
