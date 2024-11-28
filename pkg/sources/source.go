package sources

import (
	"Effective-Mobile-Music-Library/internal/models"
	"context"
)

//go:generate mockgen -source=source.go -destination=mocks/source.go
type Source interface {
	SongWithDetails(ctx context.Context, song *models.Song) error
}
