-- +goose Up
CREATE TABLE entries (
   id uuid PRIMARY KEY,
   account_id uuid NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
   amount BIGINT NOT NULL,
   created_at TIMESTAMP NOT NULL 
);

-- +goose Down
DROP TABLE entries;