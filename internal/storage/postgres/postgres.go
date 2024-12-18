package postgres

import (
	"Effective-Mobile-Music-Library/internal/models"
	"context"
	"fmt"
	"sync"
	"time"

	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	db  *pgxpool.Pool
	log *logrus.Logger
	sync.Mutex
}

func New(db *pgxpool.Pool, l *logrus.Logger) *Storage {
	return &Storage{
		db:  db,
		log: l,
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
		var t time.Time
		//ID band song release_date link lyrics
		if err := rows.Scan(&song.ID, &song.Group, &song.Song, &t, &song.Details.Link, &song.Details.Lyrics); err != nil {
			s.log.Errorf("cannot scan row: %v\n", err)
		}
		song.Details.ReleaseDate = models.CustomTime(t)
		songs = append(songs, song)
	}
	tx.Commit(ctx)

	return songs, nil
}

func (s *Storage) Verses(ctx context.Context, id int, offset uint64, limit uint64) ([]string, error) {
	var verse string
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot start transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(ctx, `SELECT lyrics FROM songs WHERE id = $1`, id).Scan(&verse)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return []string{}, nil
		}
		return nil, fmt.Errorf("cannot scan row: %v", err)
	}
	verses := strings.Split(verse, "\n")
	s.log.Debug("verses: ", verses)
	if int(offset) > len(verses) {
		return []string{}, nil
	}
	if int(offset+limit) > len(verses) {
		return verses[offset:], nil
	}
	tx.Commit(ctx)
	return verses[offset : offset+limit], nil
}

func (s *Storage) DeleteSong(ctx context.Context, id int) error {
	s.Lock()
	defer s.Unlock()
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cannot start transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	cmdTag, err := tx.Exec(ctx, `DELETE FROM songs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("cannot delete value: %v", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("there is no value with id: %v", id)
	}

	tx.Commit(ctx)
	return nil
}

func (s *Storage) AddSong(ctx context.Context, song models.Song) (int, error) {
	s.Lock()
	defer s.Unlock()
	var id int
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("cannot start transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(ctx, `INSERT INTO songs (band, song, release_date, lyrics, link) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		song.Group, song.Song, time.Time(song.Details.ReleaseDate), song.Details.Lyrics, song.Details.Link).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot insert value: %v", err)
	}
	tx.Commit(ctx)
	return id, nil
}

func (s *Storage) ChangeSong(ctx context.Context, song *models.SongDTO) (*models.SongDTO, error) {
	prep := prepSQLUpdate(song)
	sql, args, err := prep.ToSql()
	if err != nil {
		return nil, fmt.Errorf("cannot build query: %v", err)
	}
	var t time.Time
	s.log.Debugf("sql:%s args:%s", sql, args)
	s.Lock()
	defer s.Unlock()
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot start transaction: %v", err)
	}
	defer tx.Rollback(context.Background())
	err = tx.QueryRow(ctx, sql, args...).Scan(&song.ID, &song.Group, &song.Song, &t, &song.Details.Link, &song.Details.Lyrics)
	if err != nil {
		return nil, fmt.Errorf("cannot update value: %v", err)
	}
	song.Details.ReleaseDate = models.CustomTime(t)
	tx.Commit(ctx)
	return song, nil
}

func buildSQLWithFilter(reqSong models.SongDTO, offset uint64, limit uint64) (string, []interface{}, error) {
	c := 1
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	prep := qb.Select("id", "band", "song", "release_date", "link", "lyrics").From("songs").Offset(offset).Limit(limit)
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

func prepSQLUpdate(song *models.SongDTO) sq.UpdateBuilder {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	//counter for fillter to put correct var in query
	c := 1
	prep := qb.Update("songs")
	if song.Group != "" {
		c++
		prep = prep.Set("band", song.Group)
	}
	if song.Song != "" {
		c++
		prep = prep.Set("song", song.Song)
	}
	if !time.Time(song.Details.ReleaseDate).IsZero() {
		c++
		prep = prep.Set("release_date", time.Time(song.Details.ReleaseDate))
	}
	if song.Details.Lyrics != "" {
		c++
		prep = prep.Set("lyrics", song.Details.Lyrics)
	}
	if song.Details.Link != "" {
		c++
		prep = prep.Set("link", song.Details.Lyrics)
	}
	prep = prep.Where(fmt.Sprintf("id=$%d", c), song.ID).Suffix("RETURNING id, band, song, release_date, link, lyrics")
	return prep
}
