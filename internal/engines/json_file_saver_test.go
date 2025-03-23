package engines

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonFileSaverEngine_Save(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "test_storage.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	storage := &JsonFileStorage[map[string]string]{File: tempFile}
	sut := &JsonFileSaverEngine[map[string]string]{storage: storage}

	data := []map[string]string{
		{"key1": "value1"},
		{"key2": "value2"},
	}
	assert.True(t, sut.Save(data))

	// Закрываем файл перед чтением
	require.NoError(t, tempFile.Close())

	// Переоткрываем файл для чтения
	tempFile, err = os.Open(tempFile.Name())
	require.NoError(t, err)
	defer tempFile.Close()

	var results []map[string]string
	decoder := json.NewDecoder(tempFile)
	for {
		var record map[string]string
		if err := decoder.Decode(&record); err != nil {
			break
		}
		results = append(results, record)
	}

	assert.Equal(t, data, results)
}
