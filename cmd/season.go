package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/marcus-crane/wwenet/api"
	"github.com/marcus-crane/wwenet/config"
	"github.com/marcus-crane/wwenet/login"
	"github.com/marcus-crane/wwenet/storage"
)

func CacheSeason(ctx context.Context, cmd *cli.Command, cfg config.Config, db *storage.Queries) error {
	return cacheSeason(ctx, int64(cmd.Int("id")), cfg, db)
}

func cacheSeason(ctx context.Context, seasonID int64, cfg config.Config, db *storage.Queries) error {
	if existingSeason, err := db.GetSeason(ctx, seasonID); err == nil {
		fmt.Printf("Season %d (%d) is already cached\n", existingSeason.SeasonNumber.Int64, existingSeason.ID)
		return nil
	}

	token, err := login.GetAuthToken(ctx, cfg, db)
	if err != nil {
		return fmt.Errorf("failed to get auth token: %w", err)
	}

	client := api.NewClient(token, cfg)
	season, err := client.GetSeason(ctx, seasonID)
	if err != nil {
		return fmt.Errorf("failed to fetch season %d: %w", seasonID, err)
	}

	params := storage.CreateSeasonParams{
		ID:              int64(season.Id),
		Title:           season.Title,
		Description:     sqlNullString(season.Description),
		LongDescription: sqlNullString(season.LongDescription),
		SmallCoverUrl:   sqlNullString(season.SmallCoverUrl),
		CoverUrl:        sqlNullString(season.CoverUrl),
		TitleUrl:        sqlNullString(season.TitleUrl),
		PosterUrl:       sqlNullString(season.PosterUrl),
		SeasonNumber:    sqlNullInt64(int64(season.SeasonNumber)),
		EpisodeCount:    sqlNullInt64(int64(season.EpisodeCount)),
		SeriesID:        sqlNullInt64(int64(season.Series.SeriesId)),
	}

	_, err = db.CreateSeason(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to cache season: %w", err)
	}

	for _, ep := range season.Episodes {
		if err := cacheEpisode(ctx, int64(ep.Id), cfg, db); err != nil {
			fmt.Printf("failed to cache S%dE%d\n", ep.EpisodeInformation.SeasonNumber, ep.EpisodeInformation.EpisodeNumber)
		}
	}

	fmt.Printf("Cached %s Season %d\n", season.Title, season.SeasonNumber)

	return nil
}
