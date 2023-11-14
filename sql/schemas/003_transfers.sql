-- +goose Up
CREATE TABLE transfers (
   id uuid PRIMARY KEY,
   sender_id uuid NOT NULL REFERENCES accounts(id),
   recipient_id uuid NOT NULL REFERENCES accounts(id),
   amount BIGINT NOT NULL,
   created_at TIMESTAMP NOT NULL 
);

-- +goose Down
DROP TABLE transfers;