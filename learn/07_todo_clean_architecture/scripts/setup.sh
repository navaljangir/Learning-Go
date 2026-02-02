#!/bin/bash

# Database setup script for TODO application
set -e

echo "======================================"
echo "TODO App - Database Setup"
echo "======================================"
echo ""

# Load environment variables
if [ -f .env ]; then
    echo "Loading environment variables from .env file..."
    source .env
else
    echo "Error: .env file not found."
    echo "Please copy .env.example to .env and configure your database settings."
    exit 1
fi

# Check if required environment variables are set
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_NAME" ]; then
    echo "Error: Missing required environment variables."
    echo "Please ensure DB_HOST, DB_PORT, DB_USER, and DB_NAME are set in .env file."
    exit 1
fi

echo "Database Configuration:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  User: $DB_USER"
echo "  Database: $DB_NAME"
echo ""

# Create database if it doesn't exist
echo "Creating database if it doesn't exist..."

# Check if database exists
DB_EXISTS=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" 2>/dev/null | grep -c 1 || echo "0")

if [ "$DB_EXISTS" -eq "0" ]; then
    echo "Database '$DB_NAME' does not exist. Creating..."
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME"
    echo "✓ Database '$DB_NAME' created successfully!"
else
    echo "✓ Database '$DB_NAME' already exists."
fi

echo ""
echo "======================================"
echo "Database setup complete!"
echo "======================================"
echo ""
echo "Next steps:"
echo "  1. Run './scripts/migrate.sh up' to apply database migrations"
echo "  2. Run 'make run' to start the application"
echo ""
