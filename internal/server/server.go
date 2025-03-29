package server

import (
	"context"
	"go-yandex-practicum-metrics/internal/logger"
	"go-yandex-practicum-metrics/internal/router"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Addresser interface {
	GetAddress() string
}

type Router interface {
	AddHandler(method, path string, handler httprouter.Handle)
	GetRoutes() []router.Route
	ServerHTTP() http.Handler
}

type Server struct {
	server *http.Server
	router Router
}

func NewServer(a Addresser) *Server {
	logger.Info("Initializing server", "address", a.GetAddress())
	return &Server{
		server: &http.Server{
			Addr:    a.GetAddress(),
			Handler: http.NewServeMux(),
		},
		router: router.NewRouter(),
	}
}

func (s *Server) AddRouter(rtr Router) {
	logger.Info("Adding router to server")
	if s.router == nil {
		s.router = rtr
		s.server.Handler = rtr.ServerHTTP()
		logger.Info("Assigned router as the main handler")
		return
	}
	for _, r := range rtr.GetRoutes() {
		logger.Info("Adding route", "method", r.Method, "path", r.Path)
		s.router.AddHandler(r.Method, r.Path, r.Handler)
	}
	s.server.Handler = s.router.ServerHTTP()
	logger.Info("Updated server handler with new routes")
}

func (s *Server) Start(ctx context.Context) error {
	logger.Info("Starting server", "address", s.server.Addr)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error running server", "error", err)
		}
	}()
	<-ctx.Done()
	logger.Info("Shutting down gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error shutting down HTTP server", "error", err)
		return err
	}
	logger.Info("Shutdown complete.")
	return nil
}
