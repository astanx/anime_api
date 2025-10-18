package model

type SearchAnilibriaAnime struct {
	ID   int `json:"id"`
	Type struct {
		Value string `json:"value"`
	} `json:"type"`
	Year   int `json:"year"`
	Poster struct {
		Optimized struct {
			Thumbnail string `json:"thumbnail"`
		} `json:"optimized"`
	} `json:"poster"`
	Name struct {
		Main string `json:"main"`
	} `json:"name"`
}

type AnilibriaAnime struct {
	ID   int `json:"id"`
	Type struct {
		Value string `json:"value"`
	} `json:"type"`
	Year int `json:"year"`
	Name struct {
		Main string `jaon:"main"`
	} `json:"name"`
	Poster struct {
		Optimized struct {
			Thumbnail string `json:"thumbnail"`
		} `json:"optimized"`
	} `json:"poster"`
	IsOngoing     bool   `json:"is_ongoing"`
	Description   string `json:"description"`
	EpisodesTotal int    `json:"episodes_total"`
	Genres        []struct {
		Name string `json:"name"`
	} `json:"genres"`
	Episodes []struct {
		ID      string  `json:"id"`
		Name    string  `json:"name"`
		Ordinal int     `json:"ordinal"`
		HLS480  string  `json:"hls_480"`
		HLS720  string  `json:"hls_720"`
		HLS1080 *string `json:"hls_1080"`
	} `json:"episodes"`
}
