package usecases

import (
	"context"
	"fmt"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
	"strings"
)

type MetricListHTMLService interface {
	List(ctx context.Context) ([]*types.Metrics, error)
}

type MetricListHTMLUsecase struct {
	svc MetricListHTMLService
}

func NewMetricListHTMLUsecase(
	svc MetricListHTMLService,
) *MetricListHTMLUsecase {
	return &MetricListHTMLUsecase{svc: svc}
}

func (uc *MetricListHTMLUsecase) Execute(
	ctx context.Context,
) (string, error) {
	metrics, err := uc.svc.List(ctx)
	if err != nil {
		return "", errors.ErrMetricInternal
	}

	if len(metrics) == 0 {
		return "<h1>No metrics available</h1>", nil
	}

	var sb strings.Builder
	sb.WriteString("<html><body><h1>Metrics List</h1><table border='1'>")
	sb.WriteString("<tr><th>ID</th><th>Value</th></tr>")
	for _, metric := range metrics {
		var value string
		switch metric.Type {
		case types.Counter:
			if metric.Delta != nil {
				value = fmt.Sprintf("%d", *metric.Delta)
			}
		case types.Gauge:
			if metric.Value != nil {
				value = fmt.Sprintf("%.2f", *metric.Value)
			}
		}

		sb.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", metric.ID, value))
	}
	sb.WriteString("</table></body></html>")

	return sb.String(), nil
}
