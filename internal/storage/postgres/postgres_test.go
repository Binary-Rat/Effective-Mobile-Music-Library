package postgres_test

import (
	"Effective-Mobile-Music-Library/internal/models"
	"Effective-Mobile-Music-Library/internal/storage/postgres"
	p "Effective-Mobile-Music-Library/pkg/clients/postgres"
	"context"
	"fmt"
	"os"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

const dbTestURL = "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

var db *postgres.Storage

func TestMain(m *testing.M) {
	pool, _ := p.NewClient(context.Background(), dbTestURL)

	db = postgres.New(pool)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestQueryBuilder(t *testing.T) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	users := qb.Select("*").From("users")

	active := users.Where("active = $1", "A")
	active = active.Where("abc = $1", "B")

	sql, args, err := active.ToSql()

	t.Log(fmt.Sprintf(sql, args), err)
}

func TestSongs(t *testing.T) {
	if db == nil {
		t.Fatal("db is not initialized")
	}
	song := models.Song{
		Group: "Test",
		Song:  "Test",
		Details: models.SongDetails{
			ReleaseDate: "",
			Lyrics:      "Test",
			Link:        "http://test.com",
		},
	}

	id, err := db.AddSong(context.Background(), song)
	if err != nil {
		t.Fatal(err)
	}

	reqSong := models.SongDTO{
		ID: id,
	}
	songAdded, err := db.Songs(context.Background(), reqSong, 0, 1)
	if err != nil {
		t.Fatal(err)
	}

	err = db.DeleteSong(context.Background(), reqSong.ID)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(songAdded)

	assert.Equal(t, songAdded[0].Song, song.Song)
}
