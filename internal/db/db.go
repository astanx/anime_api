package db

import (
	"context"
	"crypto/tls"
	"database/sql"
	"log"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	Postgres   *sql.DB
	ClickHouse clickhouse.Conn
	Redis      *redis.Client
}

func Connect(postgresDSN, clickhouseUser, clickhousePass, clickhouseHost, redisURL string) (*DB, error) {
	pg, err := sql.Open("pgx", postgresDSN)
	if err != nil {
		return nil, err
	}

	pg.SetMaxOpenConns(25)
	pg.SetMaxIdleConns(25)
	pg.SetConnMaxLifetime(5 * time.Minute)

	if err := pg.Ping(); err != nil {
		return nil, err
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("failed to parse Redis URL: %v", err)
	}

	rdb := redis.NewClient(opt)
	ctx := context.Background()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to ping Redis: %v", err)
	}
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connected:", pong)

	ch, _ := clickhouse.Open(&clickhouse.Options{
		Addr: []string{clickhouseHost},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: clickhouseUser,
			Password: clickhousePass,
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialTimeout:  5 * time.Second,
		MaxOpenConns: 10,
		MaxIdleConns: 5,
	})

	log.Println("Connected to PostgreSQL and ClickHouse Cloud")
	return &DB{
		Postgres:   pg,
		ClickHouse: ch,
		Redis:      rdb,
	}, nil
}
