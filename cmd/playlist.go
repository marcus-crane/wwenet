package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/marcus-crane/wwenet/config"
	"github.com/marcus-crane/wwenet/storage"
)

func CachePlaylist(ctx context.Context, cmd *cli.Command, cfg config.Config, db *storage.Queries) error {
	return cachePlaylist(ctx, int64(cmd.Int("id")), cfg, db)
}

func cachePlaylist(ctx context.Context, playlistID int64, cfg config.Config, db *storage.Queries) error {
	return nil
}
