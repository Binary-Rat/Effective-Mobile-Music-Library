package api

import (
	"Effective-Mobile-Music-Library/internal/storage"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type api struct {
	r       *mux.Router
	l       *log.Logger
	storage storage.Interface
}

func New(router *mux.Router, logger *log.Logger, storage storage.Interface) *api {
	return &api{
		r:       router,
		l:       logger,
		storage: storage,
	}
}

func (a *api) Start() error {
	a.registerHandlers()
	return http.ListenAndServe(":8080", a.r)
}
