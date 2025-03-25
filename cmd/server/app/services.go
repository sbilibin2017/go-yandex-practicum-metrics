package app

import (
	"context"
	"database/sql"
	"go-yandex-practicum-metrics/internal/services"
	"go-yandex-practicum-metrics/internal/types"
)

type ServiceFilterRepository interface {
	Filter(ctx context.Context, filter types.MetricID) (*types.Metrics, error)
}

type ServiceSaveRepository interface {
	Save(ctx context.Context, metric *types.Metrics) error
}

type ServiceListRepository interface {
	List(ctx context.Context) ([]*types.Metrics, error)
}

type RepoProvider interface {
	GetSaveRepository() ServiceSaveRepository
	GetFilterRepository() ServiceFilterRepository
	GetListRepository() ServiceListRepository
}

type Transaction interface {
	Begin() (*sql.Tx, error)
}

// NewServices initializes and returns the services based on the provided Repos.
func NewServices(rp RepoProvider, tx Transaction) (*services.MetricUpdateService, *services.MetricUpdatesService, *services.MetricGetService, *services.MetricListService, error) {
	var metricUpdateService *services.MetricUpdateService
	var metricUpdatesService *services.MetricUpdatesService
	var metricGetService *services.MetricGetService
	var metricListService *services.MetricListService

	metricUpdateService = services.NewMetricUpdateService(
		rp.GetSaveRepository(),
		rp.GetFilterRepository(),
	)
	metricUpdatesService = services.NewMetricUpdatesService(
		rp.GetSaveRepository(),
		rp.GetFilterRepository(),
		tx,
	)
	metricGetService = services.NewMetricGetService(rp.GetFilterRepository())
	metricListService = services.NewMetricListService(rp.GetListRepository())

	return metricUpdateService, metricUpdatesService, metricGetService, metricListService, nil
}
