# 🚀 CI/CD & Deployment Guide - TeaLinuxOS Web

Panduan lengkap untuk setup CI/CD dan deployment aplikasi TeaLinuxOS ke VPS menggunakan GitHub Actions, Docker, dan Nginx.

---

## 📋 Daftar Isi

1. [Arsitektur Deployment](#arsitektur-deployment)
2. [Prasyarat](#prasyarat)
3. [Setup VPS](#setup-vps)
4. [Konfigurasi GitHub Secrets](#konfigurasi-github-secrets)
5. [Setup SSL Certificate](#setup-ssl-certificate)
6. [Deployment Manual](#deployment-manual)
7. [CI/CD Pipeline](#cicd-pipeline)
8. [Monitoring & Maintenance](#monitoring--maintenance)
9. [Troubleshooting](#troubleshooting)

---

## 🏗️ Arsitektur Deployment

```
┌─────────────────────────────────────────────────────┐
│                    Internet                          │
└────────────────────┬────────────────────────────────┘
                     │
                     ▼
         ┌───────────────────────┐
         │   Nginx (Port 80/443) │
         │   - Reverse Proxy     │
         │   - SSL Termination   │
         │   - Rate Limiting     │
         └───────────┬───────────┘
                     │
         ┌───────────┴───────────┐
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌──────────────────┐
│  Frontend       │    │  Backend API     │
│  Astro SSR      │    │  Go Fiber        │
│  (Port 4321)    │    │  (Port 3000)     │
└─────────────────┘    └────────┬─────────┘
                                │
                                ▼
                    ┌───────────────────────┐
                    │   PostgreSQL 16       │
                    │   (Port 5432)         │
                    └───────────────────────┘
```

**Stack:**
- **Frontend**: Astro (SSR mode) + Tailwind CSS
- **Backend**: Go (Fiber framework) + GORM
- **Database**: PostgreSQL 16
- **Web Server**: Nginx (reverse proxy)
- **SSL**: Let's Encrypt (Certbot)
- **Container**: Docker + Docker Compose
- **CI/CD**: GitHub Actions

---

## ✅ Prasyarat

### 1. VPS Requirements
- **OS**: Ubuntu 22.04 LTS atau lebih baru
- **RAM**: Minimal 2GB (recommended 4GB)
- **Storage**: Minimal 20GB
- **CPU**: Minimal 2 cores
- **Network**: Public IP address

### 2. Domain
- Domain yang sudah di-pointing ke IP VPS
- Contoh: `tealinuxos.org` → `123.456.789.0`
- DNS A record untuk `@` dan `www`

### 3. Software di VPS
- Docker Engine (latest)
- Docker Compose (latest)
- Git
- SSH access

---

## 🖥️ Setup VPS

### Langkah 1: Koneksi ke VPS

```bash
ssh root@YOUR_VPS_IP
# atau
ssh your_username@YOUR_VPS_IP
```

### Langkah 2: Update System

```bash
sudo apt update && sudo apt upgrade -y
```

### Langkah 3: Install Docker

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group (optional, untuk non-root user)
sudo usermod -aG docker $USER

# Install Docker Compose
sudo apt install docker-compose-plugin -y

# Verify installation
docker --version
docker compose version
```

### Langkah 4: Install Git

```bash
sudo apt install git -y
git --version
```

### Langkah 5: Setup Project Directory

```bash
# Create project directory
sudo mkdir -p /opt/tealinux
sudo chown $USER:$USER /opt/tealinux
cd /opt/tealinux

# Clone repository
git clone https://github.com/YOUR_USERNAME/web-tealinux-astro.git .

# Atau jika sudah ada repository
git init
git remote add origin https://github.com/YOUR_USERNAME/web-tealinux-astro.git
git pull origin main
```

### Langkah 6: Setup Environment Variables

```bash
# Copy environment template
cp .env.production .env

# Edit environment variables
nano .env
```

**Isi file `.env`:**

```env
# Database
DB_NAME=tealinux
DB_USER=tealinux_user
DB_PASSWORD=YOUR_SECURE_PASSWORD_HERE

# Backend
DB_HOST=postgres
DB_PORT=5432
DB_SSLMODE=disable

# OAuth - Google
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=https://tealinuxos.org/api/auth/google/callback

# OAuth - GitHub
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=https://tealinuxos.org/api/auth/github/callback

# JWT Secret (generate dengan: openssl rand -base64 32)
JWT_SECRET=YOUR_RANDOM_JWT_SECRET_HERE

# Frontend
NODE_ENV=production
HOST=0.0.0.0
PORT=4321
API_BASE_URL=http://backend:3000

# Domain
DOMAIN=tealinuxos.org
EMAIL=admin@tealinuxos.org
```

**Generate JWT Secret:**

```bash
openssl rand -base64 32
```

### Langkah 7: Setup Backend Environment

```bash
# Copy backend environment
cp tealinuxbe/.env tealinuxbe/.env.backup
nano tealinuxbe/.env
```

Pastikan isi sesuai dengan `.env` utama.

---

## 🔐 Konfigurasi GitHub Secrets

Buka repository GitHub → **Settings** → **Secrets and variables** → **Actions** → **New repository secret**

Tambahkan secrets berikut:

| Secret Name | Description | Example |
|-------------|-------------|---------|
| `VPS_HOST` | IP address VPS | `123.456.789.0` |
| `VPS_USERNAME` | SSH username | `root` atau `ubuntu` |
| `VPS_SSH_KEY` | Private SSH key | Isi dari `~/.ssh/id_rsa` |
| `VPS_PORT` | SSH port (optional) | `22` (default) |
| `VPS_PROJECT_PATH` | Path project di VPS | `/opt/tealinux` |

### Generate SSH Key (jika belum ada)

**Di komputer lokal:**

```bash
# Generate SSH key pair
ssh-keygen -t rsa -b 4096 -C "github-actions@tealinuxos.org"

# Copy public key ke VPS
ssh-copy-id -i ~/.ssh/id_rsa.pub your_username@YOUR_VPS_IP

# Copy private key untuk GitHub Secret
cat ~/.ssh/id_rsa
# Copy seluruh output (termasuk BEGIN dan END)
```

**Paste private key ke GitHub Secret `VPS_SSH_KEY`**

---

## 🔒 Setup SSL Certificate

### Opsi 1: Menggunakan Script Otomatis (Recommended)

```bash
cd /opt/tealinux

# Jalankan script setup SSL
./setup-ssl.sh
```

### Opsi 2: Manual Setup

```bash
# Create directories
mkdir -p certbot/conf certbot/www

# Temporary nginx config (HTTP only)
# Edit nginx/nginx.conf, comment SSL section temporarily

# Start services
docker compose up -d

# Obtain certificate
docker compose run --rm certbot certonly \
  --webroot \
  --webroot-path=/var/www/certbot \
  --email admin@tealinuxos.org \
  --agree-tos \
  --no-eff-email \
  -d tealinuxos.org \
  -d www.tealinuxos.org

# Restart nginx with SSL
docker compose restart nginx
```

### Verifikasi SSL

```bash
# Check certificate
docker compose exec certbot certbot certificates

# Test auto-renewal
docker compose exec certbot certbot renew --dry-run
```

---

## 🚀 Deployment Manual

### First Time Deployment

```bash
cd /opt/tealinux

# Build and start all services
docker compose up -d --build

# Check logs
docker compose logs -f

# Check running containers
docker compose ps
```

### Update Deployment

```bash
cd /opt/tealinux

# Pull latest code
git pull origin main

# Rebuild and restart
docker compose down
docker compose up -d --build

# Clean old images
docker image prune -af
```

### Useful Commands

```bash
# View logs
docker compose logs -f [service_name]

# Restart specific service
docker compose restart [service_name]

# Stop all services
docker compose down

# Start all services
docker compose up -d

# Execute command in container
docker compose exec [service_name] [command]

# Database backup
docker compose exec postgres pg_dump -U tealinux_user tealinux > backup.sql

# Database restore
docker compose exec -T postgres psql -U tealinux_user tealinux < backup.sql
```

---

## 🔄 CI/CD Pipeline

### Workflow Overview

Pipeline berjalan otomatis saat:
- **Push** ke branch `main` atau `develop`
- **Pull Request** ke branch `main`
- **Manual trigger** via GitHub Actions UI

### Pipeline Stages

```
┌─────────────────────────────────────────────────┐
│  Stage 1: Test Frontend (Astro)                │
│  - Install dependencies                         │
│  - Build production bundle                      │
└─────────────────┬───────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────┐
│  Stage 2: Test Backend (Go)                     │
│  - Install dependencies                         │
│  - Run tests                                    │
│  - Build binary                                 │
└─────────────────┬───────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────┐
│  Stage 3: Build & Push Docker Images            │
│  - Build frontend image → GHCR                  │
│  - Build backend image → GHCR                   │
│  - Tag: latest, branch name, commit SHA        │
└─────────────────┬───────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────┐
│  Stage 4: Deploy to VPS                         │
│  - SSH to VPS                                   │
│  - Pull latest code                             │
│  - Pull Docker images                           │
│  - Restart services                             │
│  - Verify deployment                            │
└─────────────────────────────────────────────────┘
```

### Manual Trigger

1. Buka repository di GitHub
2. Klik tab **Actions**
3. Pilih workflow **CI/CD Pipeline - Deploy to VPS**
4. Klik **Run workflow**
5. Pilih branch dan klik **Run workflow**

### View Logs

1. Buka tab **Actions**
2. Klik workflow run yang sedang berjalan
3. Klik job untuk melihat detail logs

---

## 📊 Monitoring & Maintenance

### Health Checks

```bash
# Check all services status
docker compose ps

# Check specific service health
docker compose exec frontend wget -q -O- http://localhost:4321/health || echo "Frontend down"
docker compose exec backend wget -q -O- http://localhost:3000/health || echo "Backend down"

# Check database
docker compose exec postgres pg_isready -U tealinux_user
```

### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f frontend
docker compose logs -f backend
docker compose logs -f postgres
docker compose logs -f nginx

# Last 100 lines
docker compose logs --tail=100 frontend
```

### Resource Usage

```bash
# Container stats
docker stats

# Disk usage
docker system df

# Clean up
docker system prune -a --volumes
```

### Database Backup

**Automated Backup Script:**

```bash
#!/bin/bash
# backup-db.sh

BACKUP_DIR="/opt/tealinux/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/tealinux_$DATE.sql"

mkdir -p $BACKUP_DIR

docker compose exec -T postgres pg_dump -U tealinux_user tealinux > $BACKUP_FILE

# Compress
gzip $BACKUP_FILE

# Keep only last 7 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_FILE.gz"
```

**Setup Cron Job:**

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /opt/tealinux/backup-db.sh
```

### SSL Certificate Renewal

Certbot container otomatis memperbaharui certificate setiap 12 jam. Untuk manual renewal:

```bash
docker compose exec certbot certbot renew
docker compose restart nginx
```

---

## 🔧 Troubleshooting

### 1. Container Tidak Bisa Start

```bash
# Check logs
docker compose logs [service_name]

# Check port conflicts
sudo netstat -tulpn | grep LISTEN

# Restart Docker daemon
sudo systemctl restart docker
```

### 2. Database Connection Error

```bash
# Check PostgreSQL is running
docker compose ps postgres

# Check database logs
docker compose logs postgres

# Test connection
docker compose exec backend ping postgres

# Verify credentials in .env
```

### 3. SSL Certificate Issues

```bash
# Check certificate validity
docker compose exec certbot certbot certificates

# Renew manually
docker compose exec certbot certbot renew --force-renewal

# Check nginx config
docker compose exec nginx nginx -t
```

### 4. Frontend/Backend Not Accessible

```bash
# Check nginx logs
docker compose logs nginx

# Test internal connectivity
docker compose exec nginx curl http://frontend:4321
docker compose exec nginx curl http://backend:3000

# Check firewall
sudo ufw status
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
```

### 5. GitHub Actions Deployment Failed

**Check SSH Connection:**

```bash
# Test from local machine
ssh -i ~/.ssh/id_rsa your_username@YOUR_VPS_IP

# Verify SSH key in GitHub Secrets
```

**Check VPS Disk Space:**

```bash
df -h
docker system df
docker system prune -a
```

### 6. High Memory Usage

```bash
# Check container memory
docker stats

# Restart specific service
docker compose restart [service_name]

# Increase VPS RAM or optimize application
```

---

## 📝 Checklist Deployment

### Pre-Deployment

- [ ] VPS sudah setup (Docker, Git installed)
- [ ] Domain sudah pointing ke VPS IP
- [ ] Environment variables sudah dikonfigurasi
- [ ] GitHub Secrets sudah ditambahkan
- [ ] SSH key sudah di-setup

### Initial Deployment

- [ ] Clone repository ke VPS
- [ ] Setup `.env` file
- [ ] Build dan start containers
- [ ] Setup SSL certificate
- [ ] Verify semua services running
- [ ] Test website accessible

### Post-Deployment

- [ ] Setup database backup cron job
- [ ] Configure monitoring (optional)
- [ ] Test CI/CD pipeline
- [ ] Document custom configurations

---

## 🆘 Support

Jika mengalami masalah:

1. Check logs: `docker compose logs -f`
2. Verify environment variables
3. Check GitHub Actions logs
4. Review this documentation
5. Contact team: admin@tealinuxos.org

---

## 📚 Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [GitHub Actions](https://docs.github.com/en/actions)
- [Let's Encrypt](https://letsencrypt.org/)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Astro SSR](https://docs.astro.build/en/guides/server-side-rendering/)

---

**Last Updated**: 2026-02-11  
**Version**: 1.0.0
