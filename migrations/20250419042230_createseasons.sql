-- +goose Up
-- +goose StatementBegin
CREATE TABLE seasons (
    id INTEGER PRIMARY KEY,
    name text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE seasons;
-- +goose StatementEnd
