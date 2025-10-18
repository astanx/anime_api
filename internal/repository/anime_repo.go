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

func (r *AnimeRepo) GetSearchAnimeByID(id string) (model.SearchAnime, error) {
	row := r.dbPostgres.QueryRow("SELECT id, title, year, poster, type, parser_type FROM search WHERE id = $1", id)
	var anime model.SearchAnime
	err := row.Scan(&anime.ID, &anime.Title, &anime.Year, &anime.Poster, &anime.Type, &anime.ParserType)
	if err != nil {
		return model.SearchAnime{}, err
	}

	return anime, nil
}

func (r *AnimeRepo) GetAnimeInfoByConsumetID(id string) (model.Anime, error) {
	url := fmt.Sprintf("https://consumet-caou.onrender.com/anime/zoro/info?id=%s", id)
	var anime model.Anime
	if err := doJSONRequest(url, &anime); err != nil {
		return model.Anime{}, err
	}
	return anime, nil
}

func (r *AnimeRepo) GetAnimeInfoByAnilibriaID(id string) (model.Anime, error) {
	url := fmt.Sprintf("https://aniliberty.top/api/v1/anime/releases/%s?include=id,type.value,year,name.main,poster.optimized.thumbnail,is_ongoing,description,episodes_total,genres.name,episodes", id)
	var result model.AnilibriaAnime
	if err := doJSONRequest(url, &result); err != nil {
		return model.Anime{}, err
	}

	var anime model.Anime
	anime.ID = fmt.Sprintf("%d", result.ID)
	anime.Title = result.Name.Main
	anime.Year = result.Year
	anime.Poster = result.Poster.Optimized.Thumbnail
	anime.Type = result.Type.Value
	anime.Status = "Completed"
	if result.IsOngoing {
		anime.Status = "Ongoing"
	}
	anime.Description = result.Description
	anime.TotalEpisodes = result.EpisodesTotal

	genres := make([]model.ConsumetGenre, len(result.Genres))
	for i, g := range result.Genres {
		genres[i] = model.ConsumetGenre{
			Name: g.Name,
		}
	}
	anime.Genres = genres

	episodes := make([]model.PreviewEpisode, len(result.Episodes))
	for i, e := range result.Episodes {
		episodes[i] = model.PreviewEpisode{
			ID:      e.ID,
			Title:   e.Name,
			Ordinal: e.Ordinal,
		}
	}
	anime.Episodes = episodes

	return anime, nil
}

func (r *AnimeRepo) GetAnilibriaEpisodeInfo(id string) (model.Episode, error) {
	episode, exists, err := getEpisode(r.dbPostgres, id)
	if err != nil {
		fmt.Printf("Error getEpisode from db: %s", err)
	}
	if exists {
		return episode, nil
	}
	url := fmt.Sprintf("https://aniliberty.top/api/v1/anime/releases/episodes/%s?include=id,name,ordinal,opening,ending,hls_480,hls_720,hls_1080", id)
	var result model.AnilibriaEpisode
	if err := doJSONRequest(url, &result); err != nil {
		return model.Episode{}, err
	}

	episode = model.Episode{
		ID:      result.ID,
		Title:   result.Name,
		Ordinal: result.Ordinal,
		Opening: model.TimeSegment{
			Start: result.Opening.Start,
			End:   result.Opening.End,
		},
		Ending: model.TimeSegment{
			Start: result.Ending.Start,
			End:   result.Ending.End,
		},
	}

	var sources []model.Source
	if result.Hls480 != "" {
		sources = append(sources, model.Source{
			Url:  result.Hls480,
			Type: "hls480",
		})
	}
	if result.Hls720 != "" {
		sources = append(sources, model.Source{
			Url:  result.Hls720,
			Type: "hls720",
		})
	}
	if result.Hls1080 != "" {
		sources = append(sources, model.Source{
			Url:  result.Hls1080,
			Type: "hls1080",
		})
	}
	episode.Sources = sources

	insertEpisode(r.dbPostgres, episode)

	return episode, nil
}

func (r *AnimeRepo) GetConsumetEpisodeInfo(id string) (model.Episode, error) {
	episode, exists, err := getEpisode(r.dbPostgres, id)
	if err != nil {
		fmt.Printf("Error getEpisode from db: %s", err)
	}
	if exists {
		return episode, nil
	}
	url := fmt.Sprintf("https://consumet-caou.onrender.com/anime/zoro/watch?episodeId=%s", id)
	var result model.ConsumetEpisode
	if err := doJSONRequest(url, &result); err != nil {
		return model.Episode{}, err
	}

	episode = model.Episode{
		ID:      result.ID,
		Title:   result.Title,
		Ordinal: result.Ordinal,
		Opening: model.TimeSegment{
			Start: result.Intro.Start,
			End:   result.Intro.End,
		},
		Ending: model.TimeSegment{
			Start: result.Outro.Start,
			End:   result.Outro.End,
		},
		Sources:   result.Sources,
		Subtitles: result.Subtitles,
	}

	insertEpisode(r.dbPostgres, episode)

	return episode, nil
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
			insertSearchAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, query, "consumet", len(result))
	return result, nil
}

func (r *AnimeRepo) SearchAnilibriaAnime(query string) ([]model.SearchAnime, error) {
	result, err := fetchAnilibriaReleases("app/search/releases", query, 0)
	if err != nil {
		return nil, err
	}

	for _, anime := range result {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertSearchAnime(r.dbPostgres, anime)
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
			insertSearchAnime(r.dbPostgres, anime)
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
			insertSearchAnime(r.dbPostgres, anime)
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
			insertSearchAnime(r.dbPostgres, anime)
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
			insertSearchAnime(r.dbPostgres, anime)
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
			insertSearchAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, "random", "anilibria", len(result))
	return result, nil
}

func (r *AnimeRepo) GetAnilibriaGenres() ([]model.Genre, error) {
	baseURL := "https://aniliberty.top/api/v1/anime/genres?include=id,name,total_releases"

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
			insertSearchAnime(r.dbPostgres, anime)
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
			insertSearchAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, fmt.Sprintf("genre-%s", genre), "consumet", len(result))
	return result, nil
}
