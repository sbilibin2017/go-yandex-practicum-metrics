package engines

import (
	"encoding/json"
	"io"
	"sync"
)

// JsonFileSaverEngine - структура для добавления данных в файл
type JsonFileSaverEngine[T any] struct {
	storage *JsonFileStorage[T]
	mu      sync.RWMutex
}

// Save - метод для добавления данных в конец файла, не очищая его
func (fs *JsonFileSaverEngine[T]) Save(data []T) bool {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	_, err := fs.storage.File.Seek(0, io.SeekEnd)
	if err != nil {
		return false
	}

	encoder := json.NewEncoder(fs.storage.File)
	encoder.SetIndent("", "  ")
	for _, row := range data {
		if err := encoder.Encode(row); err != nil {
			return false
		}
	}

	return true
}
