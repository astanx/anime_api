package repository

import (
	"context"
	"database/sql"
	"log"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/astanx/anime_api/internal/db"
	"github.com/astanx/anime_api/internal/model"
	"github.com/google/uuid"
)

type DeviceRepo struct {
	dbPostgres   *sql.DB
	dbClickhouse clickhouse.Conn
}

func NewDeviceRepo(db *db.DB) *DeviceRepo {
	return &DeviceRepo{
		dbPostgres:   db.Postgres,
		dbClickhouse: db.ClickHouse,
	}
}

func (r *DeviceRepo) AddDeviceID(deviceID uuid.UUID) (model.User, error) {
	var u model.User

	err := r.dbPostgres.QueryRow(
		"INSERT INTO devices (device_id) VALUES ($1) RETURNING device_id",
		deviceID,
	).Scan(&u.ID)
	if err != nil {
		return u, err
	}

	err = r.dbClickhouse.Exec(
		context.Background(),
		"INSERT INTO device_analytics (device_id, created_at) VALUES (?, now())",
		deviceID,
	)
	if err != nil {
		log.Println("ClickHouse insert failed:", err)
	}

	return u, nil
}
