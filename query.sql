-- name: GetSeason :one
SELECT * FROM seasons
WHERE id = ? LIMIT 1;

-- name: ListSeasons :many
SELECT * FROM seasons
ORDER BY name;

-- name: CreateSeason :one
INSERT INTO seasons (
    id, 
    name
) VALUES (
    ?, ?
)
RETURNING *;

-- name: DeleteSeason :exec
DELETE FROM seasons
WHERE id = ?;

-- name: GetToken :one
SELECT * FROM tokens
WHERE id = ? LIMIT 1;

-- name: ListTokens :many
SELECT * FROM tokens
ORDER BY id;

-- name: CreateToken :one
INSERT INTO tokens (
    id,
    value
) VALUES (
    ?, ?
)
RETURNING *;

-- name: UpdateToken :exec
UPDATE tokens
set value = ?
WHERE id = ?;

-- name: DeleteToken :exec
DELETE FROM tokens
WHERE id = ?;