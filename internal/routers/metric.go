package routers

import (
	"go-yandex-practicum-metrics/internal/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Logger interface {
	Infow(msg string, args ...any)
}

func RegisterMetricUpdatePathRoute(
	router chi.Router,
	handler http.HandlerFunc,
	logger Logger,
) {
	subRouter := setupMetricRouter(logger)
	subRouter.Post("/update/{type}/{id}/{value}", handler)
	addMetricRouter(router, subRouter)
}

func RegisterMetricUpdateBodyRoute(
	router chi.Router,
	handler http.HandlerFunc,
	logger Logger,
) {
	subRouter := setupMetricRouter(logger)
	subRouter.Post("/update/", handler)
	addMetricRouter(router, subRouter)
}

func RegisterMetricUpdatesBodyRoute(
	router chi.Router,
	handler http.HandlerFunc,
	logger Logger,
) {
	subRouter := setupMetricRouter(logger)
	subRouter.Post("/updates/", handler)
	addMetricRouter(router, subRouter)
}

func RegisterMetricGetPathRoute(
	router chi.Router,
	handler http.HandlerFunc,
	logger Logger,
) {
	subRouter := setupMetricRouter(logger)
	subRouter.Get("/value/{type}/{id}", handler)
	addMetricRouter(router, subRouter)
}

func RegisterMetricGetBodyRoute(
	router chi.Router,
	handler http.HandlerFunc,
	logger Logger,
) {
	subRouter := setupMetricRouter(logger)
	subRouter.Post("/value/", handler)
	addMetricRouter(router, subRouter)
}

func RegisterMetricsListHTMLRoute(
	router chi.Router,
	handler http.HandlerFunc,
	logger Logger,
) {
	subRouter := setupMetricRouter(logger)
	subRouter.Get("/", handler)
	addMetricRouter(router, subRouter)
}

func addMetricRouter(router chi.Router, subRouter chi.Router) {
	router.Mount("/", subRouter)
}

func setupMetricRouter(logger Logger) *chi.Mux {
	subRouter := newRouter()
	useMetricMiddlewares(subRouter, logger)
	return subRouter
}

func newRouter() *chi.Mux {
	return chi.NewRouter()
}

func useMetricMiddlewares(router chi.Router, logger Logger) {
	router.Use(middlewares.GzipMiddleware(logger))
	router.Use(middlewares.LoggingMiddleware(logger))
	router.Use(middleware.Recoverer)
}
