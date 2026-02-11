#!/bin/bash

# ============================================
# Quick Deploy Script
# ============================================
# Quick deployment script for updates
# Usage: ./deploy.sh

set -e

echo "🚀 TeaLinuxOS Quick Deploy"
echo "=========================="
echo ""

# Pull latest code
echo "📥 Pulling latest code from Git..."
git pull origin main

# Pull latest Docker images
echo "🐳 Pulling latest Docker images..."
docker compose pull

# Stop services
echo "🛑 Stopping services..."
docker compose down

# Start services
echo "▶️  Starting services..."
docker compose up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Health check
echo ""
echo "🏥 Running health check..."
./health-check.sh

echo ""
echo "=========================="
echo "✅ Deployment completed!"
echo ""
echo "🌐 Website: https://tealinuxos.org"
echo "📊 View logs: docker compose logs -f"
