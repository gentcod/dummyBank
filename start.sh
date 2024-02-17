#!/bin/sh
set -e

echo "run db migrations"
source /app/app.env
goose -dir /app/sql/schemas postgres "$DB_URL" up

echo "start the app"
exec "$@"