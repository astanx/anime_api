package repository

import (
	"database/sql"
	"fmt"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/astanx/anime_api/internal/db"
	"github.com/astanx/anime_api/internal/model"
)

type AnimeRepo struct {
	dbPostgres   *sql.DB
	dbClickhouse clickhouse.Conn
}

func NewAnimeRepo(db *db.DB) *AnimeRepo {
	return &AnimeRepo{
		dbPostgres:   db.Postgres,
		dbClickhouse: db.ClickHouse,
	}
}

func (r *AnimeRepo) GetAnimeByID(id string) (model.SearchAnime, error) {
	row := r.dbPostgres.QueryRow("SELECT id, title, year, poster, type, parser_type FROM search WHERE id = $1", id)
	var anime model.SearchAnime
	err := row.Scan(&anime.ID, &anime.Title, &anime.Year, &anime.Poster, &anime.Type, &anime.ParserType)
	if err != nil {
		return model.SearchAnime{}, err
	}

	return anime, nil
}

func (r *AnimeRepo) SearchConsumetAnime(query string) ([]model.SearchAnime, error) {
	rawResult, err := fetchConsumet(query)
	if err != nil {
		return nil, err
	}

	var result []model.SearchAnime
	for _, a := range rawResult {
		anime := model.SearchAnime{
			ID:         a.ID,
			Title:      a.Title,
			Poster:     a.Poster,
			Year:       a.Year,
			Type:       a.Type,
			ParserType: "Consumet",
		}
		result = append(result, anime)
		if !checkExists(r.dbPostgres, anime.ID) {
			insertAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, query, "consumet", len(result))
	return result, nil
}

func (r *AnimeRepo) SearchAnilibriaAnime(query string) ([]model.SearchAnime, error) {
	result, err := fetchAnilibriaReleases("releases", query, 0)
	if err != nil {
		return nil, err
	}

	for _, anime := range result {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, query, "anilibria", len(result))
	return result, nil
}

func (r *AnimeRepo) SearchAnilibriaRecommendedAnime(limit int) ([]model.SearchAnime, error) {
	result, err := fetchAnilibriaReleases("releases/recommended", "", limit)
	if err != nil {
		return nil, err
	}

	for _, anime := range result {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, "recommended", "anilibria", len(result))
	return result, nil
}

func (r *AnimeRepo) SearchConsumetRecommendedAnime(limit int) ([]model.SearchAnime, error) {
	result, err := fetchConsumetReleases("most-popular")
	if err != nil {
		return nil, err
	}

	for _, anime := range result {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, "recommended", "consumet", len(result))
	return result, nil
}

func (r *AnimeRepo) SearchConsumetLatestReleases(limit int) ([]model.SearchAnime, error) {
	result, err := fetchConsumetReleases("top-airing")
	if err != nil {
		return nil, err
	}

	for _, anime := range result {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, "latest", "consumet", len(result))
	return result, nil
}

func (r *AnimeRepo) SearchAnilibriaLatestReleases(limit int) ([]model.SearchAnime, error) {
	result, err := fetchAnilibriaReleases("releases/latest", "", limit)
	if err != nil {
		return nil, err
	}

	for _, anime := range result {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, "latest", "anilibria", len(result))
	return result, nil
}

func (r *AnimeRepo) SearchAnilibriaRandomReleases(limit int) ([]model.SearchAnime, error) {
	result, err := fetchAnilibriaReleases("releases/random", "", limit)
	if err != nil {
		return nil, err
	}

	for _, anime := range result {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, "random", "anilibria", len(result))
	return result, nil
}

func (r *AnimeRepo) GetAnilibriaGenres() ([]model.Genre, error) {
	baseURL := "https://aniliberty.top/api/v1/anime/genres/%s?include=id,name,total_releases"

	var result []model.Genre
	if err := doJSONRequest(baseURL, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AnimeRepo) GetConsumetGenres() ([]model.ConsumetGenre, error) {
	baseURL := "https://consumet-caou.onrender.com/anime/zoro/genre/list"

	var rawResult []string
	if err := doJSONRequest(baseURL, &rawResult); err != nil {
		return nil, err
	}

	result := make([]model.ConsumetGenre, 0, len(rawResult))
	for _, name := range rawResult {
		result = append(result, model.ConsumetGenre{
			Name: name,
		})
	}

	return result, nil
}

func (r *AnimeRepo) SearchAnilibriaGenreReleases(genreID, limit int) ([]model.SearchAnime, error) {
	result, err := fetchAnilibriaReleases(fmt.Sprintf("releases/genre/%d/releases", genreID), "", limit)
	if err != nil {
		return nil, err
	}

	for _, anime := range result {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, "random", "anilibria", len(result))
	return result, nil
}

func (r *AnimeRepo) SearchConsumetGenreReleases(genre string) ([]model.SearchAnime, error) {
	result, err := fetchConsumetReleases(genre)
	if err != nil {
		return nil, err
	}

	for _, anime := range result {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, fmt.Sprintf("genre-%s", genre), "consumet", len(result))
	return result, nil
}
