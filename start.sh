#!/bin/sh

set -e

echo "run db migration"
/app/migration -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"