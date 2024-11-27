package models

import "time"

type Song struct {
	Group   string     `json:"group"`
	Song    string     `json:"song"`
	Details SongDetail `json:"details"`
}

type SongDetail struct {
	ReleaseDate time.Time `json:"releaseDate"`
	Lyrics      string    `json:"text"`
	Link        string    `json:"link"`
}

type Verse string
