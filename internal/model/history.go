package model

import (
	"time"
)

type History struct {
	AnimeID            string     `json:"anime_id"`
	LastWatchedEpisode int        `json:"last_watched"`
	IsWatched          bool       `json:"is_watched"`
	WatchedAt          *time.Time `json:"watched_at"`
}
