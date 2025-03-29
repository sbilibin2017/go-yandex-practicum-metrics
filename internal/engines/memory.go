package engines

func NewMemoryEngine[K comparable, V any]() map[K]V {
	return make(map[K]V)
}
