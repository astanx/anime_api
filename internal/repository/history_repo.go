package repository

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/astanx/anime_api/internal/db"
	"github.com/astanx/anime_api/internal/model"
)

type HistoryRepo struct {
	db *sql.DB
}

func NewHistoryRepo(db *db.DB) *HistoryRepo {
	return &HistoryRepo{
		db: db.Postgres,
	}
}

func (r *HistoryRepo) AddHistory(deviceID string, history model.History) error {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM history WHERE device_id=$1 AND anime_id=$2)`,
		deviceID, history.AnimeID,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = r.db.Exec(
			`UPDATE history
             SET last_watched=$1, is_watched=$2, watched_at=now()
             WHERE device_id=$3 AND anime_id=$4`,
			history.LastWatchedEpisode, history.IsWatched, deviceID, history.AnimeID,
		)
	} else {
		_, err = r.db.Exec(
			`INSERT INTO history (device_id, anime_id, last_watched, is_watched, watched_at) VALUES ($1, $2, $3, $4, now())`,
			deviceID, history.AnimeID, history.LastWatchedEpisode, history.IsWatched,
		)
	}

	return err
}

func (r *HistoryRepo) GetAllHistory(deviceID string) ([]model.History, error) {
	rows, err := r.db.Query(
		"SELECT anime_id, last_watched, is_watched, watched_at FROM history WHERE device_id = $1 ORDER BY watched_at DESC",
		deviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	historyList := make([]model.History, 0)
	for rows.Next() {
		var h model.History
		if err := rows.Scan(&h.AnimeID, &h.LastWatchedEpisode, &h.IsWatched, &h.WatchedAt); err != nil {
			return nil, err
		}
		historyList = append(historyList, h)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return historyList, nil
}

func (r *HistoryRepo) GetHistory(deviceID string, page, limit int) (model.PaginatedHistory, error) {
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM history WHERE device_id = $1",
		deviceID,
	).Scan(&total)
	if err != nil {
		return model.PaginatedHistory{}, err
	}

	query := "SELECT anime_id, last_watched, is_watched, watched_at FROM history WHERE device_id = $1 ORDER BY watched_at DESC LIMIT $2 OFFSET $3"
	rows, err := r.db.Query(query, deviceID, limit, offset)
	if err != nil {
		fmt.Printf("Query error: %v\n", err)
		return model.PaginatedHistory{}, err
	}
	defer rows.Close()

	historyList := make([]model.History, 0)
	for rows.Next() {
		var h model.History
		if err := rows.Scan(&h.AnimeID, &h.LastWatchedEpisode, &h.IsWatched, &h.WatchedAt); err != nil {
			return model.PaginatedHistory{}, err
		}
		historyList = append(historyList, h)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	pagesLeft := int(math.Max(0, float64(totalPages-page)))

	return model.PaginatedHistory{
		Data: historyList,
		Meta: model.PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			PagesLeft:  pagesLeft,
		},
	}, nil
}
