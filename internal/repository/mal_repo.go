package repository

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/astanx/anime_api/internal/db"
	"github.com/astanx/anime_api/internal/model"
	"github.com/redis/go-redis/v9"
)

type MALRepo struct {
	dbPostgres     *sql.DB
	dbClickhouse   clickhouse.Conn
	dbRedis        *redis.Client
	collectionRepo CollectionRepo
	historyRepo    HistoryRepo
	timecodeRepo   TimecodeRepo
	animeRepo      AnimeRepo
}

func NewMALRepo(db *db.DB, collectionRepo CollectionRepo, historyRepo HistoryRepo, timecodeRepo TimecodeRepo, animeRepo AnimeRepo) *MALRepo {
	return &MALRepo{
		dbPostgres:     db.Postgres,
		dbClickhouse:   db.ClickHouse,
		dbRedis:        db.Redis,
		collectionRepo: collectionRepo,
		historyRepo:    historyRepo,
		timecodeRepo:   timecodeRepo,
		animeRepo:      animeRepo,
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
		res, err := r.animeRepo.SearchConsumetAnime(anime.SeriesTitle, 1)

		if err != nil {
			continue
		}

		now := time.Now()

		for _, a := range res.Data {
			idAnime, err := r.animeRepo.GetAnimeInfoByConsumetID(a.ID)

			if err != nil {
				continue
			}

			fmt.Println("Parsed anime malID ", idAnime.MalID, "Trying to find ", anime.SeriesAnimeDBID)

			if idAnime.MalID == anime.SeriesAnimeDBID {
				status, err := convertMalToStatus(anime.MyStatus)
				if err != nil {
					break
				}

				collection := model.Collection{
					Type:    status,
					AnimeID: a.ID,
				}

				r.collectionRepo.AddCollection(deviceID, collection)
				count++

				if anime.MyWatchedEpisodes > idAnime.TotalEpisodes || anime.MyWatchedEpisodes == 0 {
					break
				}

				history := model.History{
					AnimeID:            a.ID,
					LastWatchedEpisode: anime.MyWatchedEpisodes,
					IsWatched:          anime.MyWatchedEpisodes == idAnime.TotalEpisodes,
					WatchedAt:          &now,
				}
				r.historyRepo.AddHistory(deviceID, history)

				for number := 0; number < anime.MyWatchedEpisodes; number++ {
					episode := idAnime.Episodes[number]

					timecode := model.Timecode{
						EpisodeID: episode.ID,
						IsWatched: true,
						Time:      0,
						AnimeID:   idAnime.ID,
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
		`SELECT c.anime_id, c.type, h.last_watched FROM collections as c 
		LEFT JOIN history as h 
		ON h.device_id = c.device_id AND h.anime_id = c.anime_id
		WHERE c.device_id = $1`,
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
		var lastWatched sql.NullInt64

		if err := rows.Scan(&animeID, &animeStatus, &lastWatched); err != nil {
			continue
		}

		_, err := strconv.Atoi(animeID)
		if err == nil {
			continue
		}

		idAnime, err := r.animeRepo.GetAnimeInfoByConsumetID(animeID)

		if err != nil {
			continue
		}

		status, err := convertToMALStatus(animeStatus)
		if err != nil {
			continue
		}

		watched := 0
		if lastWatched.Valid {
			watched = int(lastWatched.Int64)
		}

		animes = append(animes, model.MALAnime{
			SeriesAnimeDBID:   idAnime.MalID,
			SeriesTitle:       idAnime.Title,
			MyWatchedEpisodes: watched,
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

	output, err := xml.MarshalIndent(malList, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(output), nil
}
