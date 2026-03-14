#!/bin/bash

# Exit on error
set -e

# Ensure Go binaries are in PATH
export PATH=$PATH:/go/bin

echo "Running network diagnostics..."
nslookup db || echo "nslookup failed"
nc -zv db 5432 || echo "nc failed"

echo "Waiting for database (${DATABASE_HOST}:${DATABASE_PORT}) to be ready..."
# Use environment variables for more flexible connection
until goose -dir services/orders/_migrations postgres "host=${DATABASE_HOST} port=${DATABASE_PORT} user=${DATABASE_USERNAME} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} sslmode=disable" status; do
  echo "Database is unavailable - sleeping"
  sleep 2
done

echo "Database is up - starting migrations..."
# Run migrations for all services
make migrate-orders ARGS=up
make migrate-kitchen ARGS=up

echo "Starting all services..."
# Run the main application
make run-all
