package models

import "time"

type Song struct {
	Group   string `json:"group"`
	Song    string `json:"song"`
	Details SongDetail
}

type SongDetail struct {
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type Verse string
