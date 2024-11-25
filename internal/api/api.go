package api

import (
	"Effective-Mobile-Music-Library/internal/storage"
	"log"

	"github.com/gorilla/mux"
)

type api struct {
	r       *mux.Router
	l       *log.Logger
	storage storage.Interface
}

func New(r *mux.Router, l *log.Logger, storage storage.Interface) *api {
	return &api{
		r:       r,
		l:       l,
		storage: storage,
	}
}

func (r *api) Start() error {
	return nil
}
