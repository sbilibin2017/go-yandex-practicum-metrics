package types

type MetricID struct {
	ID   string     `json:"id"`
	Type MetricType `json:"type"`
}
