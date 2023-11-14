-- +goose Up
CREATE TABLE accounts (
   id uuid PRIMARY KEY,
   owner VARCHAR NOT NULL,
   balance BIGINT NOT NULL,
   currency VARCHAR NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT now(),
   updated_at TIMESTAMP NOT NULL
);

CREATE TABLE entries (
   id uuid PRIMARY KEY,
   account_id uuid NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
   amount BIGINT NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE transfers (
   id uuid PRIMARY KEY,
   sender_id uuid NOT NULL REFERENCES accounts(id),
   recipient_id uuid NOT NULL REFERENCES accounts(id),
   amount BIGINT NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE transfers;
DROP TABLE entries;
DROP TABLE accounts;