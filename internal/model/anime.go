package model

type Anime struct {
	ID            string           `json:"id"`
	Title         string           `json:"title"`
	Poster        string           `json:"image"`
	Description   string           `json:"description"`
	Genres        []ConsumetGenre  `json:"genres"`
	Status        string           `json:"status"`
	Year          int              `json:"year"`
	Type          string           `json:"type"`
	TotalEpisodes int              `json:"total_episodes"`
	Episodes      []PreviewEpisode `json:"episodes"`
}
