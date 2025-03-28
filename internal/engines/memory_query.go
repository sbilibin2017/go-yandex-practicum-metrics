package engines

import (
	"context"
)

type MemoryQueryEngine[K comparable, V any] struct {
	data map[K]V
}

func NewMemoryQueryEngine[K comparable, V any]() *MemoryQueryEngine[K, V] {
	return &MemoryQueryEngine[K, V]{
		data: make(map[K]V),
	}
}

func (e *MemoryQueryEngine[K, V]) Query(ctx context.Context, query string, args ...any) ([]any, bool) {
	var results []any
	for _, value := range e.data {
		results = append(results, value)
	}
	return results, true
}
