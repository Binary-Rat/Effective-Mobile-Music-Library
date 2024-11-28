package sources

import (
	"Effective-Mobile-Music-Library/internal/models"
	"context"
)

type Source interface {
	SongWithDetails(ctx context.Context, song *models.Song) error
}
