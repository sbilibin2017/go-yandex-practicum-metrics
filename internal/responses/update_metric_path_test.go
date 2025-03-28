package responses

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricUpdatePathResponse_ToResponse(t *testing.T) {
	response := &MetricUpdatePathResponse{}
	expected := []byte("Metric updated successfully")
	result := response.ToResponse()

	assert.NotNil(t, result)
	assert.Equal(t, expected, *result)
}
