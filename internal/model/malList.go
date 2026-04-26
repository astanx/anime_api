package model

import "encoding/xml"

type MALList struct {
	XMLName xml.Name       `xml:"myanimelist"`
	MyInfo  MyInfo         `xml:"myinfo"`
	Animes  []MALListAnime `xml:"anime"`
}

type MyInfo struct {
	UserExportType int `xml:"user_export_type"`
}

type MALListAnime struct {
	SeriesAnimeDBID   int    `xml:"series_animedb_id"`
	SeriesTitle       string `xml:"series_title"`
	MyWatchedEpisodes int    `xml:"my_watched_episodes"`
	MyStatus          string `xml:"my_status"`
	UpdateOnImport    int    `xml:"update_on_import"`
}
