package engines

import (
	"encoding/json"
	"io"
	"sync"
)

// JsonFileWriterEngine - структура для записи данных в файл
type JsonFileWriterEngine[T any] struct {
	*JsonFileStorage[T]
	mu sync.RWMutex
}

// WriteRow - метод для записи одной строки (одного объекта JSON) в файл
func (fw *JsonFileWriterEngine[T]) WriteRow(data T) bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	_, err := fw.File.Seek(0, io.SeekEnd)
	if err != nil {
		return false
	}

	encoder := json.NewEncoder(fw.File)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	return err == nil
}
