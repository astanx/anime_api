package model

type PreviewEpisode struct {
	ID      string `json:"id"`
	Ordinal int    `json:"ordinal"`
	Title   string `json:"title"`
}

type PreviewConsumetEpisode struct {
	ID      string `json:"id"`
	Ordinal int    `json:"number"`
	Title   string `json:"title"`
}

type TimeSegment struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
type Source struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}
type Subtitle struct {
	Vtt      string `json:"url"`
	Language string `json:"lang"`
}

type Episode struct {
	ID      string      `json:"id"`
	Ordinal int         `json:"ordinal"`
	Title   string      `json:"title"`
	Opening TimeSegment `json:"opening"`

	Ending TimeSegment `json:"ending"`

	Sources []Source `json:"sources"`

	Subtitles []Subtitle `json:"subtitles"`
}

type AnilibriaEpisode struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Ordinal int    `json:"ordinal"`
	Opening struct {
		Start int `json:"start"`
		End   int `json:"stop"`
	} `json:"opening"`

	Ending struct {
		Start int `json:"start"`
		End   int `json:"stop"`
	} `json:"ending"`
	Hls480  string `json:"hls_480"`
	Hls720  string `json:"hls_720"`
	Hls1080 string `json:"hls_1080"`
}

type ConsumetEpisode struct {
	ID      string      `json:"id"`
	Ordinal int         `json:"ordinal"`
	Title   string      `json:"title"`
	Intro   TimeSegment `json:"intro"`

	Outro TimeSegment `json:"outro"`

	Sources []Source `json:"sources"`

	Subtitles []Subtitle `json:"subtitles"`
}
