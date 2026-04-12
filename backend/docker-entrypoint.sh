#!/bin/sh

echo "⏳ Waiting for Postgres..."

until nc -z postgres 5432; do
  sleep 1
done

echo "✅ Postgres is up"

echo "🚀 Running migrations..."

migrate -path ./migrations \
  -database "$TASKFLOW_DATABASE_ADDR" up

echo "🔥 Starting API..."

./app