package engines

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

// FileGenerator - структура для генерации данных из файла
type FileGenerator[T any] struct {
	*JsonFileStorage[T] // Структура для работы с файлом
	mu                  sync.RWMutex
}

// GenerateRow - метод для генерации данных из файла
func (fg *FileGenerator[T]) GenerateRow() (<-chan T, bool) {
	// Открываем файл с блокировкой чтения
	fg.mu.RLock()
	defer fg.mu.RUnlock()

	// Проверяем, открыт ли файл
	if fg.File == nil {
		fmt.Println("File is not opened")
		return nil, false
	}

	// Создаем канал для отправки значений
	ch := make(chan T)

	// Горутина для чтения строк из файла
	go func() {
		defer close(ch) // Закрываем канал после завершения чтения

		decoder := json.NewDecoder(fg.File)

		// Чтение файла строка за строкой
		for {
			var row T
			if err := decoder.Decode(&row); err == io.EOF {
				// Достигнут конец файла
				return
			} else if err != nil {
				// Ошибка при декодировании
				fmt.Println("Error decoding row:", err)
				return
			}
			// Отправляем строку в канал
			ch <- row
		}
	}()

	return ch, true
}
