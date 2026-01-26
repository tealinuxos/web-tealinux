# 🚀 Deployment Guide

Panduan lengkap untuk deploy TeaLinuxOS Web Platform ke production.

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Frontend Deployment](#frontend-deployment)
- [Backend Deployment](#backend-deployment)
- [Database Setup](#database-setup)
- [Environment Configuration](#environment-configuration)
- [CI/CD Pipeline](#cicd-pipeline)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

---

## ✅ Prerequisites

### Required Services

- [ ] Domain name (e.g., tealinuxos.org)
- [ ] SSL Certificate (Let's Encrypt recommended)
- [ ] Cloud hosting account (VPS/Cloud provider)
- [ ] PostgreSQL database (managed or self-hosted)
- [ ] Git repository access

### Recommended Providers

**Frontend Hosting:**
- ✅ Vercel (Recommended)
- ✅ Netlify
- ✅ Cloudflare Pages
- ⚠️ GitHub Pages (limited features)

**Backend Hosting:**
- ✅ DigitalOcean (VPS)
- ✅ AWS EC2
- ✅ Google Cloud Run
- ✅ Heroku
- ✅ Railway

**Database:**
- ✅ DigitalOcean Managed PostgreSQL
- ✅ AWS RDS
- ✅ Supabase
- ✅ ElephantSQL
- ⚠️ Self-hosted PostgreSQL

---

## 🎨 Frontend Deployment

### Option 1: Vercel (Recommended)

#### Step 1: Prepare Project

```bash
cd tealinux-fe

# Test build locally
npm run build

# Preview build
npm run preview
```

#### Step 2: Deploy to Vercel

**Via Vercel CLI:**
```bash
# Install Vercel CLI
npm i -g vercel

# Login
vercel login

# Deploy
vercel

# Deploy to production
vercel --prod
```

**Via Vercel Dashboard:**

1. Go to [vercel.com](https://vercel.com)
2. Click "New Project"
3. Import Git repository
4. Configure:
   - Framework Preset: **Astro**
   - Root Directory: **tealinux-fe**
   - Build Command: `npm run build`
   - Output Directory: `dist`
5. Add Environment Variables (if needed)
6. Click "Deploy"

#### Step 3: Configure Domain

1. Go to Project Settings → Domains
2. Add your custom domain
3. Update DNS records:
   ```
   Type: CNAME
   Name: www
   Value: cname.vercel-dns.com
   ```

---

### Option 2: Netlify

#### Deploy via Netlify CLI

```bash
# Install Netlify CLI
npm i -g netlify-cli

# Login
netlify login

# Build
cd tealinux-fe
npm run build

# Deploy
netlify deploy --prod --dir=dist
```

#### Deploy via Git

1. Push to GitHub/GitLab
2. Connect repository to Netlify
3. Configure build settings:
   ```
   Base directory: tealinux-fe
   Build command: npm run build
   Publish directory: tealinux-fe/dist
   ```

---

### Option 3: Cloudflare Pages

```bash
# Install Wrangler
npm i -g wrangler

# Login
wrangler login

# Deploy
cd tealinux-fe
npm run build
wrangler pages deploy dist
```

---

## 🔧 Backend Deployment

### Option 1: DigitalOcean VPS (Recommended)

#### Step 1: Create Droplet

1. Create Ubuntu 22.04 droplet
2. Choose plan (minimum: 1GB RAM)
3. Add SSH key
4. Create droplet

#### Step 2: Initial Server Setup

```bash
# SSH into server
ssh root@your_server_ip

# Update system
apt update && apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version

# Install PostgreSQL
apt install postgresql postgresql-contrib -y

# Install Nginx
apt install nginx -y

# Install Certbot (for SSL)
apt install certbot python3-certbot-nginx -y
```

#### Step 3: Setup Application

```bash
# Create app user
adduser tealinux
usermod -aG sudo tealinux

# Switch to app user
su - tealinux

# Clone repository
git clone https://github.com/doscom/web-tealinux-astro.git
cd web-tealinux-astro/tealinuxbe

# Install dependencies
go mod download

# Create .env file
nano .env
# Add production environment variables

# Build application
go build -o tealinux-api cmd/main.go
```

#### Step 4: Create Systemd Service

```bash
# Create service file
sudo nano /etc/systemd/system/tealinux-api.service
```

Add:
```ini
[Unit]
Description=TeaLinux API Service
After=network.target postgresql.service

[Service]
Type=simple
User=tealinux
WorkingDirectory=/home/tealinux/web-tealinux-astro/tealinuxbe
ExecStart=/home/tealinux/web-tealinux-astro/tealinuxbe/tealinux-api
Restart=on-failure
RestartSec=5s

Environment="PORT=3000"
EnvironmentFile=/home/tealinux/web-tealinux-astro/tealinuxbe/.env

[Install]
WantedBy=multi-user.target
```

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service
sudo systemctl enable tealinux-api

# Start service
sudo systemctl start tealinux-api

# Check status
sudo systemctl status tealinux-api
```

#### Step 5: Configure Nginx

```bash
# Create Nginx config
sudo nano /etc/nginx/sites-available/tealinux-api
```

Add:
```nginx
server {
    listen 80;
    server_name api.tealinuxos.org;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/tealinux-api /etc/nginx/sites-enabled/

# Test Nginx config
sudo nginx -t

# Restart Nginx
sudo systemctl restart nginx
```

#### Step 6: Setup SSL

```bash
# Get SSL certificate
sudo certbot --nginx -d api.tealinuxos.org

# Auto-renewal (already setup by certbot)
sudo certbot renew --dry-run
```

---

### Option 2: Docker Deployment

#### Create Dockerfile

```dockerfile
# tealinuxbe/Dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o tealinux-api cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/tealinux-api .

# Expose port
EXPOSE 3000

# Run
CMD ["./tealinux-api"]
```

#### Create docker-compose.yml

```yaml
version: '3.9'

services:
  backend:
    build: ./tealinuxbe
    ports:
      - "3000:3000"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=tealinux_user
      - DB_PASSWORD=tealinux123
      - DB_NAME=tealinux
    depends_on:
      - postgres
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: tealinux
      POSTGRES_USER: tealinux_user
      POSTGRES_PASSWORD: tealinux123
    volumes:
      - pgdata:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  pgdata:
```

#### Deploy with Docker

```bash
# Build and run
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop
docker-compose down
```

---

## 🗄️ Database Setup

### Option 1: Managed PostgreSQL (Recommended)

#### DigitalOcean Managed Database

1. Create managed PostgreSQL database
2. Choose region and plan
3. Configure:
   - Version: PostgreSQL 16
   - Plan: Basic (1GB RAM minimum)
4. Add trusted sources (your backend server IP)
5. Get connection details

#### Update Backend .env

```env
DB_HOST=your-db-host.db.ondigitalocean.com
DB_PORT=25060
DB_USER=doadmin
DB_PASSWORD=your-secure-password
DB_NAME=tealinux
```

---

### Option 2: Self-Hosted PostgreSQL

```bash
# Install PostgreSQL
sudo apt install postgresql postgresql-contrib

# Create database and user
sudo -u postgres psql

CREATE DATABASE tealinux;
CREATE USER tealinux_user WITH PASSWORD 'your-secure-password';
GRANT ALL PRIVILEGES ON DATABASE tealinux TO tealinux_user;
\q

# Configure PostgreSQL for remote access (if needed)
sudo nano /etc/postgresql/14/main/postgresql.conf
# Change: listen_addresses = '*'

sudo nano /etc/postgresql/14/main/pg_hba.conf
# Add: host all all 0.0.0.0/0 md5

# Restart PostgreSQL
sudo systemctl restart postgresql
```

---

## 🔐 Environment Configuration

### Frontend Environment Variables

Create `.env.production` in `tealinux-fe/`:

```env
PUBLIC_API_URL=https://api.tealinuxos.org
PUBLIC_SITE_URL=https://tealinuxos.org
```

### Backend Environment Variables

Production `.env` in `tealinuxbe/`:

```env
# Server
PORT=3000
ENVIRONMENT=production

# Database
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=tealinux_user
DB_PASSWORD=your-secure-password
DB_NAME=tealinux

# JWT (CHANGE THESE!)
JWT_SECRET=your-super-secret-jwt-key-min-32-chars
JWT_EXPIRE=24h

# OAuth Google
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=https://api.tealinuxos.org/api/auth/google/callback

# OAuth GitHub
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=https://api.tealinuxos.org/api/auth/github/callback

# Frontend
FRONTEND_URL=https://tealinuxos.org

# CORS
ALLOWED_ORIGINS=https://tealinuxos.org,https://www.tealinuxos.org
```

---

## 🔄 CI/CD Pipeline

### GitHub Actions Example

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install dependencies
        working-directory: ./tealinux-fe
        run: npm ci
      
      - name: Build
        working-directory: ./tealinux-fe
        run: npm run build
      
      - name: Deploy to Vercel
        uses: amondnet/vercel-action@v20
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
          vercel-project-id: ${{ secrets.VERCEL_PROJECT_ID }}
          working-directory: ./tealinux-fe

  deploy-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy to server
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            cd /home/tealinux/web-tealinux-astro
            git pull origin main
            cd tealinuxbe
            go build -o tealinux-api cmd/main.go
            sudo systemctl restart tealinux-api
```

---

## 📊 Monitoring

### Setup Logging

```bash
# View backend logs
sudo journalctl -u tealinux-api -f

# View Nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### Health Check Endpoint

Add to backend:

```go
app.Get("/health", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "healthy",
        "timestamp": time.Now(),
    })
})
```

### Uptime Monitoring

Use services like:
- UptimeRobot
- Pingdom
- StatusCake

---

## 🐛 Troubleshooting

### Backend Not Starting

```bash
# Check logs
sudo journalctl -u tealinux-api -n 50

# Check if port is in use
sudo lsof -i :3000

# Test database connection
psql -h localhost -U tealinux_user -d tealinux
```

### Frontend Build Fails

```bash
# Clear cache
rm -rf node_modules .astro dist
npm install
npm run build
```

### SSL Certificate Issues

```bash
# Renew certificate
sudo certbot renew

# Check certificate status
sudo certbot certificates
```

---

## 📝 Post-Deployment Checklist

- [ ] Frontend accessible via HTTPS
- [ ] Backend API responding
- [ ] Database connected
- [ ] OAuth login working
- [ ] Admin panel accessible
- [ ] SSL certificate valid
- [ ] Monitoring setup
- [ ] Backups configured
- [ ] DNS records correct
- [ ] Environment variables secure

---

## 🔒 Security Best Practices

1. **Use strong passwords** for database
2. **Enable firewall** (ufw on Ubuntu)
3. **Regular updates** for system packages
4. **Backup database** regularly
5. **Monitor logs** for suspicious activity
6. **Use environment variables** for secrets
7. **Enable HTTPS** everywhere
8. **Implement rate limiting**

---

## 📚 Additional Resources

- [Vercel Documentation](https://vercel.com/docs)
- [DigitalOcean Tutorials](https://www.digitalocean.com/community/tutorials)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Let's Encrypt](https://letsencrypt.org/)

---

Made with 🍵 by DOSCOM
