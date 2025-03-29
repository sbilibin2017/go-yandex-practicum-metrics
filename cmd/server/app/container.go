package app

import (
	"go-yandex-practicum-metrics/internal/configs"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/engines"
	"go-yandex-practicum-metrics/internal/repositories"
	"go-yandex-practicum-metrics/internal/router"
	"go-yandex-practicum-metrics/internal/unitofwork"
)

type Container struct {
	DBEngine       *engines.DBEngine
	FileEngine     *engines.FileEngine
	UnitOfWork     *unitofwork.UnitOfWork
	Router         *router.Router
	DBSaveRepo     *repositories.MetricDBSaveRepository
	DBFindRepo     *repositories.MetricDBFindRepository
	FileSaveRepo   *repositories.MetricFileSaveRepository
	FileFindRepo   *repositories.MetricFileFindRepository
	MemorySaveRepo *repositories.MetricMemorySaveRepository
	MemoryFindRepo *repositories.MetricMemoryFindRepository
}

func NewContainer(config *configs.ServerConfig) (*Container, error) {
	var db *engines.DBEngine
	var err error
	if config.GetDatabaseDSN() != "" {
		db, err = engines.NewDBEngine(config)
		if err != nil {
			return nil, err
		}
	}

	var file *engines.FileEngine
	if config.GetFileStoragePath() != "" {
		file, err = engines.NewFileEngine(config)
		if err != nil {
			return nil, err
		}
	}

	var dbSaveRepo *repositories.MetricDBSaveRepository
	var dbFindRepo *repositories.MetricDBFindRepository
	if db != nil {
		dbSaveRepo = repositories.NewMetricDBSaveRepository(db)
		dbFindRepo = repositories.NewMetricDBFindRepository(db)
	}

	var fileSaveRepo *repositories.MetricFileSaveRepository
	var fileFindRepo *repositories.MetricFileFindRepository
	if file != nil {
		fileSaveRepo = repositories.NewMetricFileSaveRepository(file)
		fileFindRepo = repositories.NewMetricFileFindRepository(file)
	}

	var memorySaveRepo *repositories.MetricMemorySaveRepository
	var memoryFindRepo *repositories.MetricMemoryFindRepository
	if db == nil && file == nil {
		data := make(map[domain.MetricID]*domain.Metrics)
		memorySaveRepo = repositories.NewMetricMemorySaveRepository(data)
		memoryFindRepo = repositories.NewMetricMemoryFindRepository(data)
	}

	uow := unitofwork.NewUnitOfWork(nil)
	if db != nil {
		uow = unitofwork.NewUnitOfWork(db)
	}

	return &Container{
		DBEngine:       db,
		FileEngine:     file,
		UnitOfWork:     uow,
		DBSaveRepo:     dbSaveRepo,
		DBFindRepo:     dbFindRepo,
		FileSaveRepo:   fileSaveRepo,
		FileFindRepo:   fileFindRepo,
		MemorySaveRepo: memorySaveRepo,
		MemoryFindRepo: memoryFindRepo,
	}, nil
}
