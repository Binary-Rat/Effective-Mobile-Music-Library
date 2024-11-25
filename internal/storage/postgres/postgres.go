package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{
		db: db,
	}
}
func (s *Storage) Song() (string, error) {
	// implement logic here
	return "", nil
}

func (s *Storage) Text() (string, error) {
	// implement logic here
	return "", nil
}

func (s *Storage) Delete() error {
	// implement logic here
	return nil
}

func (s *Storage) AddSong() error {
	// implement logic here
	return nil
}

func (s *Storage) ChangeSong() error {
	// implement logic here
	return nil
}
