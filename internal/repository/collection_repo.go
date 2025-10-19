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

type CollectionRepo struct {
	dbPostgres   *sql.DB
	dbClickhouse clickhouse.Conn
}

func NewCollectionRepo(db *db.DB) *CollectionRepo {
	return &CollectionRepo{
		dbPostgres:   db.Postgres,
		dbClickhouse: db.ClickHouse,
	}
}

func (r *CollectionRepo) AddCollection(deviceID string, collection model.Collection) error {
	_, err := r.dbPostgres.Exec(
		`INSERT INTO collections (device_id, anime_id, name)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (device_id, anime_id) DO UPDATE
		 SET name = EXCLUDED.name`,
		deviceID, collection.AnimeID, collection.Type,
	)
	if err != nil {
		return err
	}

	err = r.dbClickhouse.Exec(
		context.Background(),
		`INSERT INTO collection_analytics (anime_id, type, count)
		 VALUES (?, ?, 1)
		 ON CONFLICT (anime_id, type) DO UPDATE
		 SET count = count + 1`,
		collection.AnimeID, collection.Type,
	)
	if err != nil {
		log.Println("ClickHouse increment failed:", err)
	}

	return nil
}

func (r *CollectionRepo) RemoveCollection(deviceID, animeID, collectionType string) error {
	_, err := r.dbPostgres.Exec(
		"DELETE FROM collections WHERE device_id = $1 AND anime_id = $2",
		deviceID, animeID,
	)
	if err != nil {
		return err
	}

	err = r.dbClickhouse.Exec(
		context.Background(),
		`INSERT INTO collection_analytics (anime_id, type, count)
		 VALUES (?, ?, 0)
		 ON CONFLICT (anime_id, type) DO UPDATE
		 SET count = count - 1`,
		animeID, collectionType,
	)
	if err != nil {
		log.Println("ClickHouse decrement failed:", err)
	}

	return nil
}

func (r *CollectionRepo) GetAllCollections(deviceID string) ([]model.Collection, error) {
	rows, err := r.dbPostgres.Query(
		"SELECT device_id, anime_id, name FROM collections WHERE device_id = $1",
		deviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []model.Collection
	for rows.Next() {
		var c model.Collection
		if err := rows.Scan(&c.AnimeID, &c.Type); err != nil {
			return nil, err
		}
		collections = append(collections, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return collections, nil
}

func (r *CollectionRepo) GetCollections(deviceID string, page, limit int) (model.PaginatedCollections, error) {
	offset := (page - 1) * limit

	var total int
	err := r.dbPostgres.QueryRow(
		"SELECT COUNT(*) FROM collections WHERE device_id = $1",
		deviceID,
	).Scan(&total)
	if err != nil {
		return model.PaginatedCollections{}, err
	}

	rows, err := r.dbPostgres.Query(
		`SELECT device_id, anime_id, name
		 FROM collections
		 WHERE device_id = $1
		 ORDER BY anime_id
		 LIMIT $2 OFFSET $3`,
		deviceID, limit, offset,
	)
	if err != nil {
		return model.PaginatedCollections{}, err
	}
	defer rows.Close()

	var collections []model.Collection
	for rows.Next() {
		var c model.Collection
		if err := rows.Scan(&c.AnimeID, &c.Type); err != nil {
			return model.PaginatedCollections{}, err
		}
		collections = append(collections, c)
	}

	if err = rows.Err(); err != nil {
		return model.PaginatedCollections{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	pagesLeft := int(math.Max(0, float64(totalPages-page)))

	return model.PaginatedCollections{
		Data: collections,
		Meta: model.PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			PagesLeft:  pagesLeft,
		},
	}, nil
}
