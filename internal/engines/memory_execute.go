package engines

import (
	"context"
)

type MemoryExecuteEngine[K comparable, V any] struct {
	data map[K]V
}

func NewMemoryExecuteEngine[K comparable, V any]() *MemoryExecuteEngine[K, V] {
	return &MemoryExecuteEngine[K, V]{
		data: make(map[K]V),
	}
}

func (e *MemoryExecuteEngine[K, V]) Execute(ctx context.Context, query string, args ...any) bool {
	if len(args) < 2 {
		return false
	}
	key, ok := args[0].(K)
	if !ok {
		return false
	}
	value, ok := args[1].(V)
	if !ok {
		return false
	}
	e.data[key] = value
	return true
}
