package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/marcus-crane/wwenet/config"
	"github.com/marcus-crane/wwenet/networkentities"
)

type Client struct {
	token      string
	config     config.Config
	httpClient *http.Client
	baseURL    string
}

func NewClient(token string, config config.Config) *Client {
	return &Client{
		token:  token,
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://dce-frontoffice.imggaming.com/api",
	}
}

func (c *Client) GetEpisode(ctx context.Context, episodeID int64) (*networkentities.Episode, error) {
	url := fmt.Sprintf("%s/v4/vod/%d?includePlaybackDetails=URL", c.baseURL, episodeID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Add("Realm", "dce.wwe")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", c.config.Network.XApiKey)
	req.Header.Add("x-app-var", c.config.Network.XAppVar)
	req.Header.Add("User-Agent", c.config.Network.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &networkentities.Episode{}, fmt.Errorf("failed to read episode info. got status code %d", resp.StatusCode)
	}

	var episode networkentities.Episode
	if err := json.NewDecoder(resp.Body).Decode(&episode); err != nil {
		return nil, fmt.Errorf("failed to decode episode: %w", err)
	}

	return &episode, nil
}

func (c *Client) GetSeason(ctx context.Context, seasonID int64) (*networkentities.Season, error) {
	baseSeason, err := c.getSeasonPartial(ctx, seasonID, 0)
	if err != nil {
		return &networkentities.Season{}, fmt.Errorf("failed to fetch initial payload for season ID %d: %w", seasonID, err)
	}
	pagesRemaining := baseSeason.Paging.MoreDataAvailable
	lastSeen := int64(baseSeason.Paging.LastSeen)

	for pagesRemaining {
		nextSeasonPartial, err := c.getSeasonPartial(ctx, seasonID, lastSeen)
		if err != nil {
			return &networkentities.Season{}, fmt.Errorf("failed to fetch payload for season %s offset %d: %w", baseSeason.Title, lastSeen, err)
		}
		baseSeason.Episodes = append(baseSeason.Episodes, nextSeasonPartial.Episodes...)
		pagesRemaining = nextSeasonPartial.Paging.MoreDataAvailable
		lastSeen = int64(nextSeasonPartial.Paging.LastSeen)
	}

	return baseSeason, nil
}

func (c *Client) getSeasonPartial(ctx context.Context, seasonID int64, lastSeen int64) (*networkentities.Season, error) {
	url := fmt.Sprintf("%s/v4/season/%d", c.baseURL, seasonID)
	if lastSeen > 0 {
		url += fmt.Sprintf("?lastSeen=%d", lastSeen)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Add("Realm", "dce.wwe")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", c.config.Network.XApiKey)
	req.Header.Add("x-app-var", c.config.Network.XAppVar)
	req.Header.Add("User-Agent", c.config.Network.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &networkentities.Season{}, fmt.Errorf("failed to read episode info. got status code %d", resp.StatusCode)
	}

	var season networkentities.Season
	if err := json.NewDecoder(resp.Body).Decode(&season); err != nil {
		return nil, fmt.Errorf("failed to decode season: %w", err)
	}

	return &season, nil
}
