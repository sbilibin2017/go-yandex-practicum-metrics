package engines

import (
	"sync"
)

// MemoryRangerEngine отвечает за перебор всех элементов.
type MemoryRangerEngine[K comparable, V any] struct {
	*MemoryStorage[K, V]
	mu sync.RWMutex
}

// NewMemoryRangerEngine создает новый объект MemoryRangerEngine.
func NewMemoryRangerEngine[K comparable, V any](storage *MemoryStorage[K, V]) *MemoryRangerEngine[K, V] {
	return &MemoryRangerEngine[K, V]{MemoryStorage: storage}
}

// Range перебирает все пары ключ-значение и применяет функцию f к каждому элементу.
func (r *MemoryRangerEngine[K, V]) Range(f func(key K, value V) bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for key, value := range r.data {
		if !f(key, value) {
			break
		}
	}
}
