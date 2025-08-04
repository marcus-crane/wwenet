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

func CacheSeries(ctx context.Context, cmd *cli.Command, cfg config.Config, db *storage.Queries) error {
	seriesID := int64(cmd.Int("id"))

	return cacheSeries(ctx, seriesID, cfg, db)
}

func cacheSeries(ctx context.Context, seriesID int64, cfg config.Config, db *storage.Queries) error {
	if existingSeason, err := db.GetSeries(ctx, seriesID); err == nil {
		fmt.Printf("Series %s is already cached\n", existingSeason.Title)
		return nil
	}

	token, err := login.GetAuthToken(ctx, cfg, db)
	if err != nil {
		return fmt.Errorf("failed to get auth token: %w", err)
	}

	client := api.NewClient(token, cfg)
	series, err := client.GetSeries(ctx, seriesID)
	if err != nil {
		return fmt.Errorf("failed to fetch series %s: %w", series.Title, err)
	}

	params := storage.CreateSeriesParams{
		ID:              int64(series.Id),
		Title:           series.Title,
		Description:     sqlNullString(series.Description),
		LongDescription: sqlNullString(series.LongDescription),
		SmallCoverUrl:   sqlNullString(series.SmallCoverUrl),
		CoverUrl:        sqlNullString(series.CoverUrl),
		TitleUrl:        sqlNullString(series.TitleUrl),
		PosterUrl:       sqlNullString(series.PosterUrl),
		LogoUrl:         sqlNullString(series.LogoUrl),
	}

	_, err = db.CreateSeries(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to cache season: %w", err)
	}

	for _, season := range series.Seasons {
		if err := cacheSeason(ctx, int64(season.Id), cfg, db); err != nil {
			fmt.Printf("failed to cache Season %s\n", season.Title)
		}
	}

	fmt.Printf("Cached Series %s\n", series.Title)

	return nil
}
