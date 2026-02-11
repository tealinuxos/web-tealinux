# 📋 Ringkasan Setup CI/CD & Deployment TeaLinuxOS

## ✅ Yang Telah Dibuat

### 1. **Docker Configuration**

#### Production (`docker-compose.yml`)
- ✅ PostgreSQL 16 dengan health checks
- ✅ Go Backend API dengan dependency management
- ✅ Astro Frontend (SSR mode)
- ✅ Nginx reverse proxy dengan SSL support
- ✅ Certbot untuk auto-renewal SSL certificates
- ✅ Internal Docker network untuk komunikasi antar services

#### Development (`docker-compose.dev.yml`)
- ✅ Hot-reload untuk frontend (Astro)
- ✅ Hot-reload untuk backend (Go dengan Air)
- ✅ Volume mounts untuk live code changes
- ✅ Isolated development environment

### 2. **Dockerfiles**

#### Frontend
- ✅ `tealinux-fe/Dockerfile` - Multi-stage production build
- ✅ `tealinux-fe/Dockerfile.dev` - Development dengan hot-reload

#### Backend
- ✅ `tealinuxbe/Dockerfile` - Multi-stage production build (Go binary)
- ✅ `tealinuxbe/Dockerfile.dev` - Development dengan Air hot-reload
- ✅ `tealinuxbe/.air.toml` - Air configuration

### 3. **Nginx Configuration**

✅ `nginx/nginx.conf` dengan fitur:
- HTTP to HTTPS redirect
- SSL/TLS configuration (TLS 1.2 & 1.3)
- Security headers (HSTS, X-Frame-Options, etc.)
- Rate limiting untuk API dan general requests
- Gzip compression
- Reverse proxy ke frontend dan backend
- Certbot webroot untuk SSL challenges

### 4. **GitHub Actions CI/CD**

✅ `.github/workflows/deploy.yml` dengan stages:

**Stage 1: Test Frontend**
- Checkout code
- Setup Node.js 20
- Install dependencies
- Build production bundle

**Stage 2: Test Backend**
- Checkout code
- Setup Go 1.25
- Run tests
- Build binary

**Stage 3: Build & Push Docker Images**
- Build frontend image → GitHub Container Registry
- Build backend image → GitHub Container Registry
- Tag dengan: latest, branch name, commit SHA
- Layer caching untuk build cepat

**Stage 4: Deploy to VPS**
- SSH ke VPS
- Pull latest code dari Git
- Pull Docker images dari GHCR
- Restart services dengan docker compose
- Cleanup old images
- Verify deployment

### 5. **Helper Scripts**

✅ **`setup-ssl.sh`** - Setup SSL certificates dengan Certbot
- Automated certificate obtaining
- Domain validation
- Auto-renewal setup

✅ **`deploy.sh`** - Quick deployment script
- Pull latest code
- Pull Docker images
- Restart services
- Run health check

✅ **`health-check.sh`** - Comprehensive health monitoring
- Docker status
- Container health
- Frontend/Backend/Database connectivity
- SSL certificate validity
- Disk usage
- Memory usage
- Recent error logs

✅ **`backup-db.sh`** - Automated database backup
- PostgreSQL dump
- Gzip compression
- 7-day retention policy
- Backup size reporting

✅ **`restore-db.sh`** - Database restore
- Safety confirmations
- Connection handling
- Database recreation
- Service restart

### 6. **Environment Configuration**

✅ **`.env.production`** - Production environment template
- Database credentials
- OAuth configurations (Google, GitHub)
- JWT secrets
- Domain settings
- API URLs

✅ **`.gitignore`** - Updated untuk exclude:
- Environment files
- Build outputs
- Docker volumes
- SSL certificates
- Logs

### 7. **Documentation**

✅ **`DEPLOYMENT.md`** (Comprehensive, 500+ lines)
- Arsitektur deployment
- VPS setup step-by-step
- GitHub Secrets configuration
- SSL certificate setup
- Manual deployment guide
- CI/CD pipeline explanation
- Monitoring & maintenance
- Troubleshooting guide
- Checklist deployment

✅ **`QUICK_REFERENCE.md`**
- Quick commands
- Project structure
- Common operations
- URLs reference

✅ **`README.md`** - Updated dengan:
- DevOps & Infrastructure section
- Deployment & CI/CD section
- Links ke detailed guides
- GitHub Secrets requirements
- CI/CD pipeline overview

---

## 🚀 Cara Menggunakan

### A. Development (Lokal)

#### Opsi 1: Tanpa Docker
```bash
# Terminal 1 - Database
docker compose up postgres -d

# Terminal 2 - Backend
cd tealinuxbe
go run cmd/main.go

# Terminal 3 - Frontend
cd tealinux-fe
npm run dev
```

#### Opsi 2: Dengan Docker (Hot-reload)
```bash
docker compose -f docker-compose.dev.yml up -d
docker compose -f docker-compose.dev.yml logs -f
```

### B. Production Deployment

#### Setup Awal di VPS

**1. Persiapan VPS**
```bash
# SSH ke VPS
ssh root@YOUR_VPS_IP

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt install docker-compose-plugin -y

# Clone repository
sudo mkdir -p /opt/tealinux
sudo chown $USER:$USER /opt/tealinux
cd /opt/tealinux
git clone https://github.com/YOUR_USERNAME/web-tealinux-astro.git .
```

**2. Konfigurasi Environment**
```bash
# Copy dan edit .env
cp .env.production .env
nano .env

# Generate JWT secret
openssl rand -base64 32

# Update .env dengan:
# - Database password (secure)
# - OAuth credentials
# - JWT secret
# - Domain name
```

**3. Setup SSL Certificate**
```bash
# Edit domain di .env
export DOMAIN=tealinuxos.org
export EMAIL=admin@tealinuxos.org

# Run SSL setup
./setup-ssl.sh
```

**4. Deploy**
```bash
# Build dan start semua services
docker compose up -d --build

# Check logs
docker compose logs -f

# Health check
./health-check.sh
```

#### Setup GitHub Actions

**1. Generate SSH Key**
```bash
# Di komputer lokal
ssh-keygen -t rsa -b 4096 -C "github-actions@tealinuxos.org"

# Copy public key ke VPS
ssh-copy-id -i ~/.ssh/id_rsa.pub your_username@YOUR_VPS_IP

# Copy private key untuk GitHub Secret
cat ~/.ssh/id_rsa
```

**2. Tambahkan GitHub Secrets**

Buka: Repository → Settings → Secrets and variables → Actions

Tambahkan:
- `VPS_HOST`: IP VPS (contoh: `123.456.789.0`)
- `VPS_USERNAME`: SSH username (contoh: `root`)
- `VPS_SSH_KEY`: Private SSH key (paste seluruh isi `~/.ssh/id_rsa`)
- `VPS_PORT`: `22` (atau custom port)
- `VPS_PROJECT_PATH`: `/opt/tealinux`

**3. Test CI/CD**
```bash
# Push ke main branch
git add .
git commit -m "feat: setup CI/CD"
git push origin main

# Atau manual trigger di GitHub Actions tab
```

### C. Maintenance

#### Database Backup (Automated)
```bash
# Setup cron job untuk backup harian
crontab -e

# Tambahkan (backup setiap hari jam 2 pagi):
0 2 * * * /opt/tealinux/backup-db.sh
```

#### Manual Operations
```bash
# Backup database
./backup-db.sh

# Restore database
./restore-db.sh backups/tealinux_20260211_120000.sql.gz

# Health check
./health-check.sh

# Quick deploy
./deploy.sh

# View logs
docker compose logs -f [service_name]

# Restart service
docker compose restart [service_name]
```

---

## 📊 Monitoring

### Health Check
```bash
./health-check.sh
```

Output akan menampilkan:
- ✅ Docker status
- ✅ Container status
- ✅ Frontend health
- ✅ Backend health
- ✅ Database health
- ✅ Nginx configuration
- ✅ SSL certificate validity
- ✅ Disk usage
- ✅ Memory usage
- ✅ Recent errors

### Logs
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

# Follow with grep
docker compose logs -f | grep -i error
```

---

## 🔧 Troubleshooting

### Container tidak start
```bash
docker compose logs [service_name]
docker compose restart [service_name]
```

### SSL certificate issues
```bash
docker compose exec certbot certbot certificates
docker compose exec certbot certbot renew --force-renewal
```

### Database connection error
```bash
docker compose exec postgres pg_isready -U tealinux_user
docker compose logs postgres
```

### Disk space penuh
```bash
df -h
docker system df
docker system prune -a
```

---

## 📚 File Structure

```
web-tealinux-astro/
├── .github/
│   └── workflows/
│       └── deploy.yml              # CI/CD pipeline
├── nginx/
│   └── nginx.conf                  # Nginx configuration
├── tealinux-fe/
│   ├── Dockerfile                  # Production build
│   ├── Dockerfile.dev              # Development build
│   └── ...
├── tealinuxbe/
│   ├── Dockerfile                  # Production build
│   ├── Dockerfile.dev              # Development build
│   ├── .air.toml                   # Hot-reload config
│   └── ...
├── docker-compose.yml              # Production compose
├── docker-compose.dev.yml          # Development compose
├── .env.production                 # Environment template
├── .gitignore                      # Updated gitignore
├── deploy.sh                       # Quick deploy script
├── health-check.sh                 # Health monitoring
├── backup-db.sh                    # Database backup
├── restore-db.sh                   # Database restore
├── setup-ssl.sh                    # SSL setup
├── DEPLOYMENT.md                   # Full deployment guide
├── QUICK_REFERENCE.md              # Quick commands
└── README.md                       # Updated README
```

---

## ✅ Checklist Deployment

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
- [ ] Test CI/CD pipeline
- [ ] Monitor logs untuk errors
- [ ] Document custom configurations

---

## 🎯 Next Steps

1. **Setup VPS** - Ikuti panduan di `DEPLOYMENT.md`
2. **Configure GitHub Secrets** - Tambahkan VPS credentials
3. **Test CI/CD** - Push ke main branch atau manual trigger
4. **Setup Monitoring** - Configure health checks dan alerts
5. **Database Backup** - Setup automated backups

---

## 📞 Support

Jika ada pertanyaan atau masalah:
1. Cek `DEPLOYMENT.md` untuk panduan lengkap
2. Cek `QUICK_REFERENCE.md` untuk command reference
3. Run `./health-check.sh` untuk diagnostics
4. Check logs: `docker compose logs -f`

---

**Created**: 2026-02-11  
**Version**: 1.0.0  
**Status**: ✅ Ready for deployment
