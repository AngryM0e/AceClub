package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AngryM0e/AceClub/Backend/config"
	"github.com/AngryM0e/AceClub/Backend/internal/repository/postgres"
	"github.com/AngryM0e/AceClub/Backend/internal/transport/handlers"
	"github.com/AngryM0e/AceClub/Backend/internal/transport/middleware"
)

type Server struct {
	httpServer  *http.Server
	userHandler *handlers.UserHandler
}

func New(cfg *config.Config) (*Server, error) {
	db, err := postgres.NewDB(cfg, cfg.ConnStr())
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	userRepo := postgres.NewUserRepository(db.DB)
	userHandler := handlers.NewUserHandler(userRepo)

	mux := http.NewServeMux()
	srv := &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.Server.Port,
			Handler:      middleware.Logging(mux),
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
		userHandler: userHandler,
	}
	srv.RegisterRoutes(mux)
	return srv, nil
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	handle := func(h APIHandler) http.Handler {
		return ErrorAdapter(h)
	}

	mux.Handle("GET /health", handle(handlers.HealthCheck))
	mux.Handle("POST /api/users", handle(s.userHandler.Create))

}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
