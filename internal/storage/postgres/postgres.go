package postgres

import (
	"Effective-Mobile-Music-Library/internal/models"
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
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
func (s *Storage) Songs(ctx context.Context, reqSong models.Song) ([]models.SongDTO, error) {

	var songs []models.SongDTO

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot start transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	sql, args, err := buildSQLWithFilter(reqSong)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var song models.SongDTO
		if err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.Details.ReleaseDate, &song.Details.Lyrics); err != nil {
			log.Printf("cannot scan row: %v\\n", err)
		}
	}

	return songs, nil
}

func (s *Storage) Text() (string, error) {
	// implement logic here
	return "", nil
}

func (s *Storage) Delete() error {
	// implement logic here
	return nil
}

func (s *Storage) AddSong(ctx context.Context, song models.Song) (int, error) {
	var id int
	row, err := s.db.Query(ctx, `INSERT INTO songs (band, song, release_date, lyrics) 
		VALUES ($1, $2, $3, $4) RETURNING id`,
		song.Group, song.Song, song.Details.ReleaseDate, song.Details.Lyrics)

	if err != nil {
		return 0, fmt.Errorf("cannot insert value: %v", err)
	}
	defer row.Close()

	row.Scan(&id)
	return id, nil
}

func (s *Storage) ChangeSong() error {
	// implement logic here
	return nil
}

func buildSQLWithFilter(reqSong models.Song) (string, []interface{}, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	prep := qb.Select("id", "band", "song", "release_date", "lyrics").From("songs")
	if reqSong.Group != "" {
		prep = prep.Where("band = $1", reqSong.Group)
	}
	if reqSong.Song != "" {
		prep = prep.Where("song = $2", reqSong.Song)
	}
	if reqSong.Details.ReleaseDate.GoString() != "" {
		prep = prep.Where("release_date = $3", reqSong.Details.ReleaseDate)
	}

	sql, args, err := prep.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("cannot build query: %v", err)
	}
	return sql, args, nil
}
