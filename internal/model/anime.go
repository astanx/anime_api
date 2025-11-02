package model

type Anime struct {
	ID            string           `json:"id"`
	Title         string           `json:"title"`
	Poster        string           `json:"image"`
	Description   string           `json:"description"`
	Genres        []string         `json:"genres"`
	Status        string           `json:"status"`
	Year          int              `json:"year"`
	Type          string           `json:"type"`
	TotalEpisodes int              `json:"total_episodes"`
	Episodes      []PreviewEpisode `json:"episodes"`
}

type ConsumetAnime struct {
	ID            string                   `json:"id"`
	Title         string                   `json:"title"`
	Image         string                   `json:"image"`
	Description   string                   `json:"description"`
	Genres        []string                 `json:"genres"`
	Status        string                   `json:"status"`
	Type          string                   `json:"type"`
	TotalEpisodes int                      `json:"totalEpisodes"`
	Episodes      []PreviewConsumetEpisode `json:"episodes"`
}
