.PHONY: help dev dev-up dev-down prod-up prod-down build logs health backup restore ssl clean

# Default target
help:
	@echo "🍵 TeaLinuxOS Web Platform - Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev          - Start development environment (Docker with hot-reload)"
	@echo "  make dev-up       - Start dev containers in background"
	@echo "  make dev-down     - Stop dev containers"
	@echo "  make dev-logs     - View dev logs"
	@echo ""
	@echo "Production:"
	@echo "  make prod-up      - Start production containers"
	@echo "  make prod-down    - Stop production containers"
	@echo "  make prod-logs    - View production logs"
	@echo "  make deploy       - Quick deploy (pull + rebuild + restart)"
	@echo ""
	@echo "Database:"
	@echo "  make backup       - Backup database"
	@echo "  make restore      - Restore database (requires BACKUP_FILE=path)"
	@echo ""
	@echo "Utilities:"
	@echo "  make health       - Run health check"
	@echo "  make ssl          - Setup SSL certificates"
	@echo "  make logs         - View all logs"
	@echo "  make clean        - Clean Docker resources"
	@echo "  make build        - Build all Docker images"
	@echo ""

# Development commands
dev:
	@echo "🚀 Starting development environment..."
	docker compose -f docker-compose.dev.yml up

dev-up:
	@echo "🚀 Starting development containers in background..."
	docker compose -f docker-compose.dev.yml up -d
	@echo "✅ Development environment started!"
	@echo "Frontend: http://localhost:4321"
	@echo "Backend: http://localhost:3000"

dev-down:
	@echo "🛑 Stopping development containers..."
	docker compose -f docker-compose.dev.yml down

dev-logs:
	docker compose -f docker-compose.dev.yml logs -f

# Production commands
prod-up:
	@echo "🚀 Starting production containers..."
	docker compose up -d
	@echo "✅ Production environment started!"

prod-down:
	@echo "🛑 Stopping production containers..."
	docker compose down

prod-logs:
	docker compose logs -f

# Build
build:
	@echo "🔨 Building all Docker images..."
	docker compose build

# Deploy
deploy:
	@echo "🚀 Deploying application..."
	./deploy.sh

# Health check
health:
	@echo "🏥 Running health check..."
	./health-check.sh

# Database operations
backup:
	@echo "💾 Backing up database..."
	./backup-db.sh

restore:
ifndef BACKUP_FILE
	@echo "❌ Error: BACKUP_FILE not specified"
	@echo "Usage: make restore BACKUP_FILE=backups/tealinux_20260211_120000.sql.gz"
	@exit 1
endif
	@echo "📥 Restoring database from $(BACKUP_FILE)..."
	./restore-db.sh $(BACKUP_FILE)

# SSL setup
ssl:
	@echo "🔒 Setting up SSL certificates..."
	./setup-ssl.sh

# Logs
logs:
	docker compose logs -f

# Clean up
clean:
	@echo "🧹 Cleaning Docker resources..."
	docker compose down -v
	docker system prune -af
	@echo "✅ Cleanup completed!"

# Frontend specific
fe-dev:
	@echo "🎨 Starting frontend development..."
	cd tealinux-fe && npm run dev

fe-build:
	@echo "🔨 Building frontend..."
	cd tealinux-fe && npm run build

fe-install:
	@echo "📦 Installing frontend dependencies..."
	cd tealinux-fe && npm install

# Backend specific
be-dev:
	@echo "⚙️  Starting backend development..."
	cd tealinuxbe && go run cmd/main.go

be-build:
	@echo "🔨 Building backend..."
	cd tealinuxbe && go build -o tealinux-api cmd/main.go

be-test:
	@echo "🧪 Testing backend..."
	cd tealinuxbe && go test ./...

# Database
db-up:
	@echo "🗄️  Starting database..."
	docker compose up postgres -d

db-down:
	@echo "🛑 Stopping database..."
	docker compose stop postgres

db-logs:
	docker compose logs -f postgres

# Quick commands
start: dev-up
stop: dev-down
restart: dev-down dev-up
