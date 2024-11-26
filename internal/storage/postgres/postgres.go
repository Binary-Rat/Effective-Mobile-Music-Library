package postgres

import (
	"Effective-Mobile-Music-Library/internal/models"

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
func (s *Storage) Songs() ([]models.Song, error) {
	// implement logic here
	return nil, nil
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
