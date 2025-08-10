-- +goose Up
-- +goose StatementBegin
CREATE TABLE playlist_episodes (
    playlist_id INTEGER REFERENCES playlists(id),
    episode_id INTEGER REFERENCES episodes(id),
    PRIMARY KEY (playlist_id, episode_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE playlist_episodes;
-- +goose StatementEnd
