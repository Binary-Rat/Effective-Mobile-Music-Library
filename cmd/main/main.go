package main

import (
	"Effective-Mobile-Music-Library/internal/api"
	"Effective-Mobile-Music-Library/internal/storage/postgres"
	p "Effective-Mobile-Music-Library/pkg/clients/postgres"
	"Effective-Mobile-Music-Library/pkg/logger"
	"Effective-Mobile-Music-Library/pkg/sources/songlib"
	"context"
	"fmt"
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
	connectionString := mustConnString()
	if connectionString == "" {
		log.Fatal("Connection string is empty")
	}
	pool, err := p.NewClient(context.Background(), connectionString)
	if err != nil {
		log.Fatal(err)
	}
	storage := postgres.New(pool, logger)

	server := api.New(r, logger, storage, songlib.New())

	port := os.Getenv("PORT")
	logger.Infof("Server is running on port %s", port)
	log.Fatal(server.Start(port))
}

func mustConnString() string {

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
}
