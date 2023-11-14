-- +goose Up
CREATE TABLE accounts (
   id uuid PRIMARY KEY,
   owner VARCHAR NOT NULL,
   balance BIGINT NOT NULL,
   currency VARCHAR NOT NULL,
   created_at TIMESTAMP NOT NULL 
);

-- +goose Down
DROP TABLE accounts;