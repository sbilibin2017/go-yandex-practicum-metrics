package usecases

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/logger"
)

type MetricUpdatePathService interface {
	Update(ctx context.Context, metrics []*domain.Metrics) ([]*domain.Metrics, error)
}

type MetricUpdatePathUsecase struct {
	svc MetricUpdatePathService
}

func NewMetricUpdatePathUsecase(
	svc MetricUpdatePathService,
) *MetricUpdatePathUsecase {
	return &MetricUpdatePathUsecase{
		svc: svc,
	}
}

type MetricUpdatePathRequest struct {
	Type  string
	Name  string
	Value string
}

type MetricUpdatePathResponse struct {
	Message string
}

func (uc *MetricUpdatePathUsecase) Execute(
	ctx context.Context,
	req *MetricUpdatePathRequest,
) (*MetricUpdatePathResponse, error) {
	logger.Info("Executing MetricUpdatePathUsecase", "type", req.Type, "name", req.Name)
	metrics, err := domain.NewMetrics(req.Type, req.Name, req.Value)
	if err != nil {
		logger.Error("Error creating new metrics", "error", err)
		return nil, err
	}
	_, err = uc.svc.Update(ctx, []*domain.Metrics{metrics})
	if err != nil {
		logger.Error("Error updating metrics", "error", err)
		return nil, err
	}
	logger.Info("Metric updated successfully", "name", req.Name)
	return &MetricUpdatePathResponse{
		Message: "Metric updated successfully",
	}, nil
}
