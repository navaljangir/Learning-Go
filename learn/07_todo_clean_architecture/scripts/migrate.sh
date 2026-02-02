#!/bin/bash

# Database migration script for TODO application
set -e

# Load environment variables
if [ -f .env ]; then
    source .env
else
    echo "Error: .env file not found."
    echo "Please copy .env.example to .env and configure your database settings."
    exit 1
fi

# Build database URL
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Check if migrate command is available
if ! command -v migrate &> /dev/null; then
    echo "Error: 'migrate' command not found."
    echo ""
    echo "Please install golang-migrate:"
    echo "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    echo ""
    exit 1
fi

# Run migration based on command
case "$1" in
    up)
        echo "Running migrations up..."
        migrate -path migrations -database "$DB_URL" up
        echo "✓ Migrations applied successfully!"
        ;;
    down)
        echo "Rolling back migrations..."
        if [ -z "$2" ]; then
            migrate -path migrations -database "$DB_URL" down 1
        else
            migrate -path migrations -database "$DB_URL" down "$2"
        fi
        echo "✓ Migrations rolled back successfully!"
        ;;
    force)
        if [ -z "$2" ]; then
            echo "Error: version number required for force command"
            echo "Usage: $0 force <version>"
            exit 1
        fi
        echo "Forcing migration version to $2..."
        migrate -path migrations -database "$DB_URL" force "$2"
        echo "✓ Migration version forced to $2"
        ;;
    version)
        echo "Current migration version:"
        migrate -path migrations -database "$DB_URL" version
        ;;
    create)
        if [ -z "$2" ]; then
            echo "Error: migration name required for create command"
            echo "Usage: $0 create <migration_name>"
            exit 1
        fi
        echo "Creating migration: $2"
        migrate create -ext sql -dir migrations -seq "$2"
        echo "✓ Migration files created successfully!"
        ;;
    *)
        echo "Usage: $0 {up|down [steps]|force <version>|version|create <name>}"
        echo ""
        echo "Commands:"
        echo "  up              Apply all pending migrations"
        echo "  down [steps]    Rollback migrations (default: 1 step)"
        echo "  force <version> Force database to a specific version"
        echo "  version         Show current migration version"
        echo "  create <name>   Create a new migration file"
        echo ""
        exit 1
        ;;
esac
