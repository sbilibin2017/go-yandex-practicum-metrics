package domain

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)
