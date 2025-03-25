package app

import (
	"database/sql"
	"go-yandex-practicum-metrics/internal/configs"
	"go-yandex-practicum-metrics/internal/engines"
	"go-yandex-practicum-metrics/internal/types"
	"net/http"
	"os"
)

type MetricServerApp struct {
	config *configs.ServerConfig
	db     *sql.DB
	file   *os.File
	memory map[types.MetricID]*types.Metrics

	loadMetricsWorker *LoadMetricsFromFileWorker
	dumpMetricsWorker *DumpMetricsToFileWorker

	srv *http.Server
}

func NewMetricServerApp(config *configs.ServerConfig) (*MetricServerApp, error) {
	db := engines.NewDBEngine(config)
	file := engines.NewFileEngine(config)

	_, err := NewStorages(config, db, file)

	if err != nil {
		return nil, err
	}

	return nil, nil

}

// 	// Инициализация памяти
// 	if config.FileStoragePath == "" && config.DatabaseDSN == "" {
// 		memory = make(map[types.MetricID]*types.Metrics)
// 	}

// 	// Репозитории для БД, файла или памяти
// 	var dbSaveRepo *repositories.MetricDBSaveRepository
// 	var dbFilterRepo *repositories.MetricDBFilterRepository
// 	var dbListRepo *repositories.MetricDBListRepository
// 	if db != nil {
// 		dbSaveRepo = repositories.NewMetricDBSaveRepository(db)
// 		dbFilterRepo = repositories.NewMetricDBFilterRepository(db)
// 		dbListRepo = repositories.NewMetricDBListRepository(db)
// 	}

// 	var fileSaveRepo *repositories.MetricFileSaveRepository
// 	var fileFilterRepo *repositories.MetricFileFilterRepository
// 	var fileListRepo *repositories.MetricFileListRepository
// 	if file != nil {
// 		fileSaveRepo = repositories.NewMetricFileSaveRepository(file)
// 		fileFilterRepo = repositories.NewMetricFileFilterRepository(file)
// 		fileListRepo = repositories.NewMetricFileListRepository(file)
// 	}

// 	var memorySaveRepo *repositories.MetricMemorySaveRepository
// 	var memoryFilterRepo *repositories.MetricMemoryFilterRepository
// 	var memoryListRepo *repositories.MetricMemoryListRepository
// 	if db == nil && file == nil {
// 		memorySaveRepo = repositories.NewMetricMemorySaveRepository(memory)
// 		memoryFilterRepo = repositories.NewMetricMemoryFilterRepository(memory)
// 		memoryListRepo = repositories.NewMetricMemoryListRepository(memory)
// 	}

// 	// Сервисы
// 	var metricUpdateService *services.MetricUpdateService
// 	var metricUpdatesService *services.MetricUpdatesService
// 	var metricGetService *services.MetricGetService
// 	var metricListService *services.MetricListService

// 	if db != nil {
// 		metricUpdateService = services.NewMetricUpdateService(dbSaveRepo, dbFilterRepo)
// 		metricUpdatesService = services.NewMetricUpdatesService(dbSaveRepo, dbFilterRepo, db)
// 		metricGetService = services.NewMetricGetService(dbFilterRepo)
// 		metricListService = services.NewMetricListService(dbListRepo)
// 	} else if file != nil {
// 		metricUpdateService = services.NewMetricUpdateService(fileSaveRepo, fileFilterRepo)
// 		metricUpdatesService = services.NewMetricUpdatesService(fileSaveRepo, fileFilterRepo, nil)
// 		metricGetService = services.NewMetricGetService(fileFilterRepo)
// 		metricListService = services.NewMetricListService(fileListRepo)
// 	} else {
// 		metricUpdateService = services.NewMetricUpdateService(memorySaveRepo, memoryFilterRepo)
// 		metricUpdatesService = services.NewMetricUpdatesService(memorySaveRepo, memoryFilterRepo, nil)
// 		metricGetService = services.NewMetricGetService(memoryFilterRepo)
// 		metricListService = services.NewMetricListService(memoryListRepo)
// 	}

// 	r := chi.NewRouter()
// 	logger.InitializeLogger(logger.InfoLevel)
// 	routers.RegisterMetricUpdatePathRoute(r, handlers.MetricUpdatePathHandler(usecases.NewMetricUpdatePathUsecase(metricUpdateService)), logger.Logger)
// 	routers.RegisterMetricUpdateBodyRoute(r, handlers.MetricUpdateBodyHandler(usecases.NewMetricUpdateBodyUsecase(metricUpdateService)), logger.Logger)
// 	routers.RegisterMetricUpdatesBodyRoute(r, handlers.MetricUpdatesBodyHandler(usecases.NewMetricUpdatesBodyUsecase(metricUpdatesService)), logger.Logger)
// 	routers.RegisterMetricGetPathRoute(r, handlers.MetricGetPathHandler(usecases.NewMetricGetPathUsecase(metricGetService)), logger.Logger)
// 	routers.RegisterMetricGetBodyRoute(r, handlers.MetricGetBodyHandler(usecases.NewMetricGetBodyUsecase(metricGetService)), logger.Logger)
// 	routers.RegisterMetricsListHTMLRoute(r, handlers.MetricListHTMLHandler(usecases.NewMetricListHTMLUsecase(metricListService)), logger.Logger)

// 	srv := &http.Server{
// 		Addr:    config.Address,
// 		Handler: r,
// 	}

// 	var loadMetricsWorker *LoadMetricsFromFileWorker
// 	var dumpMetricsWorker *DumpMetricsToFileWorker
// 	if db != nil {
// 		loadMetricsWorker = NewLoadMetricsFromFileWorker(fileListRepo, dbSaveRepo)
// 		dumpMetricsWorker = NewDumpMetricsToFileWorker(config, fileListRepo, dbSaveRepo)
// 	}

// 	return &MetricServerApp{
// 		config:            config,
// 		db:                db,
// 		file:              file,
// 		memory:            memory,
// 		loadMetricsWorker: loadMetricsWorker,
// 		dumpMetricsWorker: dumpMetricsWorker,
// 		srv:               srv,
// 	}, nil
// }

// func (app *MetricServerApp) Start(ctx context.Context) int {
// 	go func() {
// 		// Запуск HTTP сервера
// 		if err := app.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			logger.Logger.Fatalf("Server failed: %v", err)
// 		}
// 	}()

// 	// Запуск воркеров
// 	if app.loadMetricsWorker != nil {
// 		go func() {
// 			app.loadMetricsWorker.Start(ctx)
// 		}()
// 	}

// 	if app.dumpMetricsWorker != nil {
// 		go func() {
// 			app.dumpMetricsWorker.Start(ctx)
// 		}()
// 	}

// 	// Ожидаем сигнала на завершение
// 	<-ctx.Done()
// 	logger.Logger.Info("Shutting down server...")

// 	// Создаем контекст для завершения с таймаутом
// 	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// Шutdown HTTP сервера
// 	if err := app.srv.Shutdown(shutdownCtx); err != nil {
// 		logger.Logger.Errorf("Server forced to shutdown: %v", err)
// 		// Возвращаем код ошибки 1 при неудачном завершении
// 		return 1
// 	}

// 	// Закрытие соединений с базой данных, если они есть
// 	if app.db != nil {
// 		logger.Logger.Info("Closing database connection...")
// 		if err := app.db.Close(); err != nil {
// 			logger.Logger.Errorf("Error closing database connection: %v", err)
// 			// Возвращаем код ошибки 1 при ошибке закрытия базы данных
// 			return 1
// 		} else {
// 			logger.Logger.Info("Database connection closed")
// 		}
// 	}

// 	// Закрытие файла, если он был открыт
// 	if app.file != nil {
// 		logger.Logger.Info("Closing file connection...")
// 		if err := app.file.Close(); err != nil {
// 			logger.Logger.Errorf("Error closing file connection: %v", err)
// 			// Возвращаем код ошибки 1 при ошибке закрытия файла
// 			return 1
// 		} else {
// 			logger.Logger.Info("File connection closed")
// 		}
// 	}

// 	// Завершаем работу успешно
// 	logger.Logger.Info("Server exited gracefully")
// 	// Возвращаем 0, если все прошло успешно
// 	return 0
// }
