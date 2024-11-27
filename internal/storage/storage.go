package storage

import (
	"Effective-Mobile-Music-Library/internal/models"
	"context"
)

type Interface interface {
	Songs(context.Context, models.Song) ([]models.SongDTO, error)
	AddSong(ctx context.Context, song models.Song) (int, error)
	ChangeSong() error
	Delete() error
	Text() (string, error)
}
