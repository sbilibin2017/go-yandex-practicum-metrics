package app_test

import (
	"testing"

	"go-yandex-practicum-metrics/cmd/server/app"
	"go-yandex-practicum-metrics/internal/types"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Test case when DBRepo is provided
func TestNewRepositories_DBRepoProvided(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocking DBRepo
	mockDBRepo := app.NewMockDBRepo(ctrl)

	// Create the repositories with DBRepo
	repos, err := app.NewRepositories(nil, mockDBRepo, nil, nil)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, repos.DBSaveRepo)
	assert.NotNil(t, repos.DBFilterRepo)
	assert.NotNil(t, repos.DBListRepo)
	assert.Nil(t, repos.FileSaveRepo)
	assert.Nil(t, repos.FileFilterRepo)
	assert.Nil(t, repos.FileListRepo)
	assert.Nil(t, repos.MemorySaveRepo)
	assert.Nil(t, repos.MemoryFilterRepo)
	assert.Nil(t, repos.MemoryListRepo)
}

// Test case when FileRepo is provided
func TestNewRepositories_FileRepoProvided(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocking FileRepo
	mockFileRepo := app.NewMockFileRepo(ctrl)

	// Create the repositories with FileRepo
	repos, err := app.NewRepositories(nil, nil, mockFileRepo, nil)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, repos.FileSaveRepo)
	assert.NotNil(t, repos.FileFilterRepo)
	assert.NotNil(t, repos.FileListRepo)
	assert.Nil(t, repos.DBSaveRepo)
	assert.Nil(t, repos.DBFilterRepo)
	assert.Nil(t, repos.DBListRepo)
	assert.Nil(t, repos.MemorySaveRepo)
	assert.Nil(t, repos.MemoryFilterRepo)
	assert.Nil(t, repos.MemoryListRepo)
}

// Test case when memory map is provided
func TestNewRepositories_MemoryProvided(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock memory map
	mockMemory := make(map[types.MetricID]*types.Metrics)

	// Create the repositories with memory
	repos, err := app.NewRepositories(nil, nil, nil, mockMemory)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, repos.MemorySaveRepo)
	assert.NotNil(t, repos.MemoryFilterRepo)
	assert.NotNil(t, repos.MemoryListRepo)
	assert.Nil(t, repos.DBSaveRepo)
	assert.Nil(t, repos.DBFilterRepo)
	assert.Nil(t, repos.DBListRepo)
	assert.Nil(t, repos.FileSaveRepo)
	assert.Nil(t, repos.FileFilterRepo)
	assert.Nil(t, repos.FileListRepo)
}

func TestGetFilterRepository_FileFilterRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock for FilterRepository (specifically for the FileFilterRepo)
	mockFileFilterRepo := app.NewMockFilterRepository(ctrl)

	// Initialize Repos with only the FileFilterRepo set
	repos := &app.Repos{
		FileFilterRepo: mockFileFilterRepo, // Only FileFilterRepo is set here
	}

	// Call the method to get the FilterRepository
	filterRepo := repos.GetFilterRepository()

	// Assert that the correct repository is returned (i.e., the mock passed in)
	assert.NotNil(t, filterRepo, "Expected FileFilterRepo to be returned")

	// Explicitly check if the returned repo is the mock we passed in
	assert.Equal(t, mockFileFilterRepo, filterRepo, "Expected the returned repo to be the mock FileFilterRepo")
}

// Test case when no repositories are provided
func TestNewRepositories_NoRepoProvided(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the repositories with no data
	repos, err := app.NewRepositories(nil, nil, nil, nil)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, repos)
}

func TestGetSaveRepository_DBRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для DB репозитория
	mockDBSaveRepo := app.NewMockSaveRepository(ctrl)

	// Инициализируем repos с DB репозиторием
	repos := &app.Repos{
		DBSaveRepo: mockDBSaveRepo,
	}

	// Получаем репозиторий через метод
	saveRepo := repos.GetSaveRepository()

	// Проверяем, что репозиторий вернулся именно из DB
	assert.Equal(t, mockDBSaveRepo, saveRepo)
}

func TestGetSaveRepository_FileRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для File репозитория
	mockFileSaveRepo := app.NewMockSaveRepository(ctrl)

	// Инициализируем repos с File репозиторием
	repos := &app.Repos{
		FileSaveRepo: mockFileSaveRepo,
	}

	// Получаем репозиторий через метод
	saveRepo := repos.GetSaveRepository()

	// Проверяем, что репозиторий вернулся именно из File
	assert.Equal(t, mockFileSaveRepo, saveRepo)
}

func TestGetSaveRepository_MemoryRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для Memory репозитория
	mockMemorySaveRepo := app.NewMockSaveRepository(ctrl)

	// Инициализируем repos с Memory репозиторием
	repos := &app.Repos{
		MemorySaveRepo: mockMemorySaveRepo,
	}

	// Получаем репозиторий через метод
	saveRepo := repos.GetSaveRepository()

	// Проверяем, что репозиторий вернулся именно из Memory
	assert.Equal(t, mockMemorySaveRepo, saveRepo)
}

func TestGetFilterRepository_DBRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для DB репозитория
	mockDBFilterRepo := app.NewMockFilterRepository(ctrl)

	// Инициализируем repos с DB репозиторием
	repos := &app.Repos{
		DBFilterRepo: mockDBFilterRepo,
	}

	// Получаем репозиторий через метод
	filterRepo := repos.GetFilterRepository()

	// Проверяем, что репозиторий вернулся именно из DB
	assert.Equal(t, mockDBFilterRepo, filterRepo)
}

func TestGetFilterRepository_MemoryRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для Memory репозитория
	mockMemoryFilterRepo := app.NewMockFilterRepository(ctrl)

	// Инициализируем repos с Memory репозиторием
	repos := &app.Repos{
		MemoryFilterRepo: mockMemoryFilterRepo,
	}

	// Получаем репозиторий через метод
	filterRepo := repos.GetFilterRepository()

	// Проверяем, что репозиторий вернулся именно из Memory
	assert.Equal(t, mockMemoryFilterRepo, filterRepo)
}

func TestGetListRepository_DBRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для DB репозитория
	mockDBListRepo := app.NewMockListRepository(ctrl)

	// Инициализируем repos с DB репозиторием
	repos := &app.Repos{
		DBListRepo: mockDBListRepo,
	}

	// Получаем репозиторий через метод
	listRepo := repos.GetListRepository()

	// Проверяем, что репозиторий вернулся именно из DB
	assert.Equal(t, mockDBListRepo, listRepo)
}

func TestGetListRepository_FileRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для File репозитория
	mockFileListRepo := app.NewMockListRepository(ctrl)

	// Инициализируем repos с File репозиторием
	repos := &app.Repos{
		FileListRepo: mockFileListRepo,
	}

	// Получаем репозиторий через метод
	listRepo := repos.GetListRepository()

	// Проверяем, что репозиторий вернулся именно из File
	assert.Equal(t, mockFileListRepo, listRepo)
}

func TestGetListRepository_MemoryRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки для Memory репозитория
	mockMemoryListRepo := app.NewMockListRepository(ctrl)

	// Инициализируем repos с Memory репозиторием
	repos := &app.Repos{
		MemoryListRepo: mockMemoryListRepo,
	}

	// Получаем репозиторий через метод
	listRepo := repos.GetListRepository()

	// Проверяем, что репозиторий вернулся именно из Memory
	assert.Equal(t, mockMemoryListRepo, listRepo)
}
