ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";

ALTER TABLE "accounts" ADD COLUMN owner_tempr VARCHAR;

UPDATE "accounts" SET owner_tempr = owner::varchar;

ALTER TABLE "accounts" DROP COLUMN "owner";

ALTER TABLE "accounts" RENAME COLUMN owner_tempr TO owner;

ALTER TABLE "accounts" ALTER COLUMN updated_at SET DEFAULT now();

DROP TABLE users;