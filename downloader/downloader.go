package downloader

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/marcus-crane/wwenet/api"
	"github.com/marcus-crane/wwenet/config"
	"github.com/marcus-crane/wwenet/networkentities"
	"github.com/marcus-crane/wwenet/storage"
)

const (
	QUALITY_720BEST  int = 0
	QUALITY_1080BEST int = 1
	QUALITY_480      int = 4
	QUALITY_360      int = 5
	QUALITY_240      int = 6
)

type DownloadOptions struct {
	Quality string // "best", "worst", "720p", "1080p"
}

type Downloader struct {
	client *api.Client
	config config.Config
	db     *storage.Queries
}

func New(client *api.Client, config config.Config, db *storage.Queries) *Downloader {
	return &Downloader{
		client: client,
		config: config,
		db:     db,
	}
}

func (d *Downloader) getEpisodeDetails(ctx context.Context, episodeID int64) (*networkentities.Episode, error) {
	// We have the episode details cached but we need to fetch them fresh
	// in order to get a valid PlayerUrlCallback field
	episode, err := d.client.GetEpisode(ctx, episodeID)
	if err != nil {
		return nil, err
	}
	return episode, nil
}

func (d *Downloader) getPlaybackURL(ctx context.Context, callbackURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", callbackURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", d.config.Network.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get playback URL, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var playbackManifest networkentities.Playback

	if err := json.Unmarshal(body, &playbackManifest); err != nil {
		return "", fmt.Errorf("failed to deserialize playback manifest: %w", err)
	}

	return playbackManifest.HLS[0].Url, nil
}

func (d *Downloader) generateOutputPath(episode *networkentities.Episode) string {
	sanitizedTitle := sanitizeFilename(episode.Title)

	filename := fmt.Sprintf("S%02dE%02d - %s.mp4",
		episode.EpisodeInformation.SeasonNumber,
		episode.EpisodeInformation.EpisodeNumber,
		sanitizedTitle)

	return filepath.Join(d.config.Download.StorageDirectory, filename)
}

func (d *Downloader) downloadWithFFmpeg(ctx context.Context, m3u8URL, outputPath string, opts DownloadOptions, episode *networkentities.Episode) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH. Please install ffmpeg")
	}

	args := []string{
		"-loglevel", "error",
		"-nostats",
		"-progress", "pipe:2",
		"-hide_banner",

		// Network reliability
		"-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "5",
		"-http_persistent", "1",
		"-multiple_requests", "1",
		"-threads", "0",

		"-i", m3u8URL,

		// Re-encode to avoid occasional stream issues
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "23",

		"-c:a", "copy",

		"-bsf:a", "aac_adtstoasc",

		// Timestamp corrections
		"-avoid_negative_ts", "make_zero",
		"-fflags", "+genpts",
		outputPath,
	}

	if d.config.Network.UserAgent != "" {
		args = append([]string{"-user_agent", d.config.Network.UserAgent}, args...)
	}

	// "best" and "worst" are handled automatically by ffmpeg
	switch opts.Quality {
	case "1080p":
		args = append(args, "-map", fmt.Sprintf("0:v:%d", QUALITY_1080BEST))
	case "720p":
		args = append(args, "-map", fmt.Sprintf("0:v:%d", QUALITY_720BEST))
	case "480p":
		args = append(args, "-map", fmt.Sprintf("0:v:%d", QUALITY_480))
	case "360p":
		args = append(args, "-map", fmt.Sprintf("0:v:%d", QUALITY_360))
	case "240p":
		args = append(args, "-map", fmt.Sprintf("0:v:%d", QUALITY_240))
	default:
		// Default to the best quality
		args = append(args, "-map", fmt.Sprintf("0:v:%d", QUALITY_1080BEST))
	}

	args = append(args, "-map", "0:a") // Map all audio streams

	args = append(args, "-y", outputPath)

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	cmd.Stdout = os.Stdout
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	bar := progressbar.NewOptions(
		int(episode.Duration),
		progressbar.OptionSetDescription(fmt.Sprintf("S%02dE%02d",
			episode.EpisodeInformation.SeasonNumber,
			episode.EpisodeInformation.EpisodeNumber)),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowBytes(false),
		progressbar.OptionOnCompletion(func() {
			fmt.Printf("\nDownload completed: %s\n", filepath.Base(outputPath))
		}),
		progressbar.OptionShowCount(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionThrottle(time.Millisecond*100),
		progressbar.OptionSetPredictTime(true),
	)

	if err := cmd.Start(); err != nil {
		os.Remove(outputPath)
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	progressDone := make(chan error, 1)
	go func() {
		progressDone <- d.parseProgressWithBar(stderr, episode.Duration, bar)
	}()

	cmdErr := cmd.Wait()

	<-progressDone

	bar.Finish()

	if cmdErr != nil {
		os.Remove(outputPath)
		return fmt.Errorf("ffmpeg failed: %w", cmdErr)
	}

	return nil
}

func (d *Downloader) parseProgressWithBar(reader io.Reader, totalDurationSec int, bar *progressbar.ProgressBar) error {
	scanner := bufio.NewScanner(reader)

	totalDurationMicros := int64(totalDurationSec) * 1000000

	lastSeconds := int64(-1)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "out_time_us=") {
			timeStr := strings.TrimPrefix(line, "out_time_us=")
			currentMicros, err := strconv.ParseInt(timeStr, 10, 64)
			if err != nil {
				// Can occur if slow init ie; time=N/A
				continue
			}

			if totalDurationMicros > 0 && currentMicros > 0 {
				currentSeconds := currentMicros / 1000000

				if currentSeconds > lastSeconds {
					bar.Set(int(currentSeconds))
					lastSeconds = currentSeconds
				}
			}
		}
	}

	return scanner.Err()
}

func sanitizeFilename(filename string) string {
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "-",
	)

	return strings.TrimSpace(replacer.Replace(filename))
}

func (d *Downloader) DownloadEpisode(ctx context.Context, episodeID int64, opts DownloadOptions) error {
	if download, err := d.db.GetDownload(ctx, sql.NullInt64{Int64: episodeID, Valid: true}); err == nil {
		if _, err := os.Stat(download.FilePath.String); err == nil {
			fmt.Printf("Episode %d already downloaded at %s\n", episodeID, download.FilePath.String)
			return nil
		}
		d.db.DeleteDownload(ctx, sql.NullInt64{Int64: episodeID, Valid: true})
	}

	episode, err := d.getEpisodeDetails(ctx, episodeID)
	if err != nil {
		return fmt.Errorf("failed to get episode details: %w", err)
	}

	fmt.Printf("Starting download: %s (S%dE%d). It will take a minute to progress...\n",
		episode.Title,
		episode.EpisodeInformation.SeasonNumber,
		episode.EpisodeInformation.EpisodeNumber,
	)

	playbackUrl, err := d.getPlaybackURL(ctx, episode.PlayerUrlCallback)
	if err != nil {
		return fmt.Errorf("failed to get playback URL: %w", err)
	}

	outputPath := d.generateOutputPath(episode)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := d.downloadWithFFmpeg(ctx, playbackUrl, outputPath, opts, episode); err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}

	params := storage.CreateDownloadParams{
		EpisodeID:    sql.NullInt64{Int64: episodeID, Valid: true},
		FilePath:     sql.NullString{String: outputPath, Valid: true},
		DownloadedAt: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
	}

	if _, err := d.db.CreateDownload(ctx, params); err != nil {
		return fmt.Errorf("failed to record download: %w", err)
	}

	fmt.Printf("Successfully downloaded: %s\n", outputPath)
	return nil
}

func (d *Downloader) DownloadSeason(ctx context.Context, seasonID int64, opts DownloadOptions) error {
	season, err := d.db.GetSeason(ctx, seasonID)
	if err != nil {
		return fmt.Errorf("season %d not found in cache. Run 'wwenet cache season --id %d' first", seasonID, seasonID)
	}

	fmt.Printf("Downloading season: %s\n", season.Title)

	episodes, err := d.db.ListEpisodes(ctx)
	if err != nil {
		return fmt.Errorf("failed to list episodes: %w", err)
	}

	var seasonEpisodes []storage.Episode
	for _, ep := range episodes {
		if ep.SeasonNumber.Valid && ep.SeasonNumber.Int64 == season.SeasonNumber.Int64 {
			seasonEpisodes = append(seasonEpisodes, ep)
		}
	}

	if len(seasonEpisodes) == 0 {
		return fmt.Errorf("no episodes found for season %d. Run 'wwenet cache season --id %d' first", seasonID, seasonID)
	}

	for _, ep := range seasonEpisodes {
		if err := d.DownloadEpisode(ctx, ep.ID, opts); err != nil {
			fmt.Printf("failed to download episode %d: %v\n", ep.ID, err)
			continue
		}
	}

	return nil
}

func (d *Downloader) DownloadSeries(ctx context.Context, seriesID int64, opts DownloadOptions) error {
	series, err := d.db.GetSeries(ctx, seriesID)
	if err != nil {
		return fmt.Errorf("series %d not found in cache. Run 'wwenet cache series --id %d' first", seriesID, seriesID)
	}

	fmt.Printf("Downloading series: %s\n", series.Title)

	seasons, err := d.db.ListSeasons(ctx)
	if err != nil {
		return fmt.Errorf("failed to list seasons: %w", err)
	}

	for _, season := range seasons {
		if err := d.DownloadSeason(ctx, season.ID, opts); err != nil {
			fmt.Printf("Failed to download season %d: %v\n", season.ID, err)
			continue
		}
	}

	return nil
}
