package storage

import "Effective-Mobile-Music-Library/internal/models"

type Interface interface {
	Songs() ([]models.Song, error)
	AddSong() error
	ChangeSong() error
	Delete() error
	Text() (string, error)
}
