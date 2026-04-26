package model

type Anime struct {
	ID            string           `json:"id"`
	MalID         int              `json:"malID"`
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

type ConsumetAnimeWithMAL struct {
	ConsumetAnime
	MalID int `json:"malID"`
}

type MALAnime struct {
	ID     int    `json:"mal_id"`
	Title  string `json:"title"`
	Titles []struct {
		Type  string `json:"type"`
		Title string `json:"title"`
	} `json:"titles"`
	Images struct {
		Jpg struct {
			ImageURL string `json:"image_url"`
		} `json:"jpg"`
		Webp struct {
			ImageURL string `json:"image_url"`
		} `json:"webp"`
	} `json:"images"`
	Year        int    `json:"year"`
	Type        string `json:"type"`
	Description string `json:"synopsis"`
	Status      string `json:"status"`
	Genres      []struct {
		Name string `json:"name"`
	} `json:"genres"`
	TotalEpisodes int `json:"episodes"`
}

type MALLatestAnime struct {
	Entry struct {
		ID     int    `json:"mal_id"`
		Title  string `json:"title"`
		Images struct {
			Jpg struct {
				ImageURL string `json:"image_url"`
			} `json:"jpg"`
			Webp struct {
				ImageURL string `json:"webp"`
			} `json:"webp"`
		} `json:"images"`
	} `json:"entry"`
}

type ProwlarrAnime struct {
	GUID      string `json:"guid"`
	Size      int64  `json:"size"`
	Grabs     int    `json:"grabs"`
	Seeders   int    `json:"seeders"`
	Title     string `json:"title"`
	SortTitle string `json:"sortTitle"`
	MagnetURL string `json:"magnetUrl"`
	Hash      string `json:"infoHash"`
}
