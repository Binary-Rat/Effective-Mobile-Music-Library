package models

type SongDTO struct {
	ID      int         `json:"id"`
	Group   string      `json:"group"`
	Song    string      `json:"song"`
	Details SongDetails `json:"details"`
}
