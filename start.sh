#!/bin/sh
set -e

DBHOST=postgres
DBUSER=root
DBPASSWORD=secret
DBNAME=dummy_bank
DBSSL=disable

DBSTRING="host=$DBHOST user=$DBUSER password=$DBPASSWORD dbname=$DBNAME sslmode=$DBSSL"
echo "run db migrations"
goose -dir /app/sql/schemas postgres "$DBSTRING" up

echo "start the app"
exec "$@"