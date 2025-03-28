package engines

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryQueryEngine_Query(t *testing.T) {
	engine := NewMemoryQueryEngine[string, int]()
	engine.data["key1"] = 42
	engine.data["key2"] = 99
	results, success := engine.Query(context.Background(), "SELECT")
	assert.True(t, success)
	assert.Len(t, results, 2)
	assert.Contains(t, results, 42)
	assert.Contains(t, results, 99)
}

func TestMemoryQueryEngine_Query_EmptyData(t *testing.T) {
	engine := NewMemoryQueryEngine[string, int]()
	results, success := engine.Query(context.Background(), "SELECT")
	assert.True(t, success)
	assert.Len(t, results, 0)
}

func TestMemoryQueryEngine_Query_InvalidData(t *testing.T) {
	engine := NewMemoryQueryEngine[string, string]()
	engine.data["key1"] = "value1"
	engine.data["key2"] = "value2"
	results, success := engine.Query(context.Background(), "SELECT")
	assert.True(t, success)
	assert.Len(t, results, 2)
	assert.Contains(t, results, "value1")
	assert.Contains(t, results, "value2")
}
