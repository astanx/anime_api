package model

type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	PagesLeft  int `json:"pagesLeft"`
	TotalPages int `json:"totalPages"`
}

type PaginatedFavourites struct {
	Data []Favourite    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type PaginatedCollections struct {
	Data []Collection   `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type PaginatedHistory struct {
	Data []History      `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
