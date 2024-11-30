package api

import (
	"Effective-Mobile-Music-Library/internal/models"
	"Effective-Mobile-Music-Library/pkg/middleware"
	aErr "Effective-Mobile-Music-Library/pkg/middleware/app-err"
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
	api.r.Use(middleware.ErrorMiddleware)
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
// @Param        offset   query int    false "Номер страницы для пагинации" default(1)
// @Param        limit  query int    false "Количество записей на странице" default(10)
// @Success      200 {array} models.Song "Список песен"
// @Failure      400 {object} appErr.Error "Ошибка запроса"
// @Failure      500 {object} appErr.Error "Внутренняя ошибка сервера"
// @Router       /music [get]
func (api *api) Songs(w http.ResponseWriter, r *http.Request) {
	api.l.Debug("Geting songs...")
	songFilter, offset, limit := parseSongFilterFromURL(r)

	songs, err := api.storage.Songs(context.Background(), songFilter, offset, limit)
	if err != nil {
		api.l.Error(err)
		panic(aErr.New(http.StatusInternalServerError, err.Error()))
	}

	body, err := json.Marshal(songs)
	if err != nil {
		api.l.Error(err)
		panic(aErr.New(http.StatusInternalServerError, err.Error()))
	}

	//TODO: Make function to create response
	writeRespone(w, body)
}

// AddSong
// @Summary      Add a song
// @Description  Добавление новой песни
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        song body models.Song true "Данные новой песни"
// @Success      201 {object} string "Song added with id: {id}"
// @Failure      400 {object} appErr.Error "Ошибка запроса"
// @Failure      422 {object} appErr.Error "Невалидные данные"
// @Failure      500 {object} appErr.Error "Внутренняя ошибка сервера"
// @Router       /music [post]
func (api *api) AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		api.l.Error(err)
		panic(err)
	}

	err = api.source.SongWithDetails(context.Background(), &song)
	if err != nil {
		api.l.Error(err)
		panic(aErr.New(http.StatusInternalServerError, err.Error()))
	}
	id, err := api.storage.AddSong(context.Background(), song)
	if err != nil {
		api.l.Error(err)
		panic(aErr.New(http.StatusInternalServerError, err.Error()))
	}

	body := []byte(fmt.Sprintf("Song added - id: %d", id))
	writeRespone(w, body)
}

// DeleteSong
// @Summary      Delete a song
// @Description  Удаление песни по идентификатору
// @Tags         Songs
// @Produce      json
// @Param        id query int true "Идентификатор песни"
// @Success      200 {object} string "Song deleted"
// @Failure      400 {object} appErr.Error "Ошибка запроса"
// @Failure      404 {object} appErr.Error "Песня не найдена"
// @Failure      500 {object} appErr.Error "Внутренняя ошибка сервера"
// @Router       /music [delete]
func (api *api) DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.l.Error(err)
		panic(aErr.New(http.StatusBadRequest, err.Error()))
	}
	err = api.storage.DeleteSong(context.Background(), id)
	if err != nil {
		api.l.Error(err)
		panic(aErr.New(http.StatusInternalServerError, err.Error()))
	}

	body := []byte(fmt.Sprintf("Song with id %d deleted", id))
	writeRespone(w, body)
}

// SongVerse
// @Summary      Get song verses
// @Description  Получение стихов песни по идентификатору
// @Tags         Verses
// @Produce      json
// @Param        id    path int true  "Идентификатор песни"
// @Param        offset  query int false "Номер страницы для пагинации" default(1)
// @Param        limit query int false "Количество стихов на странице" default(10)
// @Success      200 {array} models.Verse "Список стихов"
// @Failure      400 {object} appErr.Error "Ошибка запроса"
// @Failure      404 {object} appErr.Error "Песня не найдена"
// @Failure      500 {object} appErr.Error "Внутренняя ошибка сервера"
// @Router       /music/{id}/verse [get]
func (api *api) SongVerse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.l.Error(err)
		panic(aErr.New(http.StatusBadRequest, err.Error()))
	}
	offset, limit := parseOffsetLimit(r)
	verse, err := api.storage.Verses(context.Background(), id, offset, limit)
	if err != nil {
		api.l.Error(err)
		panic(aErr.New(http.StatusInternalServerError, err.Error()))
	}
	body, err := json.Marshal(verse)
	if err != nil {
		api.l.Error(err)
		panic(aErr.New(http.StatusInternalServerError, err.Error()))
	}

	writeRespone(w, body)
}

// ChangeSong
// @Summary      Change a song
// @Description  Изменение данных песни
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        song body models.SongDTO true "Данные для изменения песни"
// @Success      200 {object} string "Song changed with id: {id}"
// @Failure      400 {object} appErr.Error "Ошибка запроса"
// @Failure      404 {object} appErr.Error "Песня не найдена"
// @Failure      422 {object} appErr.Error "Невалидные данные"
// @Failure      500 {object} appErr.Error "Внутренняя ошибка сервера"
// @Router       /music [patch]
func (api *api) ChangeSong(w http.ResponseWriter, r *http.Request) {
	var song models.SongDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api.l.Error(err)
		panic(aErr.New(http.StatusBadRequest, err.Error()))
	}
	err = json.Unmarshal(body, &song)
	if err != nil {
		api.l.Error(err)
		panic(aErr.New(http.StatusBadRequest, err.Error()))
	}
	newSong, err := api.storage.ChangeSong(context.Background(), &song)
	if err != nil {
		api.l.Error(err)
		panic(aErr.New(http.StatusInternalServerError, err.Error()))
	}

	body = []byte(fmt.Sprintf("Song changed with id: %d", newSong.ID))

	writeRespone(w, body)
}

func parseSongFilterFromURL(r *http.Request) (models.SongDTO, uint64, uint64) {
	group := r.URL.Query().Get("group")
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		id = -1
	}
	song := r.URL.Query().Get("song")

	offset, limit := parseOffsetLimit(r)

	return models.SongDTO{
		ID:    id,
		Group: group,
		Song:  song,
	}, offset, limit
}

func parseOffsetLimit(r *http.Request) (uint64, uint64) {
	offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}
	return offset, limit
}

func writeRespone(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
