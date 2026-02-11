#!/bin/bash

# ============================================
# Health Check Script
# ============================================
# Check the health of all services
# Usage: ./health-check.sh

set -e

echo "🏥 TeaLinuxOS Health Check"
echo "=========================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check Docker
echo "🐳 Docker Status:"
if docker info > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Docker is running${NC}"
else
    echo -e "${RED}❌ Docker is not running${NC}"
    exit 1
fi
echo ""

# Check containers
echo "📦 Container Status:"
docker compose ps
echo ""

# Check frontend
echo "🎨 Frontend (Astro SSR):"
if docker compose exec -T frontend wget -q -O- http://localhost:4321 > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Frontend is healthy${NC}"
else
    echo -e "${RED}❌ Frontend is not responding${NC}"
fi

# Check backend
echo "⚙️  Backend (Go API):"
if docker compose exec -T backend wget -q -O- http://localhost:3000 > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Backend is healthy${NC}"
else
    echo -e "${RED}❌ Backend is not responding${NC}"
fi

# Check database
echo "🗄️  Database (PostgreSQL):"
if docker compose exec -T postgres pg_isready -U tealinux_user > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Database is healthy${NC}"
else
    echo -e "${RED}❌ Database is not responding${NC}"
fi

# Check nginx
echo "🌐 Nginx:"
if docker compose exec -T nginx nginx -t > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Nginx configuration is valid${NC}"
else
    echo -e "${RED}❌ Nginx configuration has errors${NC}"
fi
echo ""

# Check SSL certificate
echo "🔒 SSL Certificate:"
if docker compose exec certbot certbot certificates 2>/dev/null | grep -q "VALID"; then
    echo -e "${GREEN}✅ SSL certificate is valid${NC}"
    docker compose exec certbot certbot certificates 2>/dev/null | grep "Expiry Date"
else
    echo -e "${YELLOW}⚠️  SSL certificate not found or expired${NC}"
fi
echo ""

# Check disk space
echo "💾 Disk Usage:"
df -h / | tail -1 | awk '{print "Used: " $3 " / " $2 " (" $5 ")"}'
echo ""

# Check Docker disk usage
echo "🐋 Docker Disk Usage:"
docker system df
echo ""

# Check memory
echo "🧠 Memory Usage:"
free -h | grep Mem | awk '{print "Used: " $3 " / " $2}'
echo ""

# Check recent logs for errors
echo "📋 Recent Errors (last 10):"
docker compose logs --tail=100 2>&1 | grep -i "error" | tail -10 || echo "No recent errors found"
echo ""

echo "=========================="
echo "✅ Health check completed!"
