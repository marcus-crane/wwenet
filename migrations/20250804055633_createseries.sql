-- +goose Up
-- +goose StatementBegin
CREATE TABLE series (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    long_description TEXT,
    small_cover_url TEXT,
    cover_url TEXT,
    title_url TEXT,
    poster_url TEXT,
    logo_url TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE series;
-- +goose StatementEnd
