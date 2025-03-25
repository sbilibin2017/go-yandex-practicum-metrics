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

func RegisterUpdateMetricPathRoute(
	router chi.Router,
	handler http.HandlerFunc,
	logger Logger,
) {
	subRouter := setupMetricRouter(logger)
	subRouter.Post("/update/{type}/{id}/{value}", handler)
	router.Mount("/", subRouter)
}

func RegisterUpdateMetricBodyRoute(
	router chi.Router,
	handler http.HandlerFunc,
	logger Logger,
) {
	subRouter := setupMetricRouter(logger)
	subRouter.Post("/update/", handler)
	router.Mount("/", subRouter)
}

func RegisterUpdatesMetricBodyRoute(
	router chi.Router,
	handler http.HandlerFunc,
	logger Logger,
) {
	subRouter := setupMetricRouter(logger)
	subRouter.Post("/updates/", handler)
	router.Mount("/", subRouter)
}

func RegisterGetMetricByTypeAndIDPathRoute(
	router chi.Router,
	handler http.HandlerFunc,
	logger Logger,
) {
	subRouter := setupMetricRouter(logger)
	subRouter.Get("/value/{type}/{id}", handler)
	addMetricRouter(router, subRouter)
}

func RegisterGetMetricByByTypeAndIDBodyRoute(
	router chi.Router,
	handler http.HandlerFunc,
	logger Logger,
) {
	subRouter := setupMetricRouter(logger)
	subRouter.Post("/value/", handler)
	addMetricRouter(router, subRouter)
}

func RegisterListMetricsHTMLRoute(
	router chi.Router,
	handler http.HandlerFunc,
	logger Logger,
) {
	subRouter := setupMetricRouter(logger)
	subRouter.Get("/", handler)
	addMetricRouter(router, subRouter)
}

func setupMetricRouter(logger Logger) *chi.Mux {
	subRouter := newRouter()
	useMetricMiddlewares(subRouter, logger)
	return subRouter
}

func newRouter() *chi.Mux {
	return chi.NewRouter()
}

func addMetricRouter(router chi.Router, subRouter chi.Router) {
	router.Mount("/", subRouter)
}

func useMetricMiddlewares(router chi.Router, logger Logger) {
	router.Use(middlewares.GzipMiddleware(logger))
	router.Use(middlewares.LoggingMiddleware(logger))
	router.Use(middleware.Recoverer)
}
