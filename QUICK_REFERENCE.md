# TeaLinuxOS Web - Quick Reference

## 🚀 Quick Commands

### Development
```bash
# Frontend (Astro)
cd tealinux-fe
npm install
npm run dev              # http://localhost:4321

# Backend (Go)
cd tealinuxbe
go mod download
go run cmd/main.go       # http://localhost:3000
```

### Docker (Local)
```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Stop all services
docker compose down
```

### Production Deployment
```bash
# Quick deploy
./deploy.sh

# Health check
./health-check.sh

# Database backup
./backup-db.sh

# Database restore
./restore-db.sh backups/tealinux_20260211_120000.sql.gz

# SSL setup (first time only)
./setup-ssl.sh
```

## 📁 Project Structure

```
web-tealinux-astro/
├── tealinux-fe/              # Astro frontend (SSR)
│   ├── src/
│   ├── public/
│   ├── Dockerfile
│   └── package.json
├── tealinuxbe/               # Go backend API
│   ├── cmd/
│   ├── internal/
│   ├── Dockerfile
│   └── go.mod
├── nginx/
│   └── nginx.conf            # Nginx configuration
├── .github/
│   └── workflows/
│       └── deploy.yml        # CI/CD pipeline
├── docker-compose.yml        # Docker orchestration
├── .env.production           # Environment template
├── deploy.sh                 # Quick deploy script
├── health-check.sh           # Health check script
├── backup-db.sh              # Database backup
├── restore-db.sh             # Database restore
├── setup-ssl.sh              # SSL setup
└── DEPLOYMENT.md             # Full documentation
```

## 🔑 GitHub Secrets Required

| Secret | Description |
|--------|-------------|
| `VPS_HOST` | VPS IP address |
| `VPS_USERNAME` | SSH username |
| `VPS_SSH_KEY` | Private SSH key |
| `VPS_PORT` | SSH port (default: 22) |
| `VPS_PROJECT_PATH` | Project path (default: /opt/tealinux) |

## 🌐 URLs

- **Production**: https://tealinuxos.org
- **API**: https://tealinuxos.org/api
- **Frontend (local)**: http://localhost:4321
- **Backend (local)**: http://localhost:3000

## 📊 Monitoring

```bash
# Container stats
docker stats

# View logs
docker compose logs -f [service]

# Check services
docker compose ps

# Health check
./health-check.sh
```

## 🆘 Troubleshooting

```bash
# Restart all services
docker compose restart

# Rebuild specific service
docker compose up -d --build [service]

# Clean Docker cache
docker system prune -a

# View error logs
docker compose logs --tail=100 [service] | grep -i error
```

## 📚 Documentation

- **Full Deployment Guide**: [DEPLOYMENT.md](./DEPLOYMENT.md)
- **Astro Docs**: https://docs.astro.build
- **Docker Compose**: https://docs.docker.com/compose/
- **GitHub Actions**: https://docs.github.com/en/actions
