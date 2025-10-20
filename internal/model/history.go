package model

type History struct {
	AnimeID            string `json:"anime_id"`
	LastWatchedEpisode int    `json:"last_watched"`
	IsWatched          bool   `json:"is_watched"`
}
