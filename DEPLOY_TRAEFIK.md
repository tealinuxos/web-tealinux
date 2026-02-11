# 🚀 Panduan Deploy TeaLinux Web (Traefik Mode)

> Panduan ini menjelaskan cara deploy web TeaLinux menggunakan **Docker Compose + Traefik**,
> dengan konsep yang sama seperti deployment ecoBite (savorbite.doscom.org).

---

## 📐 Arsitektur

```
                    ┌──────────────────────────────┐
                    │         INTERNET              │
                    └──────────┬───────────────────┘
                               │
                    ┌──────────▼───────────────────┐
                    │      TRAEFIK (Shared)         │
                    │   Port 80 (HTTP)              │
                    │   Port 443 (HTTPS + SSL Auto) │
                    │   Network: proxy              │
                    └──────────┬───────────────────┘
                               │
          ┌────────────────────┼────────────────────┐
          │                    │                     │
    ┌─────▼──────┐     ┌──────▼──────┐      ┌──────▼──────┐
    │  TeaLinux   │     │  ecoBite    │      │  Lainnya    │
    │  Frontend   │     │  Frontend   │      │  ...        │
    │  :4321      │     │  :3000      │      │             │
    │  Astro SSR  │     │  Next.js    │      │             │
    └─────┬───────┘     └─────────────┘      └─────────────┘
          │
    ┌─────▼──────┐     ┌──────────────┐
    │  TeaLinux   │─────│  PostgreSQL  │
    │  Backend    │     │  :5432       │
    │  :3000 (Go) │     │             │
    └─────────────┘     └─────────────┘
```

**Routing Rules:**
- `tealinuxos.org/*` → Frontend (Astro SSR, port 4321)
- `tealinuxos.org/api/*` → Backend (Go, port 3000) — dengan priority lebih tinggi
- `savorbite.doscom.org/*` → ecoBite Frontend
- `savorbite.doscom.org/api/*` → ecoBite Backend

---

## 📋 Prerequisites

1. **VPS** dengan Docker & Docker Compose terinstall
2. **Traefik** sudah running sebagai shared reverse proxy
3. **Domain** `tealinuxos.org` sudah pointing ke IP VPS
4. **Git** repository sudah di-clone di VPS

---

## 🔧 Langkah-Langkah Detail

### Langkah 1: SSH ke VPS

```bash
ssh root@<IP_VPS>
# atau
ssh root@doscom-server
```

### Langkah 2: Verifikasi Traefik Sudah Running

```bash
# Cek Traefik container
docker ps | grep traefik

# Cek proxy network
docker network ls | grep proxy
```

**Jika Traefik belum ada**, setup dulu:

```bash
# Buat network external
docker network create proxy

# Setup Traefik
mkdir -p /opt/traefik && cd /opt/traefik

# Buat docker-compose.yml untuk Traefik (lihat file di bawah)
# Lalu:
docker compose up -d
```

### Langkah 3: Clone/Update Repository

```bash
# Jika belum pernah clone:
cd /opt
git clone https://github.com/tealinuxos/web-tealinux.git tealinux
cd tealinux

# Jika sudah pernah clone:
cd /opt/tealinux
git pull origin main
```

### Langkah 4: Setup Environment Variables

```bash
# Copy template ke .env
cp .env.traefik .env

# Edit dan isi nilai yang benar
nano .env
```

**Yang HARUS diubah:**

| Variable | Keterangan |
|----------|-----------|
| `DB_PASSWORD` | Ganti dengan password database yang kuat |
| `JWT_SECRET` | Generate: `openssl rand -base64 64` |
| `GOOGLE_CLIENT_ID` | OAuth Google Client ID |
| `GOOGLE_CLIENT_SECRET` | OAuth Google Client Secret |
| `GITHUB_CLIENT_ID` | OAuth GitHub Client ID |
| `GITHUB_CLIENT_SECRET` | OAuth GitHub Client Secret |

### Langkah 5: Deploy!

#### Opsi A: Menggunakan Script (Recommended)

```bash
chmod +x deploy-traefik.sh
./deploy-traefik.sh
```

#### Opsi B: Manual

```bash
# 1. Pastikan proxy network ada
docker network create proxy 2>/dev/null || true

# 2. Build images
docker compose -f docker-compose.traefik.yml build

# 3. Start services
docker compose -f docker-compose.traefik.yml up -d

# 4. Cek status
docker compose -f docker-compose.traefik.yml ps
```

### Langkah 6: Verifikasi Deployment

```bash
# Cek semua container running
docker compose -f docker-compose.traefik.yml ps

# Cek logs
docker compose -f docker-compose.traefik.yml logs -f

# Cek backend health
docker exec tealinux-postgres pg_isready -U tealinux_user -d tealinux

# Cek logs per service
docker logs tealinux-backend --tail 50
docker logs tealinux-frontend --tail 50

# Test endpoint (dari server)
curl -I http://localhost:4321    # Frontend (internal)
curl -I http://localhost:3000    # Backend (internal)
```

### Langkah 7: Verifikasi dari Browser

1. Buka `https://tealinuxos.org` — harus menampilkan frontend
2. Buka `https://tealinuxos.org/api/...` — harus menampilkan response API
3. Cek SSL certificate valid (gembok hijau di browser)

---

## 🔄 CI/CD (Auto Deploy via GitHub Actions)

### Setup GitHub Secrets

Buka **GitHub Repository → Settings → Secrets and variables → Actions**, tambahkan:

| Secret Name | Nilai | Keterangan |
|---|---|---|
| `VPS_HOST` | IP atau hostname VPS | Contoh: `103.xxx.xxx.xxx` |
| `VPS_USERNAME` | Username SSH | Biasanya `root` |
| `VPS_SSH_KEY` | Private key SSH | Generate: `ssh-keygen -t ed25519` |
| `VPS_PORT` | Port SSH | Default: `22` |
| `VPS_PROJECT_PATH` | Path project di VPS | Contoh: `/opt/tealinux` |
| `GHCR_TOKEN` | GitHub PAT (packages) | Token untuk pull image |

### Setup SSH Key

```bash
# Di local machine, generate SSH key
ssh-keygen -t ed25519 -C "github-actions-tealinux" -f ~/.ssh/github_actions_tealinux

# Copy public key ke VPS
ssh-copy-id -i ~/.ssh/github_actions_tealinux.pub root@<VPS_IP>

# Copy private key content — paste ke GitHub Secret VPS_SSH_KEY
cat ~/.ssh/github_actions_tealinux
```

### Workflow

Setiap push ke branch `main`:
1. ✅ **Test Frontend** — Build Astro project
2. ✅ **Test Backend** — Run Go tests  
3. 📦 **Build & Push** — Build Docker images → Push ke GHCR
4. 🚀 **Deploy** — SSH ke VPS → Pull latest → Restart containers
5. ✅ **Verify** — Check all containers healthy

---

## 🛠️ Perintah Berguna

```bash
# Lihat semua container
docker compose -f docker-compose.traefik.yml ps

# Lihat logs realtime
docker compose -f docker-compose.traefik.yml logs -f

# Lihat logs service tertentu
docker compose -f docker-compose.traefik.yml logs -f backend
docker compose -f docker-compose.traefik.yml logs -f frontend

# Restart satu service
docker compose -f docker-compose.traefik.yml restart backend
docker compose -f docker-compose.traefik.yml restart frontend

# Rebuild dan restart
docker compose -f docker-compose.traefik.yml up -d --build backend

# Stop semua
docker compose -f docker-compose.traefik.yml down

# Stop dan hapus volume (⚠️ DATA HILANG)
docker compose -f docker-compose.traefik.yml down -v

# Masuk ke container
docker exec -it tealinux-backend sh
docker exec -it tealinux-frontend sh
docker exec -it tealinux-postgres psql -U tealinux_user -d tealinux
```

---

## 🔍 Troubleshooting

### Container tidak start

```bash
# Cek error logs
docker compose -f docker-compose.traefik.yml logs backend
docker compose -f docker-compose.traefik.yml logs frontend

# Cek apakah port bentrok
docker ps --format "table {{.Names}}\t{{.Ports}}"
```

### Traefik tidak routing

```bash
# Cek Traefik logs
docker logs traefik --tail 50

# Cek apakah container terhubung ke proxy network
docker network inspect proxy

# Cek apakah labels benar
docker inspect tealinux-frontend | grep -A 50 Labels
docker inspect tealinux-backend | grep -A 50 Labels
```

### Database connection error

```bash
# Cek postgres health
docker exec tealinux-postgres pg_isready -U tealinux_user -d tealinux

# Cek environment variables
docker exec tealinux-backend env | grep DB

# Cek network connectivity
docker exec tealinux-backend ping -c 3 postgres
```

### SSL Certificate tidak muncul

```bash
# Cek Traefik resolver logs
docker logs traefik 2>&1 | grep -i cert

# Pastikan domain sudah pointing ke VPS
dig tealinuxos.org
nslookup tealinuxos.org

# Pastikan port 80 dan 443 terbuka
sudo ufw status
```

---

## 📊 Perbandingan dengan ecoBite

| Aspek | ecoBite | TeaLinux |
|-------|---------|----------|
| Domain | `savorbite.doscom.org` | `tealinuxos.org` |
| Frontend | Next.js (:3000) | Astro SSR (:4321) |
| Backend | Go (:8080) | Go (:3000) |
| Database | PostgreSQL 15 | PostgreSQL 16 |
| API Route | `/api` → :8080 | `/api` → :3000 |
| Web Route | `/` → :3000 | `/` → :4321 |
| SSL | Traefik Let's Encrypt | Traefik Let's Encrypt |
| Network | proxy + ecobite-network | proxy + tealinux-network |
