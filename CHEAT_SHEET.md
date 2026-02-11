# 🚀 TeaLinuxOS CI/CD - Cheat Sheet

## 📦 Quick Start

### Development (Lokal)
```bash
make dev              # Start dev environment dengan hot-reload
make dev-up           # Start di background
make dev-down         # Stop
make logs             # View logs
```

### Production (VPS)
```bash
./deploy.sh           # Deploy
./health-check.sh     # Check health
./backup-db.sh        # Backup DB
make prod-up          # Start production
```

## 🔑 GitHub Secrets (Required)

```
VPS_HOST          = 123.456.789.0
VPS_USERNAME      = root
VPS_SSH_KEY       = <private SSH key>
VPS_PROJECT_PATH  = /opt/tealinux
```

## 📁 Files Created

```
✅ docker-compose.yml           Production orchestration
✅ docker-compose.dev.yml       Development with hot-reload
✅ .github/workflows/deploy.yml CI/CD pipeline
✅ nginx/nginx.conf             Reverse proxy + SSL
✅ tealinux-fe/Dockerfile       Frontend production
✅ tealinux-fe/Dockerfile.dev   Frontend development
✅ tealinuxbe/Dockerfile        Backend production
✅ tealinuxbe/Dockerfile.dev    Backend development
✅ tealinuxbe/.air.toml         Go hot-reload config
✅ deploy.sh                    Quick deploy
✅ health-check.sh              Health monitoring
✅ backup-db.sh                 Database backup
✅ restore-db.sh                Database restore
✅ setup-ssl.sh                 SSL setup
✅ Makefile                     Command shortcuts
✅ DEPLOYMENT.md                Full guide (500+ lines)
✅ QUICK_REFERENCE.md           Quick commands
✅ CICD_SETUP_SUMMARY.md        Complete summary
```

## 🎯 Deployment Steps

### 1️⃣ Setup VPS
```bash
# Install Docker
curl -fsSL https://get.docker.com | sh

# Clone repo
cd /opt/tealinux
git clone <repo-url> .
```

### 2️⃣ Configure
```bash
# Setup environment
cp .env.production .env
nano .env  # Edit credentials

# Setup SSL
./setup-ssl.sh
```

### 3️⃣ Deploy
```bash
# First deploy
docker compose up -d --build

# Or use script
./deploy.sh
```

### 4️⃣ Setup GitHub Actions
```bash
# Generate SSH key
ssh-keygen -t rsa -b 4096

# Add to VPS
ssh-copy-id user@vps

# Add secrets to GitHub
# Settings → Secrets → Actions
```

## 🔄 CI/CD Pipeline

```
Push to main
    ↓
Test Frontend (Astro build)
    ↓
Test Backend (Go tests)
    ↓
Build Docker Images → GHCR
    ↓
Deploy to VPS via SSH
    ↓
Health Check ✅
```

## 🛠️ Common Commands

```bash
# Makefile shortcuts
make dev              # Dev environment
make deploy           # Deploy
make health           # Health check
make backup           # Backup DB
make logs             # View logs
make clean            # Clean Docker

# Docker commands
docker compose ps                    # Status
docker compose logs -f [service]     # Logs
docker compose restart [service]     # Restart
docker compose down                  # Stop all
docker compose up -d                 # Start all

# Database
./backup-db.sh                       # Backup
./restore-db.sh <file>               # Restore

# SSL
./setup-ssl.sh                       # Setup
docker compose exec certbot certbot renew  # Renew
```

## 📊 Monitoring

```bash
./health-check.sh     # Full health check
docker stats          # Resource usage
docker compose ps     # Container status
```

## 🆘 Troubleshooting

```bash
# View errors
docker compose logs --tail=100 | grep -i error

# Restart service
docker compose restart [service]

# Clean and rebuild
docker compose down
docker system prune -af
docker compose up -d --build

# Check SSL
docker compose exec certbot certbot certificates
```

## 📚 Documentation

- **DEPLOYMENT.md** - Complete deployment guide
- **QUICK_REFERENCE.md** - Quick commands
- **CICD_SETUP_SUMMARY.md** - Setup summary
- **README.md** - Project overview

## 🌐 URLs

- Production: `https://tealinuxos.org`
- Frontend (dev): `http://localhost:4321`
- Backend (dev): `http://localhost:3000`
- API: `https://tealinuxos.org/api`

---

**Need help?** Check `DEPLOYMENT.md` or run `make help`
