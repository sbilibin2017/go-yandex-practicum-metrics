package engines

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemorySetterEngine_Set(t *testing.T) {
	storage, _ := NewMemoryStorage[string, int]()
	setter := NewMemorySetterEngine(storage)

	setter.Set("key1", 100)
	setter.Set("key2", 200)

	assert.Equal(t, 100, storage.data["key1"], "Значение по ключу 'key1' должно быть 100")
	assert.Equal(t, 200, storage.data["key2"], "Значение по ключу 'key2' должно быть 200")
}

func TestMemorySetterEngine_ConcurrentSet(t *testing.T) {
	storage, _ := NewMemoryStorage[int, int]()
	setter := NewMemorySetterEngine(storage)

	var wg sync.WaitGroup
	n := 1000
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			setter.Set(i, i*10)
		}(i)
	}
	wg.Wait()

	for i := 0; i < n; i++ {
		assert.Equal(t, i*10, storage.data[i], "Некорректное значение для ключа")
	}
}
