package repositories

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSeekerReader struct {
	data   []byte
	offset int64
}

func (m *mockSeekerReader) Read(p []byte) (n int, err error) {
	n = copy(p, m.data[m.offset:])
	if n == len(m.data)-int(m.offset) {
		err = io.EOF
	}
	m.offset += int64(n)
	return n, err
}

func (m *mockSeekerReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		m.offset = offset
	case io.SeekCurrent:
		m.offset += offset
	case io.SeekEnd:
		m.offset = int64(len(m.data)) + offset
	}
	return m.offset, nil
}

func TestFileFindMetricsSuccessfully(t *testing.T) {
	data := `
		{"id": "metric1", "type": "counter", "value": 100}
		{"id": "metric2", "type": "counter", "value": 200}
	`
	// Use custom mock reader and seeker
	mock := &mockSeekerReader{data: []byte(data)}
	repo := NewMetricFileFindBatchRepository(mock, mock)

	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
	}

	result, success := repo.Find(context.Background(), filters)

	assert.True(t, success)
	assert.Len(t, result, 1)
	assert.Equal(t, "metric1", result[domain.MetricID{ID: "metric1", Type: "counter"}].ID)
}

func TestFindWithNilReader(t *testing.T) {
	repo := NewMetricFileFindBatchRepository(nil, nil)

	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
	}

	result, success := repo.Find(context.Background(), filters)

	assert.False(t, success)
	assert.Nil(t, result)
}

type mockSeekerError struct{}

func (m *mockSeekerError) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("seek error")
}

func (m *mockSeekerError) Read(p []byte) (n int, err error) {
	return 0, nil
}

func TestFindWithSeekError(t *testing.T) {
	mockReader := &mockSeekerError{}
	repo := NewMetricFileFindBatchRepository(mockReader, mockReader)

	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
	}

	result, success := repo.Find(context.Background(), filters)

	assert.False(t, success)
	assert.Nil(t, result)
}

func TestFindWithEOF(t *testing.T) {
	data := `
		{"id": "metric1", "type": "counter", "value": 100}
		{"id": "metric2", "type": "counter", "value": 200}
	`
	mock := &mockSeekerReader{data: []byte(data)}
	repo := NewMetricFileFindBatchRepository(mock, mock)

	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
	}

	result, success := repo.Find(context.Background(), filters)

	assert.True(t, success)
	assert.Len(t, result, 1)
	assert.Equal(t, "metric1", result[domain.MetricID{ID: "metric1", Type: "counter"}].ID)
}
