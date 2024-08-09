#!/bin/sh

until /usr/local/bin/migrate -path=/app/migrations -database=${DB_DSN} up; do
  echo "Migration failed..."
done
