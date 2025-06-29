#!/bin/sh

set -e

echo "run db migrate"

if [ -f /app/app.env ]; then
  echo "Loading /app/app.env"
  source /app/app.env
fi

/app/migrate --path /app/migration --database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
