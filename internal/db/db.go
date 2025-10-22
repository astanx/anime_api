package db

import (
	"crypto/tls"
	"database/sql"
	"log"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	Postgres   *sql.DB
	ClickHouse clickhouse.Conn
}

func Connect(postgresDSN, clickhouseUser, clickhousePass, clickhouseHost string) (*DB, error) {
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
	}, nil
}
