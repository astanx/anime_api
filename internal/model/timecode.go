package model

type Timecode struct {
	EpisodeID string `json:"episode_id"`
	IsWatched bool   `json:"is_watched"`
	Time      int    `json:"time"`
}
