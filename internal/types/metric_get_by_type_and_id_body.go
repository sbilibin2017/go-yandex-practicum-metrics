package types

type MetricGetByTypeAndIDBodyRequest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type MetricGetByTypeAndIDBodyResponse struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
