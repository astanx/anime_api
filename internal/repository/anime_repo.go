package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/astanx/anime_api/internal/config"
	"github.com/astanx/anime_api/internal/db"
	"github.com/astanx/anime_api/internal/model"
	"github.com/redis/go-redis/v9"
)

type AnimeRepo struct {
	dbPostgres   *sql.DB
	dbClickhouse clickhouse.Conn
	dbRedis      *redis.Client
}

func NewAnimeRepo(db *db.DB) *AnimeRepo {
	return &AnimeRepo{
		dbPostgres:   db.Postgres,
		dbClickhouse: db.ClickHouse,
		dbRedis:      db.Redis,
	}
}

func (r *AnimeRepo) SearchAnimeByID(id string) (model.SearchAnime, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime:search:id:%s", id)

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var anime model.SearchAnime
		if err := json.Unmarshal([]byte(cached), &anime); err == nil {
			return anime, nil
		}
	}
	row := r.dbPostgres.QueryRow("SELECT id, title, year, poster, type, parser_type FROM search WHERE id = $1", id)
	var anime model.SearchAnime
	err = row.Scan(&anime.ID, &anime.Title, &anime.Year, &anime.Poster, &anime.Type, &anime.ParserType)
	if err != nil {
		return model.SearchAnime{}, err
	}
	animeJSON, _ := json.Marshal(anime)
	r.dbRedis.Set(ctx, cacheKey, animeJSON, 1*time.Hour)

	return anime, nil
}

func (r *AnimeRepo) GetAnimeInfoByConsumetID(id string) (model.Anime, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime:consumet:id:%s", id)

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var anime model.Anime
		if err := json.Unmarshal([]byte(cached), &anime); err == nil {
			return anime, nil
		}
	}
	url := fmt.Sprintf("%s/anime/zoro/info?id=%s", config.ConsumetUrl, id)
	var result model.ConsumetAnime
	if err := doJSONRequest(url, &result); err != nil {
		return model.Anime{}, err
	}
	var anime = model.Anime{
		ID:            result.ID,
		Title:         result.Title,
		Poster:        result.Image,
		Description:   result.Description,
		Genres:        result.Genres,
		Status:        result.Status,
		Type:          result.Type,
		TotalEpisodes: result.TotalEpisodes,
		Episodes: func() []model.PreviewEpisode {
			episodes := make([]model.PreviewEpisode, len(result.Episodes))
			for i, e := range result.Episodes {
				episodes[i] = model.PreviewEpisode{
					ID:       e.ID,
					IsDubbed: e.IsDubbed,
					IsSubbed: e.IsSubbed,
					Ordinal:  e.Number,
					Title:    e.Title,
				}
			}
			return episodes
		}(),
	}
	animeJSON, _ := json.Marshal(anime)
	r.dbRedis.Set(ctx, cacheKey, animeJSON, 1*time.Hour)

	return anime, nil
}

func (r *AnimeRepo) GetAnimeInfoByAnilibriaID(id string) (model.Anime, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime:anilibria:id:%s", id)

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var anime model.Anime
		if err := json.Unmarshal([]byte(cached), &anime); err == nil {
			return anime, nil
		}
	}
	url := fmt.Sprintf("https://aniliberty.top/api/v1/anime/releases/%s?include=id,type.value,year,name.main,poster.src,is_ongoing,description,episodes_total,genres.name,episodes", id)
	var result model.AnilibriaAnime
	if err := doJSONRequest(url, &result); err != nil {
		return model.Anime{}, err
	}

	var anime model.Anime
	anime.ID = fmt.Sprintf("%d", result.ID)
	anime.Title = result.Name.Main
	anime.Year = result.Year
	anime.Poster =
		fmt.Sprintf("https://aniliberty.top%s", result.Poster.Src)
	anime.Type = result.Type.Value
	anime.Status = "Completed"
	if result.IsOngoing {
		anime.Status = "Ongoing"
	}
	anime.Description = result.Description
	anime.TotalEpisodes = result.EpisodesTotal

	genres := make([]string, len(result.Genres))
	for i, g := range result.Genres {
		genres[i] = g.Name
	}
	anime.Genres = genres

	episodes := make([]model.PreviewEpisode, len(result.Episodes))
	for i, e := range result.Episodes {
		episodes[i] = model.PreviewEpisode{
			ID:       e.ID,
			IsDubbed: true,
			IsSubbed: false,
			Title:    e.Name,
			Ordinal:  e.Ordinal,
		}
	}
	anime.Episodes = episodes

	animeJSON, _ := json.Marshal(anime)
	r.dbRedis.Set(ctx, cacheKey, animeJSON, 1*time.Hour)

	return anime, nil
}

func (r *AnimeRepo) GetAnilibriaEpisodeInfo(id string) (model.Episode, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime:anilibria:episode:id:%s", id)

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var episode model.Episode
		if err := json.Unmarshal([]byte(cached), &episode); err == nil {
			return episode, nil
		}
	}
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

	episodeJSON, _ := json.Marshal(episode)
	r.dbRedis.Set(ctx, cacheKey, episodeJSON, 1*time.Hour)

	return episode, nil
}

func (r *AnimeRepo) GetConsumetEpisodeInfo(id, title string, ordinal int, dub string) (model.Episode, error) {
	url := fmt.Sprintf("%s/anime/zoro/watch?episodeId=%s&dub=%s", config.ConsumetUrl, id, dub)
	var result model.ConsumetEpisode
	if err := doJSONRequest(url, &result); err != nil {
		return model.Episode{}, err
	}

	finalTitle := result.Title
	if title != "" {
		finalTitle = title
	}

	finalOrdinal := result.Ordinal
	if ordinal != -1 {
		finalOrdinal = ordinal
	}

	episode := model.Episode{
		ID:      id,
		Title:   finalTitle,
		Ordinal: finalOrdinal,
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

	return episode, nil
}

func (r *AnimeRepo) SearchConsumetAnime(query string, page int) (model.PaginatedSearchAnime, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime:search:consumet:query:%s:page:%d", query, page)

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var anime model.PaginatedSearchAnime
		if err := json.Unmarshal([]byte(cached), &anime); err == nil {
			return anime, nil
		}
	}
	rawResult, err := fetchConsumet(query, page)
	if err != nil {
		return model.PaginatedSearchAnime{}, err
	}

	result := model.PaginatedSearchAnime{
		Data: make([]model.SearchAnime, 0),
		Meta: rawResult.Meta,
	}
	for _, a := range rawResult.Data {
		anime := model.SearchAnime{
			ID:         a.ID,
			Title:      a.Title,
			Poster:     a.Poster,
			Year:       a.Year,
			Type:       a.Type,
			ParserType: "Consumet",
		}
		result.Data = append(result.Data, anime)
		if !checkExists(r.dbPostgres, anime.ID) {
			insertSearchAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, query, "consumet", len(result.Data))

	animeJSON, _ := json.Marshal(result)
	r.dbRedis.Set(ctx, cacheKey, animeJSON, 8*time.Hour)

	return result, nil
}

func (r *AnimeRepo) SearchAnilibriaAnime(query string, page int) (model.PaginatedSearchAnime, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime:search:anilibria:query:%s:page:%d", query, page)

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var anime model.PaginatedSearchAnime
		if err := json.Unmarshal([]byte(cached), &anime); err == nil {
			return anime, nil
		}
	}
	result, err := fetchAnilibriaReleases("app/search/releases", query, 0, page)
	if err != nil {
		return model.PaginatedSearchAnime{}, err
	}

	for _, anime := range result.Data {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertSearchAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, query, "anilibria", len(result.Data))

	animeJSON, _ := json.Marshal(result)
	r.dbRedis.Set(ctx, cacheKey, animeJSON, 8*time.Hour)

	return result, nil
}

func (r *AnimeRepo) SearchAnilibriaRecommendedAnime(limit int, page int) ([]model.SearchAnime, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime:search:anilibria:recommended:limit:%d:page:%d", limit, page)

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var anime []model.SearchAnime
		if err := json.Unmarshal([]byte(cached), &anime); err == nil {
			return anime, nil
		}
	}
	result, err := fetchAnilibriaReleases("anime/releases/recommended", "", limit, page)
	if err != nil {
		return nil, err
	}

	for _, anime := range result.Data {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertSearchAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, "recommended", "anilibria", len(result.Data))

	animeJSON, _ := json.Marshal(result)
	r.dbRedis.Set(ctx, cacheKey, animeJSON, 8*time.Hour)

	return result.Data, nil
}

func (r *AnimeRepo) SearchConsumetRecommendedAnime() ([]model.SearchAnime, error) {
	ctx := context.Background()
	cacheKey := "anime:search:consumet:recommended"

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var anime []model.SearchAnime
		if err := json.Unmarshal([]byte(cached), &anime); err == nil {
			return anime, nil
		}
	}
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

	animeJSON, _ := json.Marshal(result)
	r.dbRedis.Set(ctx, cacheKey, animeJSON, 8*time.Hour)

	return result, nil
}

func (r *AnimeRepo) SearchConsumetLatestReleases() ([]model.SearchAnime, error) {
	ctx := context.Background()
	cacheKey := "anime:search:consumet:latest"

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var anime []model.SearchAnime
		if err := json.Unmarshal([]byte(cached), &anime); err == nil {
			return anime, nil
		}
	}

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

	animeJSON, _ := json.Marshal(result)
	r.dbRedis.Set(ctx, cacheKey, animeJSON, 8*time.Hour)

	return result, nil
}

func (r *AnimeRepo) SearchAnilibriaLatestReleases(limit int) ([]model.SearchAnime, error) {
	ctx := context.Background()
	cacheKey := "anime:search:anilibria:latest"

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var anime []model.SearchAnime
		if err := json.Unmarshal([]byte(cached), &anime); err == nil {
			return anime, nil
		}
	}
	result, err := fetchAnilibriaReleases("anime/releases/latest", "", limit, 0)
	if err != nil {
		return nil, err
	}

	for _, anime := range result.Data {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertSearchAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, "latest", "anilibria", len(result.Data))

	animeJSON, _ := json.Marshal(result)
	r.dbRedis.Set(ctx, cacheKey, animeJSON, 8*time.Hour)
	return result.Data, nil
}

func (r *AnimeRepo) SearchAnilibriaRandomReleases(limit int, page int) (model.PaginatedSearchAnime, error) {
	result, err := fetchAnilibriaReleases("anime/releases/random", "", limit, page)
	if err != nil {
		return model.PaginatedSearchAnime{}, err
	}

	for _, anime := range result.Data {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertSearchAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, "random", "anilibria", len(result.Data))
	return result, nil
}

func (r *AnimeRepo) GetAnilibriaGenres() ([]model.Genre, error) {
	ctx := context.Background()
	cacheKey := "anime:anilibria:genres"

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var genres []model.Genre
		if err := json.Unmarshal([]byte(cached), &genres); err == nil {
			return genres, nil
		}
	}
	baseURL := "https://aniliberty.top/api/v1/anime/genres?include=id,name,total_releases"

	var result []model.Genre
	if err := doJSONRequest(baseURL, &result); err != nil {
		return nil, err
	}

	genreJSON, _ := json.Marshal(result)
	r.dbRedis.Set(ctx, cacheKey, genreJSON, 24*time.Hour)

	return result, nil
}

func (r *AnimeRepo) GetConsumetGenres() ([]string, error) {
	ctx := context.Background()
	cacheKey := "anime:consumet:genres"

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var genres []string
		if err := json.Unmarshal([]byte(cached), &genres); err == nil {
			return genres, nil
		}
	}
	baseURL := fmt.Sprintf("%s/anime/zoro/genre/list", config.ConsumetUrl)

	var result []string
	if err := doJSONRequest(baseURL, &result); err != nil {
		return nil, err
	}
	genreJSON, _ := json.Marshal(result)
	r.dbRedis.Set(ctx, cacheKey, genreJSON, 24*time.Hour)

	return result, nil
}

func (r *AnimeRepo) SearchAnilibriaGenreReleases(genreID, limit int, page int) (model.PaginatedSearchAnime, error) {
	result, err := fetchAnilibriaReleases(fmt.Sprintf("anime/releases/genre/%d/releases", genreID), "", limit, page)
	if err != nil {
		return model.PaginatedSearchAnime{}, err
	}

	for _, anime := range result.Data {
		if !checkExists(r.dbPostgres, anime.ID) {
			insertSearchAnime(r.dbPostgres, anime)
		}
	}

	logSearchClickhouse(r.dbClickhouse, "random", "anilibria", len(result.Data))
	return result, nil
}

func (r *AnimeRepo) SearchConsumetGenreReleases(genre string) ([]model.SearchAnime, error) {
	result, err := fetchConsumetReleases(fmt.Sprintf("genre/%s", genre))
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
