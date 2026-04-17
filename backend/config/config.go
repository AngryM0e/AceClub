package config

import (
	"os"
	"time"

	"github.com/AngryM0e/AceClub/Backend/internal/transport/server"
	"github.com/joho/godotenv"
)

type Config struct {
	Server server.Config
}

func New() (*Config, error) {
	godotenv.Load()
	port := os.Getenv("PORT")

	return &Config{
		Server: server.Config{
			Port:         port,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}, nil
}
