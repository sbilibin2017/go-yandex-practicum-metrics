package app

import (
	"go-yandex-practicum-metrics/internal/configs"
	"go-yandex-practicum-metrics/internal/handlers"
	"go-yandex-practicum-metrics/internal/router"
	"go-yandex-practicum-metrics/internal/routers"
	"go-yandex-practicum-metrics/internal/server"
	"go-yandex-practicum-metrics/internal/services"
	"go-yandex-practicum-metrics/internal/usecases"
)

// NewServer - создание нового сервера с зависимостями из контейнера
func NewServer(config *configs.ServerConfig, container *Container) (*server.Server, error) {
	var metricUpdateService *services.MetricUpdateService
	if container.DBEngine != nil {
		metricUpdateService = services.NewMetricUpdateService(
			container.DBSaveRepo,
			container.DBFindRepo,
			container.UnitOfWork,
		)
	} else if container.FileEngine != nil {
		metricUpdateService = services.NewMetricUpdateService(
			container.FileSaveRepo,
			container.FileFindRepo,
			container.UnitOfWork,
		)
	} else {
		metricUpdateService = services.NewMetricUpdateService(
			container.MemorySaveRepo,
			container.MemoryFindRepo,
			container.UnitOfWork,
		)
	}
	metricUpdateUsecase := usecases.NewMetricUpdatePathUsecase(metricUpdateService)
	metricUpdateHandler := handlers.MetricUpdatePathHandler(metricUpdateUsecase)
	router := router.NewRouter()
	routers.RegisterMetricUpdatePathRouter(router, metricUpdateHandler)
	srv := server.NewServer(config)
	srv.AddRouter(router)
	return srv, nil
}
