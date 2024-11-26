package api

import "net/http"

const (
	musicEnpoint = "/music"
	verseEnpoint = "/music/verse"
)

func (api *api) registerHandlers() {
	api.r.HandleFunc(musicEnpoint, api.Songs).Methods(http.MethodGet).
		Queries("group", "{group}", "song", "{song}", "release", "{release}")
	api.r.HandleFunc(musicEnpoint, api.DeleteSong).Methods(http.MethodDelete).
		Queries("group", "{group}", "song", "{song}")
	api.r.HandleFunc(verseEnpoint, api.SongVerse).Methods(http.MethodGet).
		Queries("group", "{group}", "song", "{song}")
}

func (api *api) Songs(w http.ResponseWriter, r *http.Request) {
	// implement logic here
}

func (api *api) DeleteSong(w http.ResponseWriter, r *http.Request) {
	// implement logic here
}

func (api *api) SongVerse(w http.ResponseWriter, r *http.Request) {
	// implement logic here
}
