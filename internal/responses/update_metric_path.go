package responses

type MetricUpdatePathResponse struct {
	Message []byte
}

func (r *MetricUpdatePathResponse) ToResponse() *[]byte {
	var successMessage = []byte("Metric updated successfully")
	return &successMessage
}
