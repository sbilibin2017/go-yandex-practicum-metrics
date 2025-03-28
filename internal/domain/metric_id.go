package domain

type MetricID struct {
	ID   string     `json:"id"`
	Type MetricType `json:"type"`
}
