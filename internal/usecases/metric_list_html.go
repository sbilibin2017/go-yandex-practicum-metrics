package usecases

import (
	"context"
	"fmt"
	"go-yandex-practicum-metrics/internal/domain"
)

type MetricListService interface {
	List(ctx context.Context) ([]*domain.Metrics, error)
}

type MetricListHTMLUsecase struct {
	svc MetricListService
}

func NewMetricListHTMLUsecase(svc MetricListService) *MetricListHTMLUsecase {
	return &MetricListHTMLUsecase{svc: svc}
}

func (uc *MetricListHTMLUsecase) Execute(ctx context.Context) (*MetricListHTMLResponse, error) {
	metrics, err := uc.svc.List(ctx)
	if err != nil {
		return nil, err
	}
	resp := &MetricListHTMLResponse{}
	err = resp.FromDomain(metrics)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type MetricListHTMLResponse struct {
	HTML string `json:"html"`
}

func (resp *MetricListHTMLResponse) FromDomain(metrics []*domain.Metrics) error {
	if metrics == nil {
		resp.HTML = "<html><body><h1>Metric List</h1><table border='1'><thead><tr><th>ID</th><th>Value</th></tr></thead><tbody></tbody></table></body></html>"
		return nil
	}
	html := "<html><body><h1>Metric List</h1><table border='1'><thead><tr><th>ID</th><th>Value</th></tr></thead><tbody>"
	for _, metric := range metrics {
		var value string
		switch metric.Type {
		case string(domain.Counter):
			if metric.Delta != nil {
				value = fmt.Sprintf("%d", *metric.Delta)
			} else {
				value = "N/A"
			}
		case string(domain.Gauge):
			if metric.Value != nil {
				value = fmt.Sprintf("%.2f", *metric.Value)
			} else {
				value = "N/A"
			}
		default:
			value = "Unknown type"
		}
		html += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", metric.ID, value)
	}
	html += "</tbody></table></body></html>"
	resp.HTML = html
	return nil
}
