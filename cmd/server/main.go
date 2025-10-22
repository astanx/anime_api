package main

import (
	"log"

	"github.com/astanx/anime_api/internal/config"
	"github.com/astanx/anime_api/internal/db"
	"github.com/astanx/anime_api/internal/router"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	databases, err := db.Connect(cfg.PostgresDSN, cfg.ClickhouseUser, cfg.ClickhousePass, cfg.ClickhouseHost, cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to connect databases: %v", err)
	}

	r := router.NewRouter(databases)

	log.Printf("starting server on %s", cfg.ServerAddress)
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
