# 🍵 TeaLinuxOS Web Platform

<div align="center">

![TeaLinuxOS](https://img.shields.io/badge/TeaLinuxOS-54CD4C?style=for-the-badge&logo=linux&logoColor=black)
![Astro](https://img.shields.io/badge/Astro-FF5D01?style=for-the-badge&logo=astro&logoColor=white)
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)

**Platform web modern untuk distribusi Linux TeaLinuxOS**

[Demo](#) • [Dokumentasi](#) • [Kontribusi](#kontribusi)

</div>

---

---

## 🎯 Tentang Proyek

**TeaLinuxOS Web Platform** adalah aplikasi web full-stack yang dirancang untuk mendukung distribusi Linux TeaLinuxOS. Platform ini menyediakan antarmuka modern untuk download ISO, manajemen pengguna, dan administrasi konten.





---


### Frontend Stack
| Teknologi | Versi | Deskripsi |
|-----------|-------|-----------|
| **Astro** | 5.16.11 | Framework web modern untuk performa optimal |
| **Tailwind CSS** | 4.1.18 | Utility-first CSS framework |
| **GSAP** | 3.14.2 | Library animasi profesional |
| **Lenis** | 1.3.17 | Smooth scrolling library |

### Backend Stack
| Teknologi | Versi | Deskripsi |
|-----------|-------|-----------|
| **Go** | 1.25.5 | Bahasa pemrograman backend |
| **Fiber** | 2.52.10 | Web framework untuk Go |
| **GORM** | 1.31.1 | ORM untuk Go |
| **PostgreSQL** | 16 | Database relasional |
| **JWT** | 5.3.0 | JSON Web Token untuk autentikasi |
| **OAuth2** | 0.34.0 | Autentikasi dengan provider eksternal |

### DevOps & Infrastructure
| Teknologi | Deskripsi |
|-----------|-----------|
| **Docker** | Containerization untuk semua services |
| **Docker Compose** | Orchestration untuk development & production |
| **Nginx** | Reverse proxy & load balancer |
| **Let's Encrypt** | SSL/TLS certificates (Certbot) |
| **GitHub Actions** | CI/CD pipeline automation |

---

##  Arsitektur

```
┌─────────────────────────────────────────────────────────┐
│                    Client Browser                        │
│                  (Astro Frontend)                        │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ HTTP/HTTPS
                     │
┌────────────────────▼────────────────────────────────────┐
│                  Go Fiber Backend                        │
│  ┌──────────────────────────────────────────────────┐   │
│  │              Routes & Handlers                   │   │
│  │  • Auth (Login, Register, OAuth)                 │   │
│  │  • Users Management                              │   │
│  │  • Categories Management                         │   │
│  │  • Admin Operations                              │   │
│  └────────────┬─────────────────────────────────────┘   │
│               │                                          │
│  ┌────────────▼─────────────────────────────────────┐   │
│  │              Middleware Layer                    │   │
│  │  • JWT Verification                              │   │
│  │  • CORS Handler                                  │   │
│  │  • Role-Based Access Control                     │   │
│  └────────────┬─────────────────────────────────────┘   │
│               │                                          │
│  ┌────────────▼─────────────────────────────────────┐   │
│  │              Business Logic                      │   │
│  │  • User Service                                  │   │
│  │  • Category Service                              │   │
│  │  • Auth Service                                  │   │
│  └────────────┬─────────────────────────────────────┘   │
└───────────────┼──────────────────────────────────────────┘
                │
                │ GORM ORM
                │
┌───────────────▼──────────────────────────────────────────┐
│              PostgreSQL Database                         │
│  • users                                                 │
│  • categories                                            │
│  • sessions                                              │
└──────────────────────────────────────────────────────────┘
```

---

##  Instalasi

### Prasyarat

Pastikan Anda telah menginstal:
- **Node.js** (v18 atau lebih tinggi)
- **npm** atau **bun**
- **Go** (v1.25 atau lebih tinggi)
- **Docker** dan **Docker Compose**
- **Git**

### 1. Clone Repository

```bash
git clone https://github.com/yourusername/web-tealinux-astro.git
cd web-tealinux-astro
```

### 2. Setup Database

Jalankan PostgreSQL menggunakan Docker Compose:

```bash
docker-compose up -d
```

Database akan berjalan di `localhost:5432` dengan kredensial:
- **Database**: `tealinux`
- **User**: `tealinux_user`
- **Password**: `tealinux123`

### 3. Setup Backend

```bash
cd tealinuxbe

# Install dependencies
go mod download

# Copy environment file (sesuaikan dengan konfigurasi Anda)
cp .env.example .env

# Edit .env file dengan konfigurasi yang sesuai
nano .env

# Jalankan backend
go run cmd/main.go
```

Backend akan berjalan di `http://localhost:3000`

### 4. Setup Frontend

```bash
cd ../tealinux-fe

# Install dependencies
npm install
# atau jika menggunakan bun
bun install

# Jalankan development server
npm run dev
# atau
bun run dev
```

Frontend akan berjalan di `http://localhost:4321`

---

## 🚀 Penggunaan

### Development Mode

1. **Start Database**
   ```bash
   docker-compose up -d
   ```

2. **Start Backend**
   ```bash
   cd tealinuxbe
   go run cmd/main.go
   ```

3. **Start Frontend**
   ```bash
   cd tealinux-fe
   npm run dev
   ```

4. **Akses Aplikasi**
   - Frontend: http://localhost:4321
   - Backend API: http://localhost:3000
   - Admin Panel: http://localhost:4321/admin

### Production Build

#### Frontend
```bash
cd tealinux-fe
npm run build
npm run preview
```

#### Backend
```bash
cd tealinuxbe
go build -o tealinux-api cmd/main.go
./tealinux-api
```

---

## 🚢 Deployment & CI/CD

### Quick Deploy (Development with Docker)

```bash
# Start all services with hot-reload
docker compose -f docker-compose.dev.yml up -d

# View logs
docker compose -f docker-compose.dev.yml logs -f

# Stop all services
docker compose -f docker-compose.dev.yml down
```

### Production Deployment to VPS

**Automated CI/CD Pipeline:**
- ✅ Automatic testing on push
- ✅ Docker image building to GitHub Container Registry
- ✅ Automated deployment to VPS via SSH
- ✅ Health checks and rollback support

**Quick Deploy Commands:**
```bash
# On VPS - Quick deploy
./deploy.sh

# Health check
./health-check.sh

# Database backup
./backup-db.sh

# SSL setup (first time)
./setup-ssl.sh
```

**📚 Detailed Guides:**
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Complete deployment guide (VPS setup, SSL, monitoring)
- **[QUICK_REFERENCE.md](./QUICK_REFERENCE.md)** - Quick command reference

**GitHub Secrets Required:**
| Secret | Description |
|--------|-------------|
| `VPS_HOST` | VPS IP address |
| `VPS_USERNAME` | SSH username |
| `VPS_SSH_KEY` | Private SSH key |
| `VPS_PROJECT_PATH` | Project path (e.g., `/opt/tealinux`) |

**CI/CD Pipeline Stages:**
1. 🧪 Test Frontend (Astro build)
2. 🧪 Test Backend (Go tests + build)
3. 🐳 Build & Push Docker images to GHCR
4. 🚀 Deploy to VPS via SSH
5. ✅ Health check verification

---

## 📁 Struktur Proyek

```
web-tealinux-astro/
├── tealinux-fe/                 # Frontend Astro
│   ├── src/
│   │   ├── assets/              # Gambar, font, dll
│   │   ├── components/          # Komponen Astro
│   │   │   ├── atoms/           # Komponen kecil
│   │   │   ├── molecules/       # Komponen menengah
│   │   │   └── organisms/       # Komponen besar
│   │   ├── layouts/             # Layout template
│   │   ├── lib/                 # Utilities dan helpers
│   │   ├── pages/               # Halaman aplikasi
│   │   │   ├── admin/           # Admin pages
│   │   │   ├── index.astro      # Landing page
│   │   │   ├── download.astro   # Download page
│   │   │   ├── login.astro      # Login page
│   │   │   └── register.astro   # Register page
│   │   └── styles/              # Global styles
│   ├── public/                  # Static assets
│   ├── package.json
│   └── astro.config.mjs
│
├── tealinuxbe/                  # Backend Go
│   ├── cmd/
│   │   └── main.go              # Entry point
│   ├── internal/
│   │   ├── config/              # Konfigurasi
│   │   ├── database/            # Database connection
│   │   ├── handlers/            # HTTP handlers
│   │   │   ├── auth.go          # Auth handlers
│   │   │   ├── user.go          # User handlers
│   │   │   ├── category.go      # Category handlers
│   │   │   └── admin.go         # Admin handlers
│   │   ├── middleware/          # Middleware
│   │   │   ├── cors.go          # CORS middleware
│   │   │   ├── jwt.go           # JWT middleware
│   │   │   └── role.go          # Role middleware
│   │   ├── models/              # Database models
│   │   │   ├── user.go
│   │   │   └── category.go
│   │   ├── routes/              # Route definitions
│   │   └── utils/               # Utilities
│   │       ├── jwt.go           # JWT utilities
│   │       └── hash.go          # Password hashing
│   ├── go.mod
│   └── go.sum
│
├── docker-compose.yml           # Docker configuration
├── package.json                 # Root package.json
└── README.md                    # Dokumentasi ini
```

---





``
## 🔐 Environment Variables

### Backend (.env)
```env
# Server
PORT=3000
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=tealinux_user
DB_PASSWORD=tealinux123
DB_NAME=tealinux

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRE=24h

# OAuth Google
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:3000/api/auth/google/callback

# OAuth GitHub
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=http://localhost:3000/api/auth/github/callback

# Frontend URL
FRONTEND_URL=http://localhost:4321
```

---

## 🧪 Testing

### Backend Testing
```bash
cd tealinuxbe
go test ./...
```

### Frontend Testing
```bash
cd tealinux-fe
npm run test
```

---

## 🤝 Kontribusi

Kami sangat menghargai kontribusi dari komunitas! Berikut cara berkontribusi:

1. **Fork** repository ini
2. **Clone** fork Anda
   ```bash
   git clone https://github.com/yourusername/web-tealinux-astro.git
   ```
3. **Buat branch** untuk fitur baru
   ```bash
   git checkout -b feature/amazing-feature
   ```
4. **Commit** perubahan Anda
   ```bash
   git commit -m 'Add some amazing feature'
   ```
5. **Push** ke branch
   ```bash
   git push origin feature/amazing-feature
   ```
6. **Buat Pull Request**

### Coding Standards

- **Frontend**: Ikuti Astro best practices dan ESLint rules
- **Backend**: Ikuti Go conventions dan gofmt
- **Commit Messages**: Gunakan conventional commits
  - `feat:` untuk fitur baru
  - `fix:` untuk bug fixes
  - `docs:` untuk dokumentasi
  - `style:` untuk formatting
  - `refactor:` untuk refactoring
  - `test:` untuk testing
