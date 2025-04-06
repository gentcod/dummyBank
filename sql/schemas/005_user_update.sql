-- +goose Up
ALTER TABLE "users" ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z';

-- +goose Down
ALTER TABLE "users" DROP COLUMN updated_at;
