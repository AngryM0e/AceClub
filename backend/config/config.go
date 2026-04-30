package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	DB         DBConfig
	BCryptCost int
}

type ServerConfig struct {
	Port         string // API_PORT for server
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: .env file not found, using system enviroment variables")
	}
	port := os.Getenv("PORT")
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASSWORD")
	dbHost := os.Getenv("DBHOST")
	dbPort := os.Getenv("DBPORT")
	dbName := os.Getenv("DBNAME")
	readTO := time.Duration(10 * time.Second)
	writeTO := time.Duration(5 * time.Second)
	BCryptCost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		return nil, err
	}
	return &Config{
		Server: ServerConfig{
			Port:         port,
			ReadTimeout:  readTO,
			WriteTimeout: writeTO,
		},
		DB: DBConfig{
			User:     dbUser,
			Password: dbPass,
			Host:     dbHost,
			Port:     dbPort,
			Name:     dbName,
		},
		BCryptCost: BCryptCost,
	}, nil
}

func (c Config) ConnStr() string {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name)
	return connStr
}
