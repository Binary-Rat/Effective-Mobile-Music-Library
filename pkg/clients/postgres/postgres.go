package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Need refactor for flexibility
func NewClient(ctx context.Context) (pool *pgxpool.Pool, err error) {
	connectionString, exists := os.LookupEnv("DB_URL")
	if !exists {
		return nil, errors.New("environment variable DB_URL is not set")
	}

	dbpool, err := pgxpool.New(ctx, connectionString)

	if err != nil {
		return nil, err
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot ping DB due to error: %v", err)
	}

	err = MigrateUP(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("cannot migrate DB due to error: %v", err)
	}

	return dbpool, err
}

// I don`t know if it is okay. But I will be glad if you tell me how to do it right
// I think pgxpool potentially is great solution for this project,
// but i need the pure sql.DB to setup migrations
// i.e. I have no chose if i want to use pgxpool but this (

func MigrateUP(ctx context.Context, connString string) error {
	var db *sql.DB
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return err
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
