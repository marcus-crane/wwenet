---- SERIES ----

-- name: GetSeries :one
SELECT * FROM series
WHERE id = ? LIMIT 1;

-- name: ListSeries :many
SELECT * FROM series
ORDER BY title;

-- name: CreateSeries :one
INSERT INTO series (
    id,
    title,
    description,
    long_description,
    small_cover_url,
    cover_url,
    title_url,
    poster_url,
    logo_url
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteSeries :exec
DELETE FROM series
WHERE id = ?;

---- PLAYLIST ----

-- name: GetPlaylist :one
SELECT * FROM playlists
WHERE id = ? LIMIT 1;

-- name: ListPlaylists :many
SELECT * FROM playlists
ORDER BY title;

-- name: CreatePlaylist :one
INSERT INTO playlists (
    id, 
    title,
    description,
    small_cover_url,
    cover_url,
    playlist_type
) VALUES (
    ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeletePlaylist :exec
DELETE FROM playlists
WHERE id = ?;

---- SEASONS ----

-- name: GetSeason :one
SELECT * FROM seasons
WHERE id = ? LIMIT 1;

-- name: GetSeasonsBySeries :many
SELECT * FROM seasons
WHERE series_id = ?
ORDER BY season_number;

-- name: ListSeasons :many
SELECT * FROM seasons
ORDER BY title;

-- name: CreateSeason :one
INSERT INTO seasons (
    id, 
    title,
    description,
    long_description,
    small_cover_url,
    cover_url,
    title_url,
    poster_url,
    season_number,
    episode_count,
    series_id
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteSeason :exec
DELETE FROM seasons
WHERE id = ?;

---- EPISODES ----

-- name: GetEpisode :one
SELECT * FROM episodes
WHERE id = ? LIMIT 1;

-- name: GetEpisodesBySeason :many
SELECT * FROM episodes
WHERE season_id = ?
ORDER BY episode_number;

-- name: GetEpisodesBySeries :many
SELECT e.* FROM episodes e
JOIN seasons s ON e.season_id = s.id
WHERE s.series_id = ?
ORDER BY s.season_number, e.episode_number;

-- name: GetEpisodesByPlaylist :many
SELECT e.* FROM episodes e
JOIN playlist_episodes pe ON e.id = pe.episode_id
WHERE pe.playlist_id = ?
ORDER BY e.id;

-- name: ListEpisodes :many
SELECT * FROM episodes
ORDER BY title;

-- name: CreateEpisode :one
INSERT INTO episodes (
    id,
    title,
    description,
    cover_url,
    thumbnail_url,
    poster_url,
    duration,
    external_asset_id,
    rating,
    descriptors,
    season_number,
    episode_number,
    season_id
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: AddEpisodeToPlaylist :exec
INSERT OR IGNORE INTO playlist_episodes (playlist_id, episode_id)
VALUES (?, ?);

-- name: DeleteEpisode :exec
DELETE FROM episodes
WHERE id = ?;

--- DOWNLOADS ---

-- name: GetDownload :one
SELECT * FROM downloads
WHERE episode_id = ? LIMIT 1;

-- name: ListDownloads :many
SELECT * FROM downloads
ORDER BY episode_id;

-- name: CreateDownload :one
INSERT INTO downloads (
    episode_id,
    file_path,
    downloaded_at
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: DeleteDownload :exec
DELETE FROM downloads
WHERE episode_id = ?;

---- TOKENS ----

-- name: GetToken :one
SELECT * FROM tokens
WHERE id = ? LIMIT 1;

-- name: ListTokens :many
SELECT * FROM tokens
ORDER BY id;

-- name: CreateToken :one
INSERT INTO tokens (
    id,
    value,
    expires_at
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: UpdateToken :exec
UPDATE tokens
set value = ?, expires_at = ?
WHERE id = ?;

-- name: DeleteToken :exec
DELETE FROM tokens
WHERE id = ?;