CREATE TABLE accounts (
   id uuid PRIMARY KEY,
   owner VARCHAR NOT NULL,
   balance BIGINT NOT NULL,
   currency VARCHAR NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE entries (
   id bigserial PRIMARY KEY,
   account_id uuid NOT NULL REFERENCES accounts(id),
   amount BIGINT NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE transfers (
   id bigserial PRIMARY KEY,
   sender_id uuid NOT NULL REFERENCES accounts(id),
   recipient_id uuid NOT NULL REFERENCES accounts(id),
   amount BIGINT NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("sender_id");

CREATE INDEX ON "transfers" ("recipient_id");

CREATE INDEX ON "transfers" ("sender_id", "recipient_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';
