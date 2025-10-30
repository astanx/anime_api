package repository

import (
	"context"
	"database/sql"
	"log"
	"math"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/astanx/anime_api/internal/db"
	"github.com/astanx/anime_api/internal/model"
)

type FavouriteRepo struct {
	dbPostgres   *sql.DB
	dbClickhouse clickhouse.Conn
}

func NewFavouriteRepo(db *db.DB) *FavouriteRepo {
	return &FavouriteRepo{
		dbPostgres:   db.Postgres,
		dbClickhouse: db.ClickHouse,
	}
}

func (r *FavouriteRepo) AddFavourite(deviceID string, favourite model.Favourite) error {
	var exists bool
	err := r.dbPostgres.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM favourites WHERE device_id=$1 AND anime_id=$2)`,
		deviceID, favourite.AnimeID,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	} else {
		_, err = r.dbPostgres.Exec(
			`INSERT INTO favourites (device_id, anime_id) VALUES ($1, $2)`,
			deviceID, favourite.AnimeID,
		)
	}
	if err != nil {
		return err
	}

	err = r.dbClickhouse.Exec(
		context.Background(),
		`
		INSERT INTO favourite_analytics (anime_id, favourites)
		ALUES (?, 1)
		`,
		favourite.AnimeID,
	)
	if err != nil {
		log.Println("ClickHouse increment failed:", err)
	}

	return nil
}

func (r *FavouriteRepo) RemoveFavourite(deviceID string, favourite model.Favourite) error {
	_, err := r.dbPostgres.Exec(
		"DELETE FROM favourites WHERE device_id = $1 AND anime_id = $2",
		deviceID, favourite.AnimeID,
	)
	if err != nil {
		return err
	}

	err = r.dbClickhouse.Exec(
		context.Background(),
		`
		INSERT INTO favourite_analytics (anime_id, favourites)
		VALUES (?, -1)
		`,
		favourite.AnimeID,
	)
	if err != nil {
		log.Println("ClickHouse decrement failed:", err)
	}

	return nil
}

func (r *FavouriteRepo) GetAllFavourites(deviceID string) ([]model.Favourite, error) {
	rows, err := r.dbPostgres.Query(
		"SELECT device_id, anime_id FROM favourites WHERE device_id = $1 ORDER BY id DESC",
		deviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	favourites := make([]model.Favourite, 0)
	for rows.Next() {
		var f model.Favourite
		if err := rows.Scan(&f.AnimeID); err != nil {
			return nil, err
		}
		favourites = append(favourites, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return favourites, nil
}

func (r *FavouriteRepo) GetFavourites(deviceID string, page, limit int) (model.PaginatedFavourites, error) {
	offset := (page - 1) * limit

	var total int
	err := r.dbPostgres.QueryRow(
		"SELECT COUNT(*) FROM favourites WHERE device_id = $1",
		deviceID,
	).Scan(&total)
	if err != nil {
		return model.PaginatedFavourites{}, err
	}

	rows, err := r.dbPostgres.Query(
		"SELECT anime_id FROM favourites WHERE device_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3",
		deviceID, limit, offset,
	)
	if err != nil {
		return model.PaginatedFavourites{}, err
	}
	defer rows.Close()

	favourites := make([]model.Favourite, 0)
	for rows.Next() {
		var f model.Favourite
		if err := rows.Scan(&f.AnimeID); err != nil {
			return model.PaginatedFavourites{}, err
		}
		favourites = append(favourites, f)
	}

	if err = rows.Err(); err != nil {
		return model.PaginatedFavourites{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	pagesLeft := int(math.Max(0, float64(totalPages-page)))

	return model.PaginatedFavourites{
		Data: favourites,
		Meta: model.PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			PagesLeft:  pagesLeft,
		},
	}, nil
}

func (r *FavouriteRepo) GetFavouriteForAnime(deviceID, animeID string) (model.Favourite, error) {
	var exists bool
	err := r.dbPostgres.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM favourites WHERE device_id=$1 AND anime_id=$2)`,
		deviceID, animeID,
	).Scan(&exists)
	if err != nil {
		return model.Favourite{}, err
	}

	if !exists {
		return model.Favourite{}, nil
	}
	var favourite model.Favourite
	err = r.dbPostgres.QueryRow(
		`SELECT anime_id FROM favourites WHERE device_id=$1 AND anime_id=$2`,
		deviceID, animeID,
	).Scan(&favourite.AnimeID)
	if err != nil {
		return model.Favourite{}, err
	}

	return favourite, nil
}
