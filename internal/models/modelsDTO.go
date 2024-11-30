package models

type SongDTO struct {
	ID      int         `json:"id" example:"1"`
	Group   string      `json:"group" example:"Nirvana"`
	Song    string      `json:"song" example:"Lithium"`
	Details SongDetails `json:"details"`
}
