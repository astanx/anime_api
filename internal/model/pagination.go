package model

type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	PagesLeft  int `json:"pages_left"`
	TotalPages int `json:"total_pages"`
}

type ShortPaginationMeta struct {
	CurrentPage int  `json:"current_page"`
	HasNextPage bool `json:"has_next_page"`
	TotalPages  int  `json:"total_pages"`
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

type PaginatedSearchAnime struct {
	Data []SearchAnime       `json:"data"`
	Meta ShortPaginationMeta `json:"meta"`
}

type PaginatedAnilibriaSearchAnime struct {
	Data []SearchAnilibriaAnime `json:"data"`
	Meta PaginationMeta         `json:"meta"`
}

type PaginatedConsumetSearchAnime struct {
	Data        []SearchAnime `json:"results"`
	CurrentPage int           `json:"currentPage"`
	TotalPages  int           `json:"totalPages"`
	HasNextPage bool          `json:"hasNextPage"`
}
