package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/astanx/anime_api/internal/config"
	"github.com/astanx/anime_api/internal/db"
	"github.com/astanx/anime_api/internal/model"
	"github.com/redis/go-redis/v9"
)

var bracketRegex = regexp.MustCompile(`\[[^\]]*\]`)

func removeBrackets(s string) string {
	return bracketRegex.ReplaceAllString(s, "")
}

func tokenize(s string) []string {
	s = strings.ToLower(s)
	s = removeBrackets(s)

	// replace separators with space
	replacer := strings.NewReplacer(".", " ", "-", " ", "_", " ")
	s = replacer.Replace(s)

	parts := strings.Fields(s)

	// remove short garbage tokens
	var result []string
	for _, p := range parts {
		if len(p) >= 3 {
			result = append(result, p)
		}
	}

	return result
}

func titleMatchScore(base, candidate string) int {
	baseTokens := tokenize(base)
	candidateTokens := tokenize(candidate)

	match := 0
	for _, b := range baseTokens {
		for _, c := range candidateTokens {
			if b == c {
				match++
				break
			}
		}
	}

	// require at least 1–2 strong matches
	if match == 0 {
		return -200 // HARD reject
	}

	return match * 80
}

func episodeMatch(title string, episode int) (bool, bool) {
	title = strings.ToLower(title)

	// 🔥 HARD NORMALIZATION
	replacer := strings.NewReplacer(
		"~", "-",
		"to", "-",
		"–", "-", // en dash
		"—", "-", // em dash
		"_", " ",
		"x", " ", // sometimes people write 01x02
	)
	title = replacer.Replace(title)

	// remove brackets content (group names, quality, etc.)
	title = removeBrackets(title)

	// remove season ranges like S01-S05 (we don't want to match those as episodes)
	seasonRangeRegex := regexp.MustCompile(`s\d{1,2}\s*-\s*s\d{1,2}`)
	title = seasonRangeRegex.ReplaceAllString(title, "")

	// collapse multiple spaces
	title = strings.Join(strings.Fields(title), " ")

	// --------------------------------------------------
	// 1. RANGE MATCH (STRONG) - This is the main fix
	// --------------------------------------------------
	// Improved regex to catch: 01-220, 1-220, 001~220, 01 ~ 220, etc.
	rangeRegex := regexp.MustCompile(`\b(\d{1,4})\s*[-~]\s*(\d{1,4})\b`)
	matches := rangeRegex.FindAllStringSubmatch(title, -1)

	for _, m := range matches {
		start, err1 := strconv.Atoi(m[1])
		end, err2 := strconv.Atoi(m[2])
		if err1 != nil || err2 != nil {
			continue
		}
		if start > end {
			start, end = end, start
		}

		// Ignore year-like ranges (e.g. 2002-2007, 1999-2003)
		if start > 1900 && end > 1900 && end-start > 10 {
			continue
		}

		// If the requested episode falls inside the range → success
		if episode >= start && episode <= end {
			return true, true // matched, isBatch=true
		}
	}

	// --------------------------------------------------
	// 2. SINGLE EPISODE MATCH (STRICT)
	// --------------------------------------------------
	// Support more patterns: 02, E02, -02-, EP02,  02  (with word boundaries)
	epStr := fmt.Sprintf("%02d", episode)     // 02
	epStrNoZero := fmt.Sprintf("%d", episode) // 2 (for single digit)

	singlePatterns := []string{
		fmt.Sprintf(`(^|[^0-9])%s([^0-9]|$)`, epStr),       // 02
		fmt.Sprintf(`(^|[^0-9])%s([^0-9]|$)`, epStrNoZero), // 2
		fmt.Sprintf(`e%s`, epStr),                          // e02
		fmt.Sprintf(`ep%s`, epStr),                         // ep02
		fmt.Sprintf(`episode\s*%s`, epStr),                 // episode 02
	}

	for _, pattern := range singlePatterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(title) {
			return true, false // matched, isBatch=false
		}
	}

	return false, false
}

type TorrentRepo struct {
	dbPostgres     *sql.DB
	dbClickhouse   clickhouse.Conn
	dbRedis        *redis.Client
	collectionRepo CollectionRepo
	historyRepo    HistoryRepo
	timecodeRepo   TimecodeRepo
	animeRepo      AnimeRepo
}

func NewTorrentRepo(db *db.DB, collectionRepo CollectionRepo, historyRepo HistoryRepo, timecodeRepo TimecodeRepo, animeRepo AnimeRepo) *TorrentRepo {
	return &TorrentRepo{
		dbPostgres:     db.Postgres,
		dbClickhouse:   db.ClickHouse,
		dbRedis:        db.Redis,
		collectionRepo: collectionRepo,
		historyRepo:    historyRepo,
		timecodeRepo:   timecodeRepo,
		animeRepo:      animeRepo,
	}
}

func (r *TorrentRepo) SearchMALAnime(query string, page int) (model.PaginatedSearchAnime, error) {
	url := fmt.Sprintf("https://api.jikan.moe/v4/anime?q=%s&limit=20&page=%d", query, page)

	var res model.PaginatedMALSearchAnime
	if err := doJSONRequest(url, &res); err != nil {
		return model.PaginatedSearchAnime{}, err
	}

	var data []model.SearchAnime
	for _, item := range res.Data {
		data = append(data, model.SearchAnime{
			ID:         fmt.Sprint(item.ID),
			Title:      item.Title,
			Poster:     item.Images.Webp.ImageURL,
			Year:       item.Year,
			Type:       item.Type,
			ParserType: "MAL",
		})
	}
	return model.PaginatedSearchAnime{
		Data: data,
		Meta: model.ShortPaginationMeta{
			CurrentPage: res.Pagination.CurrentPage,
			HasNextPage: res.Pagination.HasNextPage,
			TotalPages:  res.Pagination.LastVisiblePage,
		},
	}, nil
}

func (r *TorrentRepo) SearchMALRecommendedAnime(limit int, page int) ([]model.SearchAnime, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime:search:mal:recommended:limit:%d:page:%d", limit, page)

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var anime []model.SearchAnime
		if err := json.Unmarshal([]byte(cached), &anime); err == nil {
			return anime, nil
		}
	}

	url := fmt.Sprintf("https://api.jikan.moe/v4/top/anime?limit=%d&page=%d", limit, page)
	fmt.Print(url)
	var res model.PaginatedMALSearchAnime
	if err := doJSONRequest(url, &res); err != nil {
		return []model.SearchAnime{}, err
	}

	var data []model.SearchAnime
	for _, item := range res.Data {
		data = append(data, model.SearchAnime{
			ID:         fmt.Sprint(item.ID),
			Title:      item.Title,
			Poster:     item.Images.Webp.ImageURL,
			Year:       item.Year,
			Type:       item.Type,
			ParserType: "MAL",
		})
	}

	logSearchClickhouse(r.dbClickhouse, "recommended", "mal", len(data))

	animeJSON, _ := json.Marshal(data)
	r.dbRedis.Set(ctx, cacheKey, animeJSON, 12*time.Hour)

	return data, nil
}

func (r *TorrentRepo) SearchMALLatestReleases(limit int, page int) ([]model.SearchAnime, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime:search:mal:latest:limit:%d:page:%d", limit, page)

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var anime []model.SearchAnime
		if err := json.Unmarshal([]byte(cached), &anime); err == nil {
			return anime, nil
		}
	}

	url := fmt.Sprintf("https://api.jikan.moe/v4/watch/episodes?limit=%d&page=%d", limit, page)

	var res model.PaginatedMALLatestAnime
	if err := doJSONRequest(url, &res); err != nil {
		return []model.SearchAnime{}, err
	}

	var data []model.SearchAnime
	for _, item := range res.Data {
		data = append(data, model.SearchAnime{
			ID:         fmt.Sprint(item.Entry.ID),
			Title:      item.Entry.Title,
			Poster:     item.Entry.Images.Webp.ImageURL,
			Year:       0,
			Type:       "",
			ParserType: "MAL",
		})
	}

	logSearchClickhouse(r.dbClickhouse, "latest", "mal", len(data))

	animeJSON, _ := json.Marshal(data)
	r.dbRedis.Set(ctx, cacheKey, animeJSON, 12*time.Hour)

	return data, nil
}

func (r *TorrentRepo) SearchMALById(id string) (model.Anime, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime:search:mal_id:%s", id)

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var anime model.Anime
		if err := json.Unmarshal([]byte(cached), &anime); err == nil {
			return anime, nil
		}
	}

	url := fmt.Sprintf("https://api.jikan.moe/v4/anime/%s/full", id)

	var res model.MALAnime
	if err := doJSONRequest(url, &res); err != nil {
		return model.Anime{}, err
	}

	var genres []string
	for _, g := range res.Genres {
		genres = append(genres, g.Name)
	}

	result := model.Anime{
		ID:            fmt.Sprint(res.ID),
		MalID:         res.ID,
		Title:         res.Title,
		Poster:        res.Images.Webp.ImageURL,
		Description:   res.Description,
		Genres:        genres,
		Status:        res.Status,
		Year:          res.Year,
		Type:          res.Type,
		TotalEpisodes: res.TotalEpisodes,
		Episodes:      make([]model.PreviewEpisode, 0),
	}

	hasNext := true
	page := 1

	for {
		if !hasNext {
			break
		}
		url := fmt.Sprintf("https://api.jikan.moe/v4/anime/%s/episodes?page=%d", id, page)

		var episodesRes model.PaginatedMALPreviewEpisodes
		if err := doJSONRequest(url, &episodesRes); err != nil {
			return model.Anime{}, err
		}

		for _, e := range episodesRes.Data {
			last := path.Base(e.Url)
			ordinal, _ := strconv.Atoi(last)
			result.Episodes = append(result.Episodes, model.PreviewEpisode{
				ID:       fmt.Sprint(e.ID),
				Title:    e.Title,
				Ordinal:  ordinal,
				IsSubbed: true,
			})
		}

		hasNext = episodesRes.Pagination.HasNextPage
		page++
	}

	animeJSON, _ := json.Marshal(result)
	r.dbRedis.Set(ctx, cacheKey, animeJSON, 12*time.Hour)

	return result, nil
}

func (r *TorrentRepo) SearchMALByEpisodeId(id, episodeId string) (model.Episode, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("anime:mal:episode:mal_id:%s:episode_id:%s", id, episodeId)

	cached, err := r.dbRedis.Get(ctx, cacheKey).Result()
	if err == nil {
		var episode model.Episode
		if err := json.Unmarshal([]byte(cached), &episode); err == nil {
			return episode, nil
		}
	}

	url := fmt.Sprintf("https://api.jikan.moe/v4/anime/%s/episodes/%s", id, episodeId)

	var res model.MalPreviewEpisode
	if err := doJSONRequest(url, &res); err != nil {
		return model.Episode{}, err
	}

	result := model.Episode{
		ID:        fmt.Sprintf("%s/%d", id, res.Data.ID),
		Title:     res.Data.Title,
		Ordinal:   res.Data.ID,
		Sources:   []model.Source{},
		Subtitles: []model.Subtitle{},
	}

	animeUrl := strings.TrimPrefix(res.Data.Url, "https://")
	animeUrl = strings.TrimPrefix(animeUrl, "http://")

	parts := strings.Split(animeUrl, "/")

	var title string

	if len(parts) >= 4 {
		title = parts[3]
	}
	if title == "" {
		title = res.Data.Title
	}
	episodeNum := res.Data.ID

	queries := []string{
		fmt.Sprintf("%s %02d", title, episodeNum),
		fmt.Sprintf("%s E%02d", title, episodeNum),
		fmt.Sprintf("%s - %02d", title, episodeNum),
	}

	var prowres []model.ProwlarrAnime
	for _, q := range queries {
		searchURL := fmt.Sprintf("%s/search?apikey=%s&query=%s&categories=%s&type=search",
			config.PROWLARR_URL, config.PROWLARR_APIKEY, q, config.CATEGORY)

		var temp []model.ProwlarrAnime
		if err := doJSONRequest(searchURL, &temp); err == nil {
			prowres = append(prowres, temp...)
		}
	}

	type scoredResult struct {
		item  model.ProwlarrAnime
		score int
	}
	var scored []scoredResult

	for _, p := range prowres {
		cleanTitle := removeBrackets(p.Title)
		cleanSort := removeBrackets(p.SortTitle)

		score := 0

		// Seeder weight
		if p.Seeders > 0 {
			score += p.Seeders * 10
		}
		score += p.Grabs * 2

		// Title match (outside brackets)
		titleScore := titleMatchScore(title, p.Title)
		score += titleScore

		if titleScore < 0 {
			continue
		}

		matched1, batch1 := episodeMatch(cleanTitle, res.Data.ID)
		matched2, batch2 := episodeMatch(cleanSort, res.Data.ID)

		if !matched1 && !matched2 {
			continue
		}

		score += 200

		// Penalize batch torrents
		if batch1 || batch2 {
			score -= 80
		}

		if score > 100 {
			scored = append(scored, scoredResult{item: p, score: score})
		}
	}

	// Sort by score descending
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Take the best one (or top few if you want multiple sources)
	for _, s := range scored {
		p := s.item

		result.Sources = append(result.Sources, model.Source{
			Url:  fmt.Sprintf("%s/stream?link=%s&index=1&play", p.Title, p.Hash),
			Type: "TORRENT",
		})

		// // Add torrent
		// addURL := fmt.Sprintf("%s/torrents", config.TORR_URL)
		// payload := fmt.Sprintf(`{"action":"add","link":"%s","save_to_db":true}`, p.MagnetURL)
		// http.Post(addURL, "application/json", strings.NewReader(payload))

		// // Better index heuristic:
		// // Many single-episode torrents have index=1 (or only one file).
		// // For batches, we still often use index=1, but at least the torrent is now more likely correct.
		// // Future improvement: parse file list from torrent client API if available.
		// videoURL := fmt.Sprintf("%s/stream?link=%s&index=1&play", config.TORR_URL, p.Hash)

		// result.Sources = append(result.Sources, model.Source{
		// 	Url:  videoURL,
		// 	Type: "TORRENT",
		// })

		// break // take the best one for now
	}

	// animeJSON, _ := json.Marshal(result)
	// r.dbRedis.Set(ctx, cacheKey, animeJSON, 1*time.Hour)

	return result, nil
}
