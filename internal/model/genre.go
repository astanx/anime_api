package model

type Genre struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	TotalReleases int    `json:"total_releases"`
}
