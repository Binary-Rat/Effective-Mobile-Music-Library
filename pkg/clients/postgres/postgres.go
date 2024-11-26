package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

	return dbpool, err
}
