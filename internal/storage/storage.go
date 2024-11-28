package storage

import (
	"Effective-Mobile-Music-Library/internal/models"
	"context"
)

type Interface interface {
	Songs(context.Context, models.SongDTO) ([]models.SongDTO, error)
	AddSong(ctx context.Context, song models.Song) (int, error)
	ChangeSong(ctx context.Context, song models.SongDTO) (*models.SongDTO, error)
	Delete(ctx context.Context, id int) error
	Text(ctx context.Context, id int) (string, error)
}
