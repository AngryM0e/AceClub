package server

import (
	"context"
	"net/http"
	"time"

	"github.com/AngryM0e/AceClub/Backend/internal/transport/handlers"
	"github.com/AngryM0e/AceClub/Backend/internal/transport/middleware"
)

type Server struct {
	httpServer *http.Server
	cfg        Config
}

type Config struct {
	Port         string
	Host         string
	Name         string
	Password     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func New(cfg Config) (*Server, error) {
	mux := http.NewServeMux()
	srv := &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      middleware.Logging(mux),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
	srv.RegisterRoutes(mux)
	return srv, nil
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	handle := func(h APIHandler) http.Handler {
		return ErrorAdapter(h)
	}

	mux.Handle("GET /health", handle(handlers.HealthCheck))
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
