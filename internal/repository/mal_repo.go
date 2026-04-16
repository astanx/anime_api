package repository

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/astanx/anime_api/internal/config"
	"github.com/astanx/anime_api/internal/db"
	"github.com/astanx/anime_api/internal/model"
)

type MALRepo struct {
	dbPostgres     *sql.DB
	dbClickhouse   clickhouse.Conn
	collectionRepo CollectionRepo
	historyRepo    HistoryRepo
	timecodeRepo   TimecodeRepo
}

func NewMALRepo(db *db.DB, collectionRepo CollectionRepo, historyRepo HistoryRepo, timecodeRepo TimecodeRepo) *MALRepo {
	return &MALRepo{
		dbPostgres:     db.Postgres,
		dbClickhouse:   db.ClickHouse,
		collectionRepo: collectionRepo,
		historyRepo:    historyRepo,
		timecodeRepo:   timecodeRepo,
	}
}

func convertToMALStatus(status string) (string, error) {
	switch status {
	case "watching":
		return "Watching", nil
	case "watched":
		return "Completed", nil
	case "abandoned":
		return "Dropped", nil
	case "planned":
		return "Plan to Watch", nil
	default:
		return "", fmt.Errorf("invalid status: %s", status)
	}
}

func convertMalToStatus(mal string) (string, error) {
	switch mal {
	case "Watching":
		return "watching", nil
	case "Completed":
		return "watched", nil
	case "Dropped":
		return "abandoned", nil
	case "Plan to Watch":
		return "planned", nil
	default:
		return "", fmt.Errorf("invalid mal status: %s", mal)
	}
}

func (r *MALRepo) ImportMALList(deviceID, malList string) (int, error) {

	var mal model.MALList
	if err := xml.Unmarshal([]byte(malList), &mal); err != nil {
		return 0, fmt.Errorf("failed to unmarshal MAL XML: %w", err)
	}

	var count int

	for _, anime := range mal.Animes {
		url := fmt.Sprintf("%s/anime/hianime/%s", config.ConsumetUrl, anime.SeriesTitle)

		var res model.PaginatedConsumetSearchAnime
		if err := doJSONRequest(url, &res); err != nil {
			return 0, err
		}

		for _, a := range res.Data {
			url := fmt.Sprintf("%s/anime/hianime/info?id=%s", config.ConsumetUrl, a.ID)
			var res model.ConsumetAnimeWithMAL

			if err := doJSONRequest(url, &res); err != nil {
				continue
			}

			if res.MalID == anime.SeriesAnimeDBID {

				status, err := convertMalToStatus(anime.MyStatus)
				if err != nil {
					continue
				}

				collection := model.Collection{
					Type:    status,
					AnimeID: a.ID,
				}

				r.collectionRepo.AddCollection(deviceID, collection)
				count++

				now := time.Now()

				if anime.MyWatchedEpisodes > res.TotalEpisodes || anime.MyWatchedEpisodes == 0 {
					continue
				}

				history := model.History{
					AnimeID:            a.ID,
					LastWatchedEpisode: anime.MyWatchedEpisodes,
					IsWatched:          anime.MyWatchedEpisodes == res.TotalEpisodes,
					WatchedAt:          &now,
				}
				r.historyRepo.AddHistory(deviceID, history)

				for number := 0; number < anime.MyWatchedEpisodes; number++ {
					episode := res.Episodes[number]

					timecode := model.Timecode{
						EpisodeID: episode.ID,
						IsWatched: true,
						Time:      0,
						AnimeID:   res.ID,
					}

					r.timecodeRepo.AddTimecode(deviceID, timecode)
				}
				break
			}
		}
	}
	return count, nil
}

func (r *MALRepo) ExportMALList(deviceID string) (string, error) {
	rows, err := r.dbPostgres.Query(
		"SELECT anime_id, type FROM collections WHERE device_id = $1",
		deviceID,
	)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	animes := make([]model.MALAnime, 0)
	for rows.Next() {
		var animeID string
		var animeStatus string
		if err := rows.Scan(&animeID, &animeStatus); err != nil {
			continue
		}
		_, err := strconv.Atoi(animeID)
		if err == nil {
			continue
		}

		url := fmt.Sprintf("%s/anime/hianime/info?id=%s", config.ConsumetUrl, animeID)
		var res model.MALListAnime

		if err := doJSONRequest(url, &res); err != nil {
			continue
		}

		status, err := convertToMALStatus(animeStatus)
		if err != nil {
			continue
		}

		row, err := r.dbPostgres.Query(
			"SELECT last_watched FROM history WHERE device_id = $1 AND anime_id = $2", deviceID, animeID,
		)

		var lastWatched int
		if err == nil && row.Next() {
			if err := row.Scan(&lastWatched); err != nil {
				lastWatched = 0
			}
		} else {
			lastWatched = 0
		}

		animes = append(animes, model.MALAnime{
			SeriesAnimeDBID:   res.AnimeID,
			SeriesTitle:       res.Title,
			MyWatchedEpisodes: lastWatched,
			MyStatus:          status,
			UpdateOnImport:    1,
		})
	}

	if err = rows.Err(); err != nil {
		return "", err
	}

	malList := model.MALList{
		MyInfo: model.MyInfo{UserExportType: 1},
		Animes: animes,
	}

	output, err := xml.MarshalIndent(malList, "", "    ")
	if err != nil {
		return "", err
	}

	xmlWithHeader := xml.Header + string(output)

	return xmlWithHeader, nil
}
