package types

type MetricUpdateBodyRequest struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type MetricUpdateBodyResponse struct {
	MetricUpdateBodyRequest
}
