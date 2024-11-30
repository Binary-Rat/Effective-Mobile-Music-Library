package songlib

import (
	"Effective-Mobile-Music-Library/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"

	"errors"
)

// For tests
//
//go:generate mockgen -source=songlib.go -destination=mocks/songlib.go
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type client struct {
	client HTTPClient
}

func New() *client {
	return &client{
		client: &http.Client{},
	}
}

func (c *client) SongWithDetails(ctx context.Context, song *models.Song) error {
	var sd models.SongDetails

	var host = os.Getenv("SOURCE_HOST")
	//for making url query
	query := map[string]string{
		"group": song.Group,
		"song":  song.Song,
	}

	url := configureURL(host, query)

	body, err := c.doHTTP(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to get song detail: %v", err)
	}
	err = json.Unmarshal(body, &sd)
	if err != nil {
		return fmt.Errorf("failed to unmarshal body: %v", err)
	}

	song.Details = sd

	return nil
}

func (c *client) doHTTP(ctx context.Context, method string, url string, body interface{}) ([]byte, error) {
	log.Infof("HTTP request: %s %s", method, url)
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Sprintf("API Error: %s", resp.Status)
		return nil, errors.New(err)
	}

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return respB, nil
}

func configureURL(host string, query map[string]string) string {
	u, err := url.Parse(host)
	if err != nil {
		return ""
	}
	for k, v := range query {
		(*u).Query().Add(k, v)
	}

	return u.String()
}
