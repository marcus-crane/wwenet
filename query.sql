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

---- SEASONS ----

-- name: GetSeason :one
SELECT * FROM seasons
WHERE id = ? LIMIT 1;

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
    episode_count
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteSeason :exec
DELETE FROM seasons
WHERE id = ?;

---- EPISODES ----

-- name: GetEpisode :one
SELECT * FROM episodes
WHERE id = ? LIMIT 1;

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
    episode_number
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

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