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

	"errors"
)

type client struct {
	client *http.Client
}

func New() *client {
	return &client{
		client: &http.Client{},
	}
}

func (c *client) SongWithDetails(ctx context.Context, song *models.Song) error {
	var sd models.SongDetail

	query := map[string]string{
		"group": song.Group,
		"song":  song.Song,
	}

	body, err := c.doHTTP(ctx, "/info", query, nil)
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

var host = "https://localhost:8080"

func (c *client) doHTTP(ctx context.Context, method string, query map[string]string, body interface{}) ([]byte, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %v", err)
	}
	url := configureURL(query)

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

func configureURL(query map[string]string) string {
	u := url.URL{}
	u.Path = host
	for k, v := range query {
		u.Query().Add(k, v)
	}

	return u.String()
}
