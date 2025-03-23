package storages

import (
	"os"
)

// FileStorageConfig интерфейс для получения пути к файлу
type FileStorageConfig interface {
	GetFileStoragePath() string
}

// JsonFileEngine - основная структура, управляющая файлом
type JsonFileStorage[T any] struct {
	*os.File
}

// NewJsonFileEngine - конструктор для создания нового экземпляра JsonFileEngine
func NewJsonFileStorage[T any](e FileStorageConfig) (*JsonFileStorage[T], bool) {
	file, err := os.OpenFile(e.GetFileStoragePath(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, false
	}
	return &JsonFileStorage[T]{File: file}, true
}
