package main

import (
	"Effective-Mobile-Music-Library/internal/api"
	"Effective-Mobile-Music-Library/internal/storage/postgres"
	p "Effective-Mobile-Music-Library/pkg/clients/postgres"
	"Effective-Mobile-Music-Library/pkg/sources/songlib"
	"context"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

}

func main() {
	logger := log.New(log.Writer(), log.Prefix(), log.Flags())
	r := mux.NewRouter()

	connectionString, exists := os.LookupEnv("DB_URL")
	if !exists {
		log.Fatal("environment variable DB_URL is not set")
	}
	pool, err := p.NewClient(context.Background(), connectionString)
	if err != nil {
		log.Fatal(err)
	}

	storage := postgres.New(pool)

	server := api.New(r, logger, storage, songlib.New())

	port := os.Getenv("PORT")

	logger.Printf("Server is running on port %s", port)
	log.Fatal(server.Start())
}
