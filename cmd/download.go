package cmd

import (
	"context"

	"github.com/marcus-crane/wwenet/api"
	"github.com/marcus-crane/wwenet/config"
	"github.com/marcus-crane/wwenet/downloader"
	"github.com/marcus-crane/wwenet/login"
	"github.com/marcus-crane/wwenet/storage"
	"github.com/urfave/cli/v3"
)

func DownloadEpisode(ctx context.Context, cmd *cli.Command, cfg config.Config, db *storage.Queries) error {
	token, err := login.GetAuthToken(ctx, cfg, db)
	if err != nil {
		return err
	}

	client := api.NewClient(token, cfg)
	dl := downloader.New(client, cfg, db)

	return dl.DownloadEpisode(ctx, int64(cmd.Int("id")))
}
