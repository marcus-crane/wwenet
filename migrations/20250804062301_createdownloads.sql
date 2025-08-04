-- +goose Up
-- +goose StatementBegin
CREATE TABLE downloads (
    episode_id INTEGER REFERENCES episodes(id),
    file_path TEXT,
    downloaded_at INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE downloads;
-- +goose StatementEnd
