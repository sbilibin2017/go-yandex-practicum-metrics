package engines

import (
	"encoding/json"
	"io"
	"sync"
)

// JsonFileReaderEngine - структура для чтения данных из файла
type JsonFileReaderEngine[T any] struct {
	storage *JsonFileStorage[T]
	mu      sync.RWMutex
}

// ReadRow - метод для чтения одной строки (одного объекта JSON) из файла
func (fr *JsonFileReaderEngine[T]) ReadRow() (*T, bool) {
	fr.mu.RLock() // Чтение, блокируем на чтение
	defer fr.mu.RUnlock()

	decoder := json.NewDecoder(fr.storage.File)
	var data T
	if err := decoder.Decode(&data); err == io.EOF {
		return nil, false
	} else if err != nil {
		return nil, false
	}

	return &data, true
}
