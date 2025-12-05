#!/bin/bash

# Load environment variables
source .env

# Migration directory
MIGRATION_DIR="./migrations"

# Database connection string
DB_URL="mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"

echo "Running migrations..."

# Check if migrate tool is installed
if ! command -v migrate &> /dev/null; then
    echo "migrate tool not found. Installing..."
    go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi

# Run migrations
migrate -path ${MIGRATION_DIR} -database "${DB_URL}" up

if [ $? -eq 0 ]; then
    echo "Migrations completed successfully!"
else
    echo "Migration failed!"
    exit 1
fi