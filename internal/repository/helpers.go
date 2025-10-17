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

func fetchConsumet(query string) ([]model.SearchAnime, error) {
	url := fmt.Sprintf("https://consumet-caou.onrender.com/anime/zoro/%s", query)
	var res struct {
		Results []model.SearchAnime `json:"results"`
	}
	if err := doJSONRequest(url, &res); err != nil {
		return nil, err
	}
	return res.Results, nil
}

func fetchAnilibriaReleases(endpoint string, query string, limit int) ([]model.SearchAnime, error) {
	baseURL := fmt.Sprintf("https://aniliberty.top/api/v1/app/search/%s?include=id,type.value,year,poster.optimized.thumbnail,name.main", endpoint)

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
		return true // предотвращаем дубли
	}
	return exists
}

func insertAnime(db *sql.DB, anime model.SearchAnime) {
	_, err := db.Exec(
		"INSERT INTO search (id, title, year, poster, type, parser_type) VALUES ($1, $2, $3, $4, $5, $6)",
		anime.ID, anime.Title, anime.Year, anime.Poster, anime.Type, anime.ParserType,
	)
	if err != nil {
		log.Printf("Error inserting anime: %v", err)
	}
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
