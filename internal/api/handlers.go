package api

import (
	"Effective-Mobile-Music-Library/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	_ "Effective-Mobile-Music-Library/docs"

	httpSwagger "github.com/swaggo/http-swagger"

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
	//SWAGGER
	api.r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
}

// Songs
// @Summary      Get songs
// @Description  Получение списка песен с возможностью фильтрации, пагинации и сортировки
// @Tags         Songs
// @Produce      json
// @Param        artist query string false "Фильтр по исполнителю"
// @Param        genre  query string false "Фильтр по жанру"
// @Param        page   query int    false "Номер страницы для пагинации" default(1)
// @Param        limit  query int    false "Количество записей на странице" default(10)
// @Success      200 {array} models.Song
// @Failure      500 {object} string "Internal Server Error"
// @Router       /music [get]
func (api *api) Songs(w http.ResponseWriter, r *http.Request) {
	api.l.Println("Geting songs...")

	songFilter, page, limit := parseSongFilterFromURL(r)

	songs, err := api.storage.Songs(context.Background(), songFilter, page, limit)
	if err != nil {
		panic(err)
	}

	body, err := json.Marshal(songs)
	if err != nil {
		panic(err)
	}

	//TODO: Make function to create response
	w.Write(body)
}

// AddSong
// @Summary      Add a song
// @Description  Добавление новой песни
// @Tags         Songs
// @Accept       json
// @Produce      text/plain
// @Param        song body models.Song true "Данные новой песни"
// @Success      201 {string} string "Song added with id: {id}"
// @Failure      400 {object} string "Bad Request"
// @Failure      500 {object} string "Internal Server Error"
// @Router       /music [post]
func (api *api) AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		panic(err)
	}

	err = api.source.SongWithDetails(context.Background(), &song)
	if err != nil {
		panic(err)
	}
	id, err := api.storage.AddSong(context.Background(), song)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf("Song added with id: %d", id)))
}

// DeleteSong
// @Summary      Delete a song
// @Description  Удаление песни по идентификатору
// @Tags         Songs
// @Produce      text/plain
// @Param        id query int true "Идентификатор песни"
// @Success      200 {string} string "Song deleted"
// @Failure      400 {object} string "Bad Request"
// @Failure      500 {object} string "Internal Server Error"
// @Router       /music [delete]
func (api *api) DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	err = api.storage.DeleteSong(context.Background(), id)
	if err != nil {
		panic(err)
	}
}

// SongVerse
// @Summary      Get song verses
// @Description  Получение стихов песни по идентификатору
// @Tags         Songs
// @Produce      json
// @Param        id    path int true  "Идентификатор песни"
// @Param        page  query int false "Номер страницы для пагинации" default(1)
// @Param        limit query int false "Количество стихов на странице" default(10)
// @Success      200 {array} models.Verse
// @Failure      400 {object} string "Bad Request"
// @Failure      500 {object} string "Internal Server Error"
// @Router       /music/{id}/verse [get]
func (api *api) SongVerse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["page"])
	if err != nil {
		panic(err)
	}
	offset, limit := parsePageLimit(r)
	verse, err := api.storage.Verses(context.Background(), id, offset, limit)
	if err != nil {
		panic(err)
	}
	body, err := json.Marshal(verse)
	if err != nil {
		panic(err)
	}

	w.Write(body)
}

// ChangeSong
// @Summary      Change a song
// @Description  Изменение данных песни
// @Tags         Songs
// @Accept       json
// @Produce      text/plain
// @Param        song body models.SongDTO true "Данные для изменения песни"
// @Success      200 {string} string "Song changed with id: {id}"
// @Failure      400 {object} string "Bad Request"
// @Failure      500 {object} string "Internal Server Error"
// @Router       /music [patch]
func (api *api) ChangeSong(w http.ResponseWriter, r *http.Request) {
	var song models.SongDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &song)
	if err != nil {
		panic(err)
	}
	newSong, err := api.storage.ChangeSong(context.Background(), &song)
	if err != nil {
		panic(err)
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
