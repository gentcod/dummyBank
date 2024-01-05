-- +goose Up
ALTER TABLE users ADD COLUMN username VARCHAR NOT NULL

-- +goose Down
ALTER TABLE users DROP COLUMN username
