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

func CacheEpisode(ctx context.Context, cmd *cli.Command, cfg config.Config, db *storage.Queries) error {
	episodeID := int64(cmd.Int("id"))

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

	fmt.Printf("Retrieved %s\n", episode.Title)
	return nil
}
