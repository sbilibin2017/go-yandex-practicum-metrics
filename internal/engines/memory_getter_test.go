package engines

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryGetterEngine_Get(t *testing.T) {
	storage, _ := NewMemoryStorage[string, int]()
	storage.data["key1"] = 100
	storage.data["key2"] = 200

	getter := NewMemoryGetterEngine(storage)

	value, ok := getter.Get("key1")
	assert.True(t, ok, "Ключ 'key1' должен существовать")
	assert.Equal(t, 100, value, "Значение по ключу 'key1' должно быть 100")

	value, ok = getter.Get("key2")
	assert.True(t, ok, "Ключ 'key2' должен существовать")
	assert.Equal(t, 200, value, "Значение по ключу 'key2' должно быть 200")

	value, ok = getter.Get("key3")
	assert.False(t, ok, "Ключ 'key3' не должен существовать")
}

func TestMemoryGetterEngine_ConcurrentGet(t *testing.T) {
	storage, _ := NewMemoryStorage[int, int]()
	for i := 0; i < 1000; i++ {
		storage.data[i] = i * 10
	}

	getter := NewMemoryGetterEngine(storage)

	var wg sync.WaitGroup
	n := 1000
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			value, ok := getter.Get(i)
			assert.True(t, ok, "Ключ должен существовать")
			assert.Equal(t, i*10, value, "Некорректное значение для ключа")
		}(i)
	}
	wg.Wait()
}
