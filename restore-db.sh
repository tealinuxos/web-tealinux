#!/bin/bash

# ============================================
# Database Restore Script
# ============================================
# Restore PostgreSQL database from backup
# Usage: ./restore-db.sh <backup_file.sql.gz>

set -e

if [ -z "$1" ]; then
    echo "❌ Error: No backup file specified"
    echo "Usage: ./restore-db.sh <backup_file.sql.gz>"
    echo ""
    echo "Available backups:"
    ls -lh /opt/tealinux/backups/*.sql.gz 2>/dev/null || echo "No backups found"
    exit 1
fi

BACKUP_FILE="$1"

if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Error: Backup file not found: $BACKUP_FILE"
    exit 1
fi

echo "⚠️  WARNING: This will REPLACE the current database!"
echo "Backup file: $BACKUP_FILE"
echo ""
read -p "Are you sure you want to continue? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "❌ Restore cancelled"
    exit 0
fi

echo "🗄️  Starting database restore..."

# Decompress if needed
if [[ "$BACKUP_FILE" == *.gz ]]; then
    echo "📦 Decompressing backup..."
    gunzip -c "$BACKUP_FILE" > /tmp/restore.sql
    SQL_FILE="/tmp/restore.sql"
else
    SQL_FILE="$BACKUP_FILE"
fi

# Drop existing connections
echo "🔌 Closing existing database connections..."
docker compose exec -T postgres psql -U tealinux_user -d postgres -c \
    "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'tealinux' AND pid <> pg_backend_pid();"

# Drop and recreate database
echo "🗑️  Dropping existing database..."
docker compose exec -T postgres psql -U tealinux_user -d postgres -c "DROP DATABASE IF EXISTS tealinux;"
docker compose exec -T postgres psql -U tealinux_user -d postgres -c "CREATE DATABASE tealinux;"

# Restore backup
echo "📥 Restoring database..."
docker compose exec -T postgres psql -U tealinux_user -d tealinux < "$SQL_FILE"

# Cleanup
if [ -f "/tmp/restore.sql" ]; then
    rm /tmp/restore.sql
fi

echo "✅ Database restored successfully!"
echo "🔄 Restarting backend service..."
docker compose restart backend

echo "✅ Restore completed!"
