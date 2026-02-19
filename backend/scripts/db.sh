#!/bin/bash

set -e

echo "=== Genvid Database Setup ==="

DB_HOST="${DATABASE_HOST:-localhost}"
DB_PORT="${DATABASE_PORT:-5432}"
DB_USER="${DATABASE_USER:-postgres}"
DB_PASSWORD="${DATABASE_PASSWORD:-postgres}"
DB_NAME="${DATABASE_NAME:-genvid}"

export PGPASSWORD="$DB_PASSWORD"

echo "Database: $DB_HOST:$DB_PORT/$DB_NAME"

run_migrations() {
    local migration_dir="./migrations"
    
    if [ -d "$migration_dir" ]; then
        echo "Running migrations..."
        
        for sql_file in "$migration_dir"/*.sql; do
            if [ -f "$sql_file" ]; then
                echo "Applying: $sql_file"
                psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$sql_file"
            fi
        done
        
        echo "Migrations completed!"
    else
        echo "No migrations directory found."
    fi
}

check_connection() {
    echo "Checking database connection..."
    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1" > /dev/null 2>&1; then
        echo "Connected successfully!"
        return 0
    else
        echo "Failed to connect. Please check your database configuration."
        return 1
    fi
}

reset_database() {
    echo "WARNING: This will delete all data!"
    read -p "Are you sure? (yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        echo "Dropping database..."
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "CREATE DATABASE $DB_NAME;"
        echo "Database reset. Running migrations..."
        run_migrations
    else
        echo "Cancelled."
    fi
}

case "${1:-}" in
    setup)
        check_connection && run_migrations
        ;;
    reset)
        reset_database
        ;;
    check)
        check_connection
        ;;
    *)
        echo "Usage: $0 {setup|reset|check}"
        echo ""
        echo "Commands:"
        echo "  setup  - Check connection and run migrations"
        echo "  reset  - Drop and recreate database (WARNING: destroys data)"
        echo "  check  - Test database connection"
        exit 1
        ;;
esac
