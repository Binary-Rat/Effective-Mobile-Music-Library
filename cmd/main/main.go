package main

import (
	"Effective-Mobile-Music-Library/internal/api"
	"Effective-Mobile-Music-Library/internal/storage/postgres"
	p "Effective-Mobile-Music-Library/pkg/clients/postgres"
	"Effective-Mobile-Music-Library/pkg/sources/songlib"
	"context"
	"log"

	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(log.Writer(), log.Prefix(), log.Flags())
	r := mux.NewRouter()
	pool, err := p.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	storage := postgres.New(pool)

	server := api.New(r, logger, storage, songlib.New())

	log.Fatal(server.Start())
}
