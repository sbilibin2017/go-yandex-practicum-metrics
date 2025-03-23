package engines

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage[T any] struct {
	file *os.File
}

func (m *mockStorage[T]) File() *os.File {
	return m.file
}

func TestJsonFileWriterEngine_WriteRow(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "test_storage.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	storage := &JsonFileStorage[map[string]string]{File: tempFile}
	sut := &JsonFileWriterEngine[map[string]string]{JsonFileStorage: storage}

	data := map[string]string{"key": "value"}
	assert.True(t, sut.WriteRow(data))

	content, err := ioutil.ReadFile(tempFile.Name())
	require.NoError(t, err)

	var result map[string]string
	require.NoError(t, json.Unmarshal(content, &result))
	assert.Equal(t, data, result)
}
