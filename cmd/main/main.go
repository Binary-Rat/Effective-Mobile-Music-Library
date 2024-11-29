package main

import (
	"Effective-Mobile-Music-Library/internal/api"
	"Effective-Mobile-Music-Library/internal/storage/postgres"
	p "Effective-Mobile-Music-Library/pkg/clients/postgres"
	"Effective-Mobile-Music-Library/pkg/logger"
	"Effective-Mobile-Music-Library/pkg/sources/songlib"
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@host		localhost:8080
//	@BasePath	/

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

}

func main() {
	logger := logger.ConfigureLogger()
	r := mux.NewRouter()

	logger.Info("Connecting to database")
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

	logger.Infof("Server is running on port %s", os.Getenv("PORT"))
	log.Fatal(server.Start())
}
