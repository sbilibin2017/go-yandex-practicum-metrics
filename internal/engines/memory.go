package engines

// Storage - это структура, которая хранит данные.
type MemoryStorage[K comparable, V any] struct {
	data map[K]V
}

// NewStorage создает новый объект Storage.
func NewMemoryStorage[K comparable, V any]() (*MemoryStorage[K, V], bool) {
	return &MemoryStorage[K, V]{
		data: make(map[K]V),
	}, true
}
