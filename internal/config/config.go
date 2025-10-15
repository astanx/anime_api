package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress  string
	PostgresDSN    string
	ClickhouseUser string
	ClickhouseHost string
	ClickhousePass string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}
	return &Config{
		ServerAddress:  addr,
		PostgresDSN:    os.Getenv("POSTGRES_DSN"),
		ClickhouseUser: os.Getenv("CLICKHOUSE_USERNAME"),
		ClickhouseHost: os.Getenv("CLICKHOUSE_HOST"),
		ClickhousePass: os.Getenv("CLICKHOUSE_PASSWORD"),
	}, nil
}
