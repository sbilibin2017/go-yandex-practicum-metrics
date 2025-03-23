package engines

import (
	"sync"
)

// Setter отвечает за установку значений.
type MemorySetterEngine[K comparable, V any] struct {
	*MemoryStorage[K, V]
	mu sync.RWMutex
}

// NewSetter создает новый объект Setter.
func NewMemorySetterEngine[K comparable, V any](storage *MemoryStorage[K, V]) *MemorySetterEngine[K, V] {
	return &MemorySetterEngine[K, V]{MemoryStorage: storage}
}

// Set устанавливает значение по ключу.
func (s *MemorySetterEngine[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}
