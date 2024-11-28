package storage

import (
	"Effective-Mobile-Music-Library/internal/models"
	"context"
)

type Interface interface {
	Songs(ctx context.Context, reqSong models.SongDTO, page uint64, limit uint64) ([]models.SongDTO, error)
	AddSong(ctx context.Context, song models.Song) (int, error)
	ChangeSong(ctx context.Context, song models.SongDTO) (*models.SongDTO, error)
	DeleteSong(ctx context.Context, id int) error
	Verses(ctx context.Context, id int, page uint64, limit uint64) ([]string, error)
}
