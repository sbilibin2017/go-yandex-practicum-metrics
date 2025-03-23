package engines

import (
	"sync"
)

// MemoryGetterEngine отвечает за получение значений.
type MemoryGetterEngine[K comparable, V any] struct {
	*MemoryStorage[K, V]
	mu sync.RWMutex
}

// NewMemoryGetterEngine создает новый объект MemoryGetterEngine.
func NewMemoryGetterEngine[K comparable, V any](storage *MemoryStorage[K, V]) *MemoryGetterEngine[K, V] {
	return &MemoryGetterEngine[K, V]{MemoryStorage: storage}
}

// Get получает значение по ключу.
func (g *MemoryGetterEngine[K, V]) Get(key K) (V, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	value, ok := g.data[key]
	return value, ok
}
