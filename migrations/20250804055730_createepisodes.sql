-- +goose Up
-- +goose StatementBegin
CREATE TABLE episodes (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    cover_url TEXT,
    thumbnail_url TEXT,
    poster_url TEXT,
    duration INT,
    external_asset_id TEXT,
    rating TEXT,
    descriptors TEXT,
    season_number INT,
    episode_number INT,
    season_id INTEGER REFERENCES seasons(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE episodes;
-- +goose StatementEnd
