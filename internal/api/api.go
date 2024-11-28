package api

import (
	"Effective-Mobile-Music-Library/internal/storage"
	"Effective-Mobile-Music-Library/pkg/sources"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type api struct {
	r       *mux.Router
	l       *log.Logger
	storage storage.Interface
	source  sources.Source
}

func New(router *mux.Router, logger *log.Logger, storage storage.Interface, source sources.Source) *api {
	return &api{
		r:       router,
		l:       logger,
		storage: storage,
		source:  source,
	}
}

func (a *api) Start() error {
	a.registerHandlers()
	return http.ListenAndServe(os.Getenv("PORT"), a.r)
}
