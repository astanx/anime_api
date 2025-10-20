package repository

import (
	"database/sql"
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
	_, err := r.db.Exec(
		`INSERT INTO history (device_id, anime_id, last_watched, is_watched)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (device_id, anime_id) DO UPDATE
		 SET last_watched = EXCLUDED.last_watched`,
		deviceID, history.AnimeID, history.LastWatchedEpisode, history.IsWatched,
	)
	return err
}

func (r *HistoryRepo) GetAllHistory(deviceID string) ([]model.History, error) {
	rows, err := r.db.Query(
		"SELECT anime_id, last_watched, is_watched FROM history WHERE device_id = $1 ORDER BY anime_id",
		deviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var historyList []model.History
	for rows.Next() {
		var h model.History
		if err := rows.Scan(&h.AnimeID, &h.LastWatchedEpisode); err != nil {
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

	rows, err := r.db.Query(
		"SELECT anime_id, last_watched, is_watched FROM history WHERE device_id = $1 ORDER BY anime_id LIMIT $2 OFFSET $3",
		deviceID, limit, offset,
	)
	if err != nil {
		return model.PaginatedHistory{}, err
	}
	defer rows.Close()

	var historyList []model.History
	for rows.Next() {
		var h model.History
		if err := rows.Scan(&h.AnimeID, &h.LastWatchedEpisode); err != nil {
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
