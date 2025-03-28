package usecases

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
)

type MetricUpdatePathService interface {
	Update(ctx context.Context, metrics []*domain.Metrics) ([]*domain.Metrics, error)
}

type MetricUpdateRequester interface {
	Validate() error
	ToDomain() []*domain.Metrics
}

type MetricUpdateResponser interface {
	ToResponse() *[]byte
}

type MetricUpdatePathUsecase struct {
	svc       MetricUpdatePathService
	requester MetricUpdateRequester
	responser MetricUpdateResponser
}

func NewMetricUpdatePathUsecase(
	svc MetricUpdatePathService,
	requester MetricUpdateRequester,
	responser MetricUpdateResponser,
) *MetricUpdatePathUsecase {
	return &MetricUpdatePathUsecase{
		svc:       svc,
		requester: requester,
		responser: responser,
	}
}

func (uc MetricUpdatePathUsecase) Execute(ctx context.Context) (*[]byte, error) {
	err := uc.requester.Validate()
	if err != nil {
		return nil, err
	}
	metrics := uc.requester.ToDomain()
	_, err = uc.svc.Update(ctx, metrics)
	if err != nil {
		return nil, err
	}
	return uc.responser.ToResponse(), nil
}
