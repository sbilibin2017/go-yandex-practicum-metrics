package usecases

import (
	"context"
	"fmt"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
)

type MetricGetByIDPathService interface {
	GetByID(ctx context.Context, id *domain.MetricID) (*domain.Metrics, error)
}

type MetricGetByIDPathUsecase struct {
	svc MetricGetByIDPathService
}

func NewMetricGetByIDPathUsecase(svc MetricGetByIDPathService) *MetricGetByIDPathUsecase {
	return &MetricGetByIDPathUsecase{svc: svc}
}

func (uc MetricGetByIDPathUsecase) Execute(
	ctx context.Context, req *MetricGetByIDPathRequest,
) (*MetricGetByIDPathResponse, error) {
	id, err := req.ToDomain()
	if err != nil {
		return nil, err
	}
	metric, err := uc.svc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := MetricGetByIDPathResponse{}
	err = resp.FromDomain(metric)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type MetricGetByIDPathRequest struct {
	Type string
	Name string
}

func (req *MetricGetByIDPathRequest) ToDomain() (*domain.MetricID, error) {
	if req.Name == "" {
		return nil, errors.MissingMetricNameError
	}
	switch req.Type {
	case string(domain.Gauge), string(domain.Counter):
		return &domain.MetricID{
			ID:   req.Name,
			Type: req.Type,
		}, nil
	default:
		return nil, errors.InvalidMetricTypeError
	}
}

type MetricGetByIDPathResponse []byte

func (resp *MetricGetByIDPathResponse) FromDomain(metric *domain.Metrics) error {
	var value string
	switch metric.Type {
	case string(domain.Counter):
		if metric.Delta == nil {
			return errors.InvalidMetricValueError
		}
		value = fmt.Sprintf("%d", *metric.Delta)
	case string(domain.Gauge):
		if metric.Value == nil {
			return errors.InvalidMetricValueError
		}
		value = fmt.Sprintf("%f", *metric.Value)
	default:
		return errors.InvalidMetricTypeError
	}
	*resp = MetricGetByIDPathResponse(value)
	return nil
}
