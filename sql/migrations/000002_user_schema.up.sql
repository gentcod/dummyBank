CREATE TABLE users (
   id uuid PRIMARY KEY,
   username VARCHAR UNIQUE NOT NULL,
   harshed_password VARCHAR NOT NULL,
   full_name VARCHAR NOT NULL,
   email VARCHAR UNIQUE NOT NULL,
   password_changed_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
   created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Change owner colum from type VARCHAR to UUID:
-- Step 1: Add a new column of type UUID
ALTER TABLE "accounts" ADD COLUMN owner_temp uuid NOT NULL;

-- Step 2: Update the new column with the UUID representation of the existing varchar values
UPDATE "accounts" SET owner_temp = owner::uuid;

-- Step 3: Drop the old varchar column
ALTER TABLE "accounts" DROP COLUMN "owner";

-- Step 4: Rename the new UUID column to match the original column name
ALTER TABLE "accounts" RENAME COLUMN owner_temp TO "owner";

-- Other changes:  
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES users(id);
   
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");

ALTER TABLE "accounts" ALTER COLUMN updated_at SET DEFAULT '0001-01-01 00:00:00Z';
