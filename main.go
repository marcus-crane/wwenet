package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v3"

	"github.com/marcus-crane/wwenet/cmd"
	"github.com/marcus-crane/wwenet/config"
	"github.com/marcus-crane/wwenet/migrations"
	"github.com/marcus-crane/wwenet/storage"
)

var (
	cfg   config.Config
	store *storage.Queries
)

func main() {
	wwecmd := &cli.Command{
		Name:  "wwenet",
		Usage: "enjoy your favourite wwe network matches while offline",
		Before: func(ctx context.Context, ucmd *cli.Command) (context.Context, error) {
			if err := loadConfiguration(); err != nil {
				return ctx, err
			}

			return ctx, runMigrations()
		},
		Commands: []*cli.Command{
			{
				Name:  "cache",
				Usage: "cache metadata for series, seasons and episodes",
				Commands: []*cli.Command{
					{
						Name:  "playlist",
						Usage: "cache metadata for all seasons and episodes that make up a playlist",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "id",
								Usage:    "playlist id to cache",
								Required: true,
							},
						},
						Action: func(ctx context.Context, ucmd *cli.Command) error {
							return cmd.CachePlaylist(ctx, ucmd, cfg, store)
						},
					},
					{
						Name:  "series",
						Usage: "cache metadata for all seasons and episodes that make up a series",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "id",
								Usage:    "series id to cache",
								Required: true,
							},
						},
						Action: func(ctx context.Context, ucmd *cli.Command) error {
							return cmd.CacheSeries(ctx, ucmd, cfg, store)
						},
					},
					{
						Name:  "season",
						Usage: "cache metadata for all seasons and episodes that make up a season",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "id",
								Usage:    "season id to cache",
								Required: true,
							},
						},
						Action: func(ctx context.Context, ucmd *cli.Command) error {
							return cmd.CacheSeason(ctx, ucmd, cfg, store)
						},
					},
					{
						Name:  "episode",
						Usage: "cache metadata for all seasons and episodes that make up a episode",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "id",
								Usage:    "episode id to cache",
								Required: true,
							},
						},
						Action: func(ctx context.Context, ucmd *cli.Command) error {
							return cmd.CacheEpisode(ctx, ucmd, cfg, store)
						},
					},
				},
			},
			{
				Name:  "download",
				Usage: "download video content for offline usage",
				Commands: []*cli.Command{
					{
						Name:  "episode",
						Usage: "download a single episode",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "id",
								Usage:    "episode id to download",
								Required: true,
							},
							&cli.StringFlag{
								Name:  "quality",
								Usage: "video quality (1080p, 720p, 480p, 360p, 240p)",
								Value: "1080p",
							},
						},
						Action: func(ctx context.Context, ucmd *cli.Command) error {
							return cmd.DownloadEpisode(ctx, ucmd, cfg, store)
						},
					},
					{
						Name:  "season",
						Usage: "download all episodes in a season",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "id",
								Usage:    "season id to download",
								Required: true,
							},
							&cli.StringFlag{
								Name:  "quality",
								Usage: "video quality (1080p, 720p, 480p, 360p, 240p)",
								Value: "1080p",
							},
						},
						Action: func(ctx context.Context, ucmd *cli.Command) error {
							return cmd.DownloadSeason(ctx, ucmd, cfg, store)
						},
					},
					{
						Name:  "series",
						Usage: "download all episodes in a series",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "id",
								Usage:    "series id to download",
								Required: true,
							},
							&cli.StringFlag{
								Name:  "quality",
								Usage: "video quality (1080p, 720p, 480p, 360p, 240p)",
								Value: "1080p",
							},
						},
						Action: func(ctx context.Context, ucmd *cli.Command) error {
							return cmd.DownloadSeries(ctx, ucmd, cfg, store)
						},
					},
					{
						Name:  "playlist",
						Usage: "download all episodes in a playlist",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:     "id",
								Usage:    "playlist id to download",
								Required: true,
							},
							&cli.StringFlag{
								Name:  "quality",
								Usage: "video quality (1080p, 720p, 480p, 360p, 240p)",
								Value: "1080p",
							},
						},
						Action: func(ctx context.Context, ucmd *cli.Command) error {
							return cmd.DownloadPlaylist(ctx, ucmd, cfg, store)
						},
					},
				},
			},
			{
				Name:  "config",
				Usage: "output config",
				Action: func(ctx context.Context, ucmd *cli.Command) error {
					return cmd.OutputConfig(ctx, ucmd, cfg, store)
				},
			},
		},
	}

	if err := wwecmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func loadConfiguration() error {
	k := koanf.New(".")
	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		return err
	}

	var loadedCfg config.Config
	if err := k.Unmarshal("", &loadedCfg); err != nil {
		return err
	}
	cfg = loadedCfg
	return nil
}

func runMigrations() error {
	db, err := sql.Open("sqlite", "wwenet.db")
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrations.GetMigrations())

	if err := goose.SetDialect(string(goose.DialectSQLite3)); err != nil {
		return err
	}

	goose.SetLogger(goose.NopLogger())

	if err := goose.Up(db, "."); err != nil {
		return err
	}

	store = storage.New(db)

	return nil
}
