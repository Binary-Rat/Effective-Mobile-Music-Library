package postgres

import (
	"Effective-Mobile-Music-Library/internal/models"
	"context"
	"fmt"
	"log"
	"strings"

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
func (s *Storage) Songs(ctx context.Context, reqSong models.SongDTO, offset uint64, limit uint64) ([]models.SongDTO, error) {

	var songs []models.SongDTO

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot start transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	sql, args, err := buildSQLWithFilter(reqSong, offset, limit)
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
		songs = append(songs, song)
	}
	tx.Commit(ctx)

	return songs, nil
}

func (s *Storage) Verses(ctx context.Context, id int, page uint64, limit uint64) ([]string, error) {
	var verse string
	err := s.db.QueryRow(ctx, `SELECT lyrics FROM songs WHERE id = $1`, id).Scan(&verse)
	if err != nil {
		return nil, fmt.Errorf("cannot scan row: %v", err)
	}
	if page != 0 {
		page--
	}
	verses := strings.Split(verse, "\n")[page*limit : limit*(page+1)]
	return verses, nil
}

func (s *Storage) DeleteSong(ctx context.Context, id int) error {
	_, err := s.db.Exec(ctx, `DELETE FROM songs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("cannot delete value: %v", err)
	}
	return nil
}

func (s *Storage) AddSong(ctx context.Context, song models.Song) (int, error) {
	var id int
	err := s.db.QueryRow(ctx, `INSERT INTO songs (band, song, release_date, lyrics, link) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		song.Group, song.Song, song.Details.ReleaseDate, song.Details.Lyrics, song.Details.Link).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot insert value: %v", err)
	}

	return id, nil
}

func (s *Storage) ChangeSong(ctx context.Context, song *models.SongDTO) (*models.SongDTO, error) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	prep := qb.Update("songs")
	if song.Group != "" {
		prep = prep.Set("band", song.Group)
	}
	if song.Song != "" {
		prep = prep.Set("song", song.Song)
	}
	if !song.Details.ReleaseDate.IsZero() {
		prep = prep.Set("release_date", song.Details.ReleaseDate)
	}
	if song.Details.Lyrics != "" {
		prep = prep.Set("lyrics", song.Details.Lyrics)
	}
	prep = prep.Where("id = $1", song.ID).Suffix("RETURNING id, band, song, release_date, lyrics")
	sql, args, err := prep.ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build query: %v", err)
	}

	err = s.db.QueryRow(ctx, sql, args...).Scan(&song.ID, &song.Group, &song.Song, &song.Details.ReleaseDate, &song.Details.Lyrics)
	if err != nil {
		return nil, fmt.Errorf("cannot update value: %v", err)
	}
	return song, nil
}

func buildSQLWithFilter(reqSong models.SongDTO, offset uint64, limit uint64) (string, []interface{}, error) {
	c := 1
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	prep := qb.Select("id", "band", "song", "release_date", "lyrics").From("songs").Offset(offset).Limit(limit)
	if reqSong.ID >= 0 {
		prep = prep.Where(fmt.Sprintf("id = $%d", c), reqSong.ID)
		c++
	}
	if reqSong.Group != "" {
		prep = prep.Where(fmt.Sprintf("band = $%d", c), reqSong.Group)
		c++
	}
	if reqSong.Song != "" {
		prep = prep.Where(fmt.Sprintf("song = $%d", c), reqSong.Song)
		c++
	}

	sql, args, err := prep.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("cannot build query: %v", err)
	}
	return sql, args, nil
}
