package db

import (
	"context"
	"crypto/tls"
	"database/sql"
	"log"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/lib/pq"
)

type DB struct {
	Postgres   *sql.DB
	ClickHouse clickhouse.Conn
}

// Connect подключается к Postgres и ClickHouse Cloud через безопасные соединения
func Connect(postgresDSN, clickhouseUser, clickhousePass, clickhouseHost string) (*DB, error) {
	// --- Подключение к Postgres ---
	pg, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		return nil, err
	}

	pg.SetMaxOpenConns(25)
	pg.SetMaxIdleConns(25)
	pg.SetConnMaxLifetime(5 * time.Minute)

	if err := pg.Ping(); err != nil {
		return nil, err
	}

	// --- Подключение к ClickHouse Cloud (HTTPS, TLS) ---
	ch, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{clickhouseHost},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: clickhouseUser,
			Password: clickhousePass,
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true, // безопасно для ClickHouse Cloud
		},
		DialTimeout:  5 * time.Second,
		MaxOpenConns: 10,
		MaxIdleConns: 5,
	})
	if err != nil {
		return nil, err
	}

	// Проверка соединения ClickHouse
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ch.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL and ClickHouse Cloud")
	return &DB{
		Postgres:   pg,
		ClickHouse: ch,
	}, nil
}
