package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/astanx/anime_api/internal/model"
)

// --- HTTP helper ---

func doJSONRequest(url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed [%d]: %s\n", resp.StatusCode, body)
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	return json.Unmarshal(body, target)
}

func fetchConsumet(endpoint string) ([]model.SearchAnime, error) {
	url := fmt.Sprintf("https://consumet-caou.onrender.com/anime/zoro/%s", endpoint)
	var res struct {
		Results []model.SearchAnime `json:"results"`
	}
	if err := doJSONRequest(url, &res); err != nil {
		return nil, err
	}
	return res.Results, nil
}

func fetchAnilibriaReleases(endpoint string, query string, limit int) ([]model.SearchAnime, error) {
	baseURL := fmt.Sprintf("https://aniliberty.top/api/v1/%s?include=id,type.value,year,poster.optimized.thumbnail,name.main", endpoint)

	if query != "" {
		baseURL += fmt.Sprintf("&query=%s", query)
	}
	if limit > 0 {
		baseURL += fmt.Sprintf("&limit=%d", limit)
	}

	var rawResult []model.SearchAnilibriaAnime
	if err := doJSONRequest(baseURL, &rawResult); err != nil {
		return nil, err
	}

	result := make([]model.SearchAnime, 0, len(rawResult))
	for _, a := range rawResult {
		result = append(result, model.SearchAnime{
			ID:         fmt.Sprint(a.ID),
			Title:      a.Name.Main,
			Poster:     fmt.Sprintf("https://aniliberty.top%s", a.Poster.Optimized.Thumbnail),
			Year:       a.Year,
			Type:       a.Type.Value,
			ParserType: "Anilibria",
		})
	}

	return result, nil
}

func fetchConsumetReleases(endpoint string) ([]model.SearchAnime, error) {
	baseURL := fmt.Sprintf("https://consumet-caou.onrender.com/anime/zoro/%s", endpoint)

	var rawResult []model.SearchAnime
	if err := doJSONRequest(baseURL, &rawResult); err != nil {
		return nil, err
	}

	result := make([]model.SearchAnime, 0, len(rawResult))
	for _, a := range rawResult {
		result = append(result, model.SearchAnime{
			ID:         a.ID,
			Title:      a.Title,
			Poster:     a.Poster,
			Year:       a.Year,
			Type:       a.Type,
			ParserType: "Consumet",
		})
	}

	return result, nil
}

// --- PostgreSQL helpers ---

func checkExists(db *sql.DB, id string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM search WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		log.Printf("Error checking anime existence: %v", err)
		return true
	}
	return exists
}

func insertSearchAnime(db *sql.DB, anime model.SearchAnime) {
	_, err := db.Exec(
		"INSERT INTO search (id, title, year, poster, type, parser_type) VALUES ($1, $2, $3, $4, $5, $6)",
		anime.ID, anime.Title, anime.Year, anime.Poster, anime.Type, anime.ParserType,
	)
	if err != nil {
		log.Printf("Error inserting anime: %v", err)
	}
}

func insertEpisode(db *sql.DB, episode model.Episode) {
	_, e := db.Exec("INSERT INTO episodes (id, anime_id ordinal, title, opening_start, opening_end, ending_start, ending_end) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		episode.ID, episode.Ordinal, episode.Title,
		episode.Opening.Start, episode.Opening.End,
		episode.Ending.Start, episode.Ending.End)
	if e != nil {
		log.Printf("Error inserting anime: %v", e)
	}
	for _, source := range episode.Sources {
		_, e = db.Exec("INSERT INTO episode_sources (episode_id, url, type) VALUES ($1, $2, $3)", episode.ID, source.Url, source.Type)
		if e != nil {
			log.Printf("Error inserting source: %v", e)
		}
	}

	for _, subtitle := range episode.Subtitles {
		_, e = db.Exec("INSERT INTO episode_subtitles (episode_id, vtt, language) VALUES ($1, $2, $3)", episode.ID, subtitle.Vtt, subtitle.Language)
		if e != nil {
			log.Printf("Error inserting subtitle: %v", e)
		}
	}
}

func getEpisode(db *sql.DB, id string) (model.Episode, bool, error) {
	var episode model.Episode
	row := db.QueryRow(`
		SELECT id, ordinal, title,
		       opening_start, opening_end,
		       ending_start, ending_end
		FROM episodes
		WHERE id = $1
	`, id)

	var openingStart, openingEnd, endingStart, endingEnd sql.NullInt64

	err := row.Scan(
		&episode.ID,
		&episode.Ordinal,
		&episode.Title,
		&openingStart, &openingEnd,
		&endingStart, &endingEnd,
	)
	if err == nil {
		return model.Episode{}, false, nil
	}

	episode.Opening = model.TimeSegment{
		Start: int(openingStart.Int64),
		End:   int(openingEnd.Int64),
	}
	episode.Ending = model.TimeSegment{
		Start: int(endingStart.Int64),
		End:   int(endingEnd.Int64),
	}

	rows, err := db.Query(`SELECT url, type FROM episode_sources WHERE episode_id = $1`, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var s model.Source
			if err := rows.Scan(&s.Url, &s.Type); err == nil {
				episode.Sources = append(episode.Sources, s)
			}
		}
	}

	srows, err := db.Query(`SELECT vtt, language FROM episode_subtitles WHERE episode_id = $1`, id)
	if err == nil {
		defer rows.Close()
		for srows.Next() {
			var s model.Subtitle
			if err := srows.Scan(&s.Vtt, &s.Language); err == nil {
				episode.Subtitles = append(episode.Subtitles, s)
			}
		}
	}

	return episode, true, nil

}

// --- ClickHouse helper ---

func logSearchClickhouse(conn clickhouse.Conn, query, parserType string, resultCount int) {
	err := conn.Exec(
		context.Background(),
		"INSERT INTO search_analytics (query, type, results, searched_at) VALUES (?, ?, ?, now())",
		query, parserType, resultCount,
	)
	if err != nil {
		log.Println("ClickHouse insert failed:", err)
	}
}
