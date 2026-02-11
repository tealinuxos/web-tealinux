#!/bin/bash

# ============================================
# SSL Certificate Setup Script
# ============================================
# This script obtains SSL certificates using Certbot
# Run this ONCE on your VPS after initial deployment

set -e

DOMAIN="${DOMAIN:-tealinuxos.org}"
EMAIL="${EMAIL:-admin@tealinuxos.org}"

echo "🔐 Setting up SSL certificates for: $DOMAIN"
echo "📧 Email: $EMAIL"
echo ""

# Create directories
mkdir -p certbot/conf
mkdir -p certbot/www

# Stop nginx temporarily
docker compose stop nginx

# Obtain certificate
docker compose run --rm certbot certonly \
  --webroot \
  --webroot-path=/var/www/certbot \
  --email "$EMAIL" \
  --agree-tos \
  --no-eff-email \
  -d "$DOMAIN" \
  -d "www.$DOMAIN"

# Start nginx again
docker compose up -d nginx

echo ""
echo "✅ SSL certificates obtained successfully!"
echo "🔄 Certificates will auto-renew every 12 hours via certbot container"
