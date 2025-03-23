package engines

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryRangerEngine_Range(t *testing.T) {
	storage, _ := NewMemoryStorage[string, int]()
	storage.data["key1"] = 100
	storage.data["key2"] = 200
	storage.data["key3"] = 300

	ranger := NewMemoryRangerEngine(storage)

	collected := make(map[string]int)
	ranger.Range(func(key string, value int) bool {
		collected[key] = value
		return true
	})

	assert.Equal(t, storage.data, collected, "Все элементы должны быть перебраны корректно")
}

func TestMemoryRangerEngine_RangeStopEarly(t *testing.T) {
	storage, _ := NewMemoryStorage[string, int]()
	storage.data["key1"] = 100
	storage.data["key2"] = 200
	storage.data["key3"] = 300

	ranger := NewMemoryRangerEngine(storage)

	count := 0
	ranger.Range(func(key string, value int) bool {
		count++
		return count < 2 // Остановить после второго элемента
	})

	assert.LessOrEqual(t, count, 2, "Range должен остановиться после второго элемента")
}
