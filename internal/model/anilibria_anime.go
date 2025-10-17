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
