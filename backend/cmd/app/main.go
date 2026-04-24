package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AngryM0e/AceClub/Backend/config"
	"github.com/AngryM0e/AceClub/Backend/internal/repository/postgres"
	"github.com/AngryM0e/AceClub/Backend/internal/repository/postgres/pgutils"
	"github.com/AngryM0e/AceClub/Backend/internal/transport/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("error with loading config: %v", err)
	}
	// Create context
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM)
	defer stop()

	// Create server
	server, err := server.New(cfg)
	if err != nil {
		log.Fatalf("error with create server: %v", err)
	}

	db, err := postgres.NewDB(cfg, cfg.ConnStr())
	if err != nil {
		log.Fatalf("error with start db: %v", err)
	}
	defer db.Close()

	err = pgutils.RunMigrations("migrations", cfg.ConnStr())
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// Launch server
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error with start server: %v", err)
		}
	}()
	log.Printf("server launch on %s port...", cfg.Server.Port)
	// Graceful shutdown
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Fatalf("Shutdown server: %v", err)
	}
}
