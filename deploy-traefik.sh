#!/bin/bash
# ============================================
# deploy-traefik.sh
# Script untuk deploy TeaLinux Web dengan Traefik
# ============================================

set -e

echo "🚀 Starting TeaLinux Deployment (Traefik Mode)..."

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Check .env file exists
if [ ! -f .env ]; then
    echo -e "${RED}❌ Error: .env file not found!${NC}"
    echo -e "${YELLOW}📋 Copy .env.traefik to .env and fill in the values:${NC}"
    echo "   cp .env.traefik .env"
    echo "   nano .env"
    exit 1
fi

# Check if proxy network exists
if ! docker network ls | grep -q "proxy"; then
    echo -e "${YELLOW}⚠️  Creating 'proxy' network...${NC}"
    docker network create proxy
    echo -e "${GREEN}✅ 'proxy' network created${NC}"
else
    echo -e "${GREEN}✅ 'proxy' network already exists${NC}"
fi

# Pull latest code (if using git)
if [ -d ".git" ]; then
    echo -e "${YELLOW}📥 Pulling latest code...${NC}"
    git pull origin main
fi

# Build and deploy using Traefik compose file
echo -e "${YELLOW}🔨 Building Docker images...${NC}"
docker compose -f docker-compose.traefik.yml build --no-cache

echo -e "${YELLOW}🛑 Stopping existing containers...${NC}"
docker compose -f docker-compose.traefik.yml down

echo -e "${YELLOW}🚀 Starting services...${NC}"
docker compose -f docker-compose.traefik.yml up -d

# Wait for services to be ready
echo -e "${YELLOW}⏳ Waiting for services to start...${NC}"
sleep 10

# Health check
echo -e "${YELLOW}🏥 Running health checks...${NC}"
echo ""

# Check postgres
if docker exec tealinux-postgres pg_isready -U tealinux_user -d tealinux > /dev/null 2>&1; then
    echo -e "${GREEN}  ✅ PostgreSQL is healthy${NC}"
else
    echo -e "${RED}  ❌ PostgreSQL is NOT healthy${NC}"
fi

# Check backend
if docker ps --format '{{.Names}}' | grep -q tealinux-backend; then
    echo -e "${GREEN}  ✅ Backend is running${NC}"
else
    echo -e "${RED}  ❌ Backend is NOT running${NC}"
fi

# Check frontend
if docker ps --format '{{.Names}}' | grep -q tealinux-frontend; then
    echo -e "${GREEN}  ✅ Frontend is running${NC}"
else
    echo -e "${RED}  ❌ Frontend is NOT running${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Deployment complete!${NC}"
echo -e "📊 Running containers:"
docker compose -f docker-compose.traefik.yml ps

echo ""
echo -e "${YELLOW}📝 View logs:${NC}"
echo "   docker compose -f docker-compose.traefik.yml logs -f"
echo "   docker compose -f docker-compose.traefik.yml logs -f backend"
echo "   docker compose -f docker-compose.traefik.yml logs -f frontend"
