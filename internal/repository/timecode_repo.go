package repository

import (
	"database/sql"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/astanx/anime_api/internal/db"
	"github.com/astanx/anime_api/internal/model"
)

type TimecodeRepo struct {
	dbPostgres   *sql.DB
	dbClickhouse clickhouse.Conn
}

func NewTimecodeRepo(db *db.DB) *TimecodeRepo {
	return &TimecodeRepo{
		dbPostgres:   db.Postgres,
		dbClickhouse: db.ClickHouse,
	}
}

func (r *TimecodeRepo) GetAllTimecodes(deviceID string) ([]model.Timecode, error) {
	rows, err := r.dbPostgres.Query(
		"SELECT time, episode_id, is_watched, device_id FROM timecodes WHERE device_id = $1",
		deviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timecodes []model.Timecode

	for rows.Next() {
		var t model.Timecode
		if err := rows.Scan(&t.Time, &t.EpisodeID, &t.IsWatched); err != nil {
			return nil, err
		}
		timecodes = append(timecodes, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return timecodes, nil
}

func (r *TimecodeRepo) GetTimecode(deviceID string, episodeID string) (*model.Timecode, error) {
	row := r.dbPostgres.QueryRow(
		"SELECT time, episode_id, is_watched, device_id FROM timecodes WHERE device_id = $1 AND episode_id = $2 LIMIT 1",
		deviceID, episodeID,
	)

	var t model.Timecode
	err := row.Scan(&t.Time, &t.EpisodeID, &t.IsWatched)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}

func (r *TimecodeRepo) AddTimecode(deviceID string, timecode model.Timecode) error {
	_, err := r.dbPostgres.Exec(
		`INSERT INTO timecodes (time, episode_id, is_watched, device_id)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (device_id, episode_id) DO UPDATE
		 SET time = EXCLUDED.time,
		     is_watched = EXCLUDED.is_watched`,
		timecode.Time, timecode.EpisodeID, timecode.IsWatched, deviceID,
	)
	return err
}
