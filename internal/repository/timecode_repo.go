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
		"SELECT time, episode_id, is_watched, anime_id FROM timecodes WHERE device_id = $1",
		deviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	timecodes := make([]model.Timecode, 0)
	for rows.Next() {
		var t model.Timecode
		if err := rows.Scan(&t.Time, &t.EpisodeID, &t.IsWatched, &t.AnimeID); err != nil {
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
		"SELECT time, episode_id, is_watched, anime_id FROM timecodes WHERE device_id = $1 AND episode_id = $2 LIMIT 1",
		deviceID, episodeID,
	)

	var t model.Timecode
	err := row.Scan(&t.Time, &t.EpisodeID, &t.IsWatched, &t.AnimeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}

func (r *TimecodeRepo) AddTimecode(deviceID string, timecode model.Timecode) error {
	var exists bool
	err := r.dbPostgres.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM timecodes WHERE device_id=$1 AND episode_id=$2)`,
		deviceID, timecode.EpisodeID,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = r.dbPostgres.Exec(
			`UPDATE timecodes
             SET time=$1, is_watched=$2
             WHERE device_id=$3 AND episode_id=$4`,
			timecode.Time, timecode.IsWatched, deviceID, timecode.EpisodeID,
		)
	} else {
		_, err = r.dbPostgres.Exec(
			`INSERT INTO timecodes (time, episode_id, is_watched, device_id, anime_id)
             VALUES ($1, $2, $3, $4, $5)`,
			timecode.Time, timecode.EpisodeID, timecode.IsWatched, deviceID, timecode.AnimeID,
		)
	}
	return err
}

func (r *TimecodeRepo) GetTimecodesForAnime(deviceID string, animeID string) ([]model.Timecode, error) {
	rows, err := r.dbPostgres.Query(
		"SELECT time, episode_id, is_watched, anime_id FROM timecodes WHERE device_id = $1 AND anime_id = $2",
		deviceID, animeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timecodes []model.Timecode

	for rows.Next() {
		var t model.Timecode
		if err := rows.Scan(&t.Time, &t.EpisodeID, &t.IsWatched, &t.AnimeID); err != nil {
			return nil, err
		}
		timecodes = append(timecodes, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return timecodes, nil
}
