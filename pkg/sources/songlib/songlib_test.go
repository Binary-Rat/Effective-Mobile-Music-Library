package songlib

import (
	"Effective-Mobile-Music-Library/internal/models"
	mock_songlib "Effective-Mobile-Music-Library/pkg/sources/songlib/mocks"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_client_SongWithDetails(t *testing.T) {
	ctl := gomock.NewController(t)
	hc := mock_songlib.NewMockHTTPClient(ctl)
	client := client{
		client: hc,
	}
	details := &struct {
		ReleaseDate models.CustomTime `json:"releaseDate"`
		Lyrics      string            `json:"text"`
		Link        string            `json:"link"`
	}{
		ReleaseDate: models.CustomTime(time.Now()),
		Lyrics:      "lyrics",
		Link:        "link",
	}
	body, _ := json.Marshal(details)

	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)), //body,
	}

	ctx := context.Background()
	song := &models.Song{
		Group: "group",
		Song:  "song",
	}
	hc.EXPECT().Do(gomock.Any()).Return(mockResp, nil)

	t.Log(mockResp)

	client.SongWithDetails(ctx, song)

	if song.Details.ReleaseDate == (details.ReleaseDate) {
		t.Errorf("got %v, want %v", song.Details.ReleaseDate, details.ReleaseDate)
	}

	assert.Equal(t, details.Link, song.Details.Link)

	assert.Equal(t, details.Lyrics, song.Details.Lyrics)

}
