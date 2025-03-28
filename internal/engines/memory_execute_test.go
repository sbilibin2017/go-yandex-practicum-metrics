package engines

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryExecuteEngine_Execute(t *testing.T) {
	engine := NewMemoryExecuteEngine[string, int]()
	success := engine.Execute(context.Background(), "INSERT", "key1", 42)
	assert.True(t, success)
	value, exists := engine.data["key1"]
	assert.True(t, exists)
	assert.Equal(t, 42, value)
}

func TestMemoryExecuteEngine_Execute_InvalidArgs(t *testing.T) {
	engine := NewMemoryExecuteEngine[string, int]()
	success := engine.Execute(context.Background(), "INSERT", "key1")
	assert.False(t, success)
	success = engine.Execute(context.Background(), "INSERT", 42, "value")
	assert.False(t, success)
}

func TestMemoryExecuteEngine_Execute_EmptyArgs(t *testing.T) {
	engine := NewMemoryExecuteEngine[string, int]()
	success := engine.Execute(context.Background(), "INSERT")
	assert.False(t, success)
}
