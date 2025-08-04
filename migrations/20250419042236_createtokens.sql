-- +goose Up
-- +goose StatementBegin
CREATE TABLE tokens (
    id VARCHAR(50) PRIMARY KEY,
    value TEXT NOT NULL,
    expires_at INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tokens;
-- +goose StatementEnd
