package api

import (
	"Effective-Mobile-Music-Library/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	musicEnpoint = "/music"
)

func (api *api) registerHandlers() {
	//Get songs wtih filter
	api.r.HandleFunc(musicEnpoint, api.Songs).Methods(http.MethodGet)
	//Create Song
	api.r.HandleFunc(musicEnpoint, api.AddSong).Methods(http.MethodPost)
	//Delete song
	api.r.HandleFunc(musicEnpoint, api.DeleteSong).Methods(http.MethodDelete).Queries("id", "{id}")
	//Change song
	api.r.HandleFunc(musicEnpoint, api.ChangeSong).Methods(http.MethodPatch)
	//Get Verses of song
	api.r.HandleFunc(fmt.Sprintf("%s/{id}/verse", musicEnpoint), api.SongVerse).Methods(http.MethodGet)
}

func (api *api) Songs(w http.ResponseWriter, r *http.Request) {
	api.l.Println("Geting songs...")

	songFilter, page, limit := parseSongFilterFromURL(r)

	songs, err := api.storage.Songs(context.Background(), songFilter, page, limit)
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

	err = api.source.SongWithDetails(context.Background(), &song)
	if err != nil {
		api.l.Println(err)
	}
	id, err := api.storage.AddSong(context.Background(), song)
	if err != nil {
		api.l.Println(err)
	}

	w.Write([]byte(fmt.Sprintf("Song added with id: %d", id)))
}

func (api *api) DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.l.Println(err)
	}
	err = api.storage.DeleteSong(context.Background(), id)
	if err != nil {
		api.l.Println(err)
	}
}

func (api *api) SongVerse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["page"])
	if err != nil {
		api.l.Println(err)
	}
	offset, limit := parsePageLimit(r)
	verse, err := api.storage.Verses(context.Background(), id, offset, limit)
	if err != nil {
		api.l.Println(err)
	}
	body, err := json.Marshal(verse)
	if err != nil {
		api.l.Println(err)
		return
	}

	w.Write(body)
}

func (api *api) ChangeSong(w http.ResponseWriter, r *http.Request) {
	var song models.SongDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api.l.Println(err)
		return
	}
	err = json.Unmarshal(body, &song)
	if err != nil {
		api.l.Println(err)
		return
	}
	newSong, err := api.storage.ChangeSong(context.Background(), song)
	if err != nil {
		api.l.Println(err)
	}

	w.Write([]byte(fmt.Sprintf("Song changed with id: %d", newSong.ID)))
}

func parseSongFilterFromURL(r *http.Request) (models.SongDTO, uint64, uint64) {
	group := r.URL.Query().Get("group")
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		id = -1
	}
	song := r.URL.Query().Get("song")

	page, limit := parsePageLimit(r)

	return models.SongDTO{
		ID:    id,
		Group: group,
		Song:  song,
	}, page, limit
}

func parsePageLimit(r *http.Request) (uint64, uint64) {
	page, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		page = 0
	}
	limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}
	return page, limit
}
