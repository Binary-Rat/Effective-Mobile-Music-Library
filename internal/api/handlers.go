package api

import (
	"Effective-Mobile-Music-Library/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	musicEnpoint = "/music"
	verseEnpoint = "/music/verse"
)

func (api *api) registerHandlers() {
	//Get songs wtih filter
	api.r.HandleFunc(musicEnpoint, api.Songs).Methods(http.MethodGet).
		Queries("group", "{group}", "song", "{song}", "release", "{release}")
	//Create Song
	api.r.HandleFunc(musicEnpoint, api.AddSong).Methods(http.MethodPost)
	//Delete song
	api.r.HandleFunc(musicEnpoint, api.DeleteSong).Methods(http.MethodDelete).
		Queries("group", "{group}", "song", "{song}")
	//Get Verses of song
	api.r.HandleFunc(verseEnpoint, api.SongVerse).Methods(http.MethodGet).
		Queries("group", "{group}", "song", "{song}")
}

func (api *api) Songs(w http.ResponseWriter, r *http.Request) {
	songFilter := parseSongFilterFromURL(r)

	songs, err := api.storage.Songs(context.Background(), songFilter)
	if err != nil {
		api.l.Println(err)
	}

	body, err := json.Marshal(songs)
	if err != nil {
		api.l.Println(err)
	}

	//TODO: Make function to create response
	w.Write(body)
}

func (api *api) AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		api.l.Println(err)
	}

	api.source.SongWithDetails(context.Background(), &song)

	id, err := api.storage.AddSong(context.Background(), song)
	if err != nil {
		api.l.Println(err)
	}

	w.Write([]byte(fmt.Sprintf("Song added with id: %d", id)))
}

func (api *api) DeleteSong(w http.ResponseWriter, r *http.Request) {
	// implement logic here
}

func (api *api) SongVerse(w http.ResponseWriter, r *http.Request) {
	// implement logic here
}

func parseSongFilterFromURL(r *http.Request) models.Song {
	vars := mux.Vars(r)
	return models.Song{
		Group: vars["group"],
		Song:  vars["song"],
	}
}
