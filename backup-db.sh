#!/bin/bash

# ============================================
# Database Backup Script
# ============================================
# Automatically backup PostgreSQL database
# Usage: ./backup-db.sh

set -e

BACKUP_DIR="/opt/tealinux/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/tealinux_$DATE.sql"

echo "🗄️  Starting database backup..."

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Dump database
docker compose exec -T postgres pg_dump -U tealinux_user tealinux > "$BACKUP_FILE"

# Compress backup
gzip "$BACKUP_FILE"

echo "✅ Backup completed: $BACKUP_FILE.gz"

# Keep only last 7 days of backups
echo "🧹 Cleaning old backups (older than 7 days)..."
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +7 -delete

# Show backup size
BACKUP_SIZE=$(du -h "$BACKUP_FILE.gz" | cut -f1)
echo "📦 Backup size: $BACKUP_SIZE"

# List all backups
echo ""
echo "📋 Available backups:"
ls -lh "$BACKUP_DIR"/*.sql.gz 2>/dev/null || echo "No backups found"
