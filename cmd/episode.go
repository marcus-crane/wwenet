package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/marcus-crane/wwenet/api"
	"github.com/marcus-crane/wwenet/config"
	"github.com/marcus-crane/wwenet/login"
	"github.com/marcus-crane/wwenet/storage"
)

func CacheEpisode(ctx context.Context, cmd *cli.Command, cfg config.Config, db *storage.Queries) error {
	return cacheEpisode(ctx, int64(cmd.Int("id")), cfg, db)
}

func cacheEpisode(ctx context.Context, episodeID int64, cfg config.Config, db *storage.Queries) error {
	if existingEpisode, err := db.GetEpisode(ctx, episodeID); err == nil {
		fmt.Printf("Episode %d (%s) is already cached\n", existingEpisode.ID, existingEpisode.Title)
		return nil
	}

	token, err := login.GetAuthToken(ctx, cfg, db)
	if err != nil {
		return fmt.Errorf("failed to get auth token: %w", err)
	}

	client := api.NewClient(token, cfg)
	episode, err := client.GetEpisode(ctx, episodeID)
	if err != nil {
		return fmt.Errorf("failed to fetch episode %d: %w", episodeID, err)
	}

	params := storage.CreateEpisodeParams{
		ID:              int64(episode.Id),
		Title:           episode.Title,
		Description:     sqlNullString(episode.Description),
		CoverUrl:        sqlNullString(episode.CoverUrl),
		ThumbnailUrl:    sqlNullString(episode.ThumbnailUrl),
		PosterUrl:       sqlNullString(episode.PosterUrl),
		Duration:        sqlNullInt64(int64(episode.Duration)),
		ExternalAssetID: sqlNullString(episode.ExternalAssetId),
		Rating:          sqlNullString(episode.Rating.Rating),
		Descriptors:     sqlNullString(strings.Join(episode.Rating.Descriptors, ",")),
		SeasonNumber:    sqlNullInt64(int64(episode.EpisodeInformation.SeasonNumber)),
		EpisodeNumber:   sqlNullInt64(int64(episode.EpisodeInformation.EpisodeNumber)),
	}

	cachedEpisode, err := db.CreateEpisode(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to cache episode: %w", err)
	}

	fmt.Printf("Cached %s (S%dE%d)\n", cachedEpisode.Title, cachedEpisode.SeasonNumber.Int64, cachedEpisode.EpisodeNumber.Int64)

	return nil
}

func cacheEpisodeWithSeason(ctx context.Context, episodeID, seasonID int64, cfg config.Config, db *storage.Queries) error {
	if existingEpisode, err := db.GetEpisode(ctx, episodeID); err == nil {
		fmt.Printf("Episode %d (%s) is already cached\n", existingEpisode.ID, existingEpisode.Title)
		return nil
	}

	token, err := login.GetAuthToken(ctx, cfg, db)
	if err != nil {
		return fmt.Errorf("failed to get auth token: %w", err)
	}

	client := api.NewClient(token, cfg)
	episode, err := client.GetEpisode(ctx, episodeID)
	if err != nil {
		return fmt.Errorf("failed to fetch episode %d: %w", episodeID, err)
	}

	params := storage.CreateEpisodeParams{
		ID:              int64(episode.Id),
		Title:           episode.Title,
		Description:     sqlNullString(episode.Description),
		CoverUrl:        sqlNullString(episode.CoverUrl),
		ThumbnailUrl:    sqlNullString(episode.ThumbnailUrl),
		PosterUrl:       sqlNullString(episode.PosterUrl),
		Duration:        sqlNullInt64(int64(episode.Duration)),
		ExternalAssetID: sqlNullString(episode.ExternalAssetId),
		Rating:          sqlNullString(episode.Rating.Rating),
		Descriptors:     sqlNullString(strings.Join(episode.Rating.Descriptors, ",")),
		SeasonNumber:    sqlNullInt64(int64(episode.EpisodeInformation.SeasonNumber)),
		EpisodeNumber:   sqlNullInt64(int64(episode.EpisodeInformation.EpisodeNumber)),
		SeasonID:        sqlNullInt64(seasonID),
	}

	cachedEpisode, err := db.CreateEpisode(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to cache episode: %w", err)
	}

	fmt.Printf("Cached %s (S%dE%d)\n", cachedEpisode.Title, cachedEpisode.SeasonNumber.Int64, cachedEpisode.EpisodeNumber.Int64)

	return nil
}
