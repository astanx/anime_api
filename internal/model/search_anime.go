package model

type SearchAnime struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Poster     string `json:"image"`
	Year       int    `json:"year"`
	Type       string `json:"type"`
	ParserType string `json:"parser_type"`
}
