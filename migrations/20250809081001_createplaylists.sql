-- +goose Up
-- +goose StatementBegin
CREATE TABLE playlists (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    small_cover_url TEXT,
    cover_url TEXT,
    playlist_type TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE playlists;
-- +goose StatementEnd
