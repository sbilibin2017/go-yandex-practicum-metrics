package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для DB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Open() error {
	args := m.Called()
	return args.Error(0)
}

// Мок для File
type MockFile struct {
	mock.Mock
}

func (m *MockFile) Open() error {
	args := m.Called()
	return args.Error(0)
}

// Мок для Config
type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) GetDatabaseDSN() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfig) GetFileStoragePath() string {
	args := m.Called()
	return args.String(0)
}

func TestNewStorages_DBAndFileProvided(t *testing.T) {
	config := new(MockConfig)
	config.On("GetDatabaseDSN").Return("some-dsn")
	config.On("GetFileStoragePath").Return("/path/to/file")

	db := new(MockDB)
	db.On("Open").Return(nil)

	file := new(MockFile)
	file.On("Open").Return(nil)

	storages, err := NewStorages(config, db, file)

	assert.NoError(t, err)
	assert.Nil(t, storages.Memory)
	assert.NotNil(t, storages.DB)
	assert.NotNil(t, storages.File)

	db.AssertExpectations(t)
	file.AssertExpectations(t)
}

func TestNewStorages_OnlyDBProvided(t *testing.T) {
	config := new(MockConfig)
	config.On("GetDatabaseDSN").Return("some-dsn")
	config.On("GetFileStoragePath").Return("")

	db := new(MockDB)
	db.On("Open").Return(nil)

	storages, err := NewStorages(config, db, nil)

	assert.NoError(t, err)
	assert.Nil(t, storages.Memory)
	assert.NotNil(t, storages.DB)
	assert.Nil(t, storages.File)

	db.AssertExpectations(t)
}

func TestNewStorages_OnlyFileProvided(t *testing.T) {
	config := new(MockConfig)
	config.On("GetDatabaseDSN").Return("")
	config.On("GetFileStoragePath").Return("some-file-path")

	file := new(MockFile)
	file.On("Open").Return(nil)

	storages, err := NewStorages(config, nil, file)

	assert.NoError(t, err)
	assert.Nil(t, storages.Memory)
	assert.Nil(t, storages.DB)
	assert.NotNil(t, storages.File)

	file.AssertExpectations(t)
}

func TestNewStorages_DBError(t *testing.T) {
	config := new(MockConfig)
	config.On("GetDatabaseDSN").Return("some-dsn")
	config.On("GetFileStoragePath").Return("")

	db := new(MockDB)
	db.On("Open").Return(errors.New("db error"))

	storages, err := NewStorages(config, db, nil)

	assert.EqualError(t, err, "db error")
	assert.Nil(t, storages.DB)
	assert.Nil(t, storages.File)
	assert.Nil(t, storages.Memory)

	db.AssertExpectations(t)
}

func TestNewStorages_FileError(t *testing.T) {
	config := new(MockConfig)
	config.On("GetDatabaseDSN").Return("")
	config.On("GetFileStoragePath").Return("/path/to/file")

	file := new(MockFile)
	file.On("Open").Return(errors.New("file error"))

	storages, err := NewStorages(config, nil, file)

	assert.EqualError(t, err, "file error")
	assert.Nil(t, storages.DB)
	assert.Nil(t, storages.File)
	assert.Nil(t, storages.Memory)

	file.AssertExpectations(t)
}

func TestNewStorages_NoDBNoFile(t *testing.T) {
	config := new(MockConfig)
	config.On("GetDatabaseDSN").Return("")
	config.On("GetFileStoragePath").Return("")

	storages, err := NewStorages(config, nil, nil)

	assert.NoError(t, err)
	assert.Nil(t, storages.DB)
	assert.Nil(t, storages.File)
	assert.NotNil(t, storages.Memory) // Память должна быть инициализирована
}
