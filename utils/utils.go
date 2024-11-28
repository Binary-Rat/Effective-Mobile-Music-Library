package utils

import "Effective-Mobile-Music-Library/internal/models"

func parseSongToMap(song models.Song) map[string]interface{} {
	return map[string]interface{}{
		"group":   song.Group,
		"song":    song.Song,
		"details": song.Details,
	}
}
