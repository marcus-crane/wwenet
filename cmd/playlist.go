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

func CachePlaylist(ctx context.Context, cmd *cli.Command, cfg config.Config, db *storage.Queries) error {
	return cachePlaylist(ctx, int64(cmd.Int("id")), cfg, db)
}

func cachePlaylist(ctx context.Context, playlistID int64, cfg config.Config, db *storage.Queries) error {
	if existingPlaylist, err := db.GetPlaylist(ctx, playlistID); err == nil {
		fmt.Printf("Playlist %s (%d) is already cached\n", existingPlaylist.Title, existingPlaylist.ID)
	}

	token, err := login.GetAuthToken(ctx, cfg, db)
	if err != nil {
		return fmt.Errorf("failed to get auth token: %w", err)
	}

	client := api.NewClient(token, cfg)
	playlist, err := client.GetPlaylist(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to fetch playlist %d: %w", playlistID, err)
	}

	params := storage.CreatePlaylistParams{
		ID:            int64(playlist.Id),
		Title:         playlist.Title,
		Description:   sqlNullString(playlist.Description),
		SmallCoverUrl: sqlNullString(playlist.SmallCoverUrl),
		CoverUrl:      sqlNullString(playlist.CoverUrl),
		PlaylistType:  sqlNullString(playlist.PlaylistType),
	}

	_, err = db.CreatePlaylist(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to cache playlist: %w", err)
	}

	for _, ep := range playlist.VODs {
		if err := cacheEpisode(ctx, int64(ep.Id), cfg, db); err != nil {
			fmt.Printf("failed to cache episode S%dE%d\n", ep.EpisodeInformation.SeasonNumber, ep.EpisodeInformation.EpisodeNumber)
		}

		addParams := storage.AddEpisodeToPlaylistParams{
			PlaylistID: sqlNullInt64(playlistID),
			EpisodeID:  sqlNullInt64(int64(ep.Id)),
		}

		if err := db.AddEpisodeToPlaylist(ctx, addParams); err != nil {
			fmt.Printf("failed to link episode %d to playlist %d\n", ep.Id, playlistID)
		}
	}

	fmt.Printf("Cached playlist %s\n", playlist.Title)

	return nil
}
