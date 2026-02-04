#!/bin/bash

# Database setup script for TODO application (MySQL in Docker)
set -e

echo "======================================"
echo "TODO App - Database Setup (MySQL)"
echo "======================================"
echo ""

echo "Checking for MySQL Docker container..."

# Check if container exists
if docker ps -a --format '{{.Names}}' | grep -q "^todo_mysql$"; then
    echo "✓ MySQL container 'todo_mysql' found"
    
    # Check if it's running
    if docker ps --format '{{.Names}}' | grep -q "^todo_mysql$"; then
        echo "✓ MySQL container is running"
    else
        echo "Starting MySQL container..."
        docker start todo_mysql
        echo "✓ MySQL container started"
    fi
else
    echo "MySQL container not found. Please run:"
    echo ""
    echo "  docker run -d --name todo_mysql -e MYSQL_ROOT_PASSWORD=rootpassword -e MYSQL_DATABASE=todo_db -p 3306:3306 mysql:8.0"
    echo ""
    exit 1
fi

echo ""
echo "======================================"
echo "Database setup complete!"
echo "======================================"
echo ""
echo "Next steps:"
echo "  1. Run 'make migrate-up' to apply database migrations"
echo "  2. Run 'make run' to start the application"
echo ""
