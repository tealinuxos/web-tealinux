================================================================================
PERBAIKAN ERROR: ERR_MODULE_NOT_FOUND - server-destroy
================================================================================

KRONOLOGI MASALAH:
------------------
1. ✅ Error 401 Fixed → Frontend sekarang proxy ke Go backend
2. ❌ Error baru muncul: Cannot find package 'server-destroy'
3. 🔍 Root cause: Astro SSR mode butuh adapter (@astrojs/node) yang memerlukan 
   package 'server-destroy', tapi package ini tidak ada di dependencies

SOLUSI YANG DITERAPKAN:
-----------------------

📦 1. package.json
   ✅ Tambah "@astrojs/node": "^9.1.1"
   ✅ Tambah "server-destroy": "^1.0.1"

⚙️ 2. astro.config.mjs
   ✅ Import adapter: import node from '@astrojs/node'
   ✅ Configure adapter: adapter: node({ mode: 'standalone' })

🐳 3. Dockerfile.prod
   ✅ Multi-stage build (builder + runtime)
   ✅ Builder: Install all deps + build app
   ✅ Runtime: Install all deps (bukan hanya production!)
   ✅ Copy dist dari builder stage

🐙 4. docker-compose.yml (sebelumnya)
   ✅ Tambah BACKEND_URL: http://tealinux-backend:3000
   ✅ Hapus env_file dari postgres service

================================================================================
LANGKAH DEPLOYMENT LENGKAP
================================================================================

# 1. Stop semua container
docker compose down

# 2. Hapus volume postgres lama (PENTING - untuk reset password!)
docker volume rm web-tealinux_postgres_data

# 3. Rebuild frontend dengan dependencies baru
docker compose build frontend --no-cache

# 4. Rebuild backend (optional, tapi direkomendasikan)
docker compose build backend --no-cache

# 5. Jalankan semua service
docker compose up -d

# 6. Monitor logs
docker compose logs -f

# Tunggu sampai muncul:
# - "Admin user seeded: admin@tealinux.org / admin123" (dari backend)
# - Frontend running tanpa error di port 4321

# 7. Test login
# Buka: https://tealinuxos.org/login
# Email: admin@tealinux.org
# Password: admin123

================================================================================
ARSITEKTUR AKHIR
================================================================================

┌─────────────────────────────────────────────────────────────────┐
│                         BROWSER                                  │
│                  https://tealinuxos.org                         │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                      TRAEFIK (Reverse Proxy)                     │
│   - /api/* → backend:3000                                       │
│   - /* → frontend:4321                                          │
└─────────────┬───────────────────────────────┬───────────────────┘
              │                               │
              ▼                               ▼
┌─────────────────────────┐   ┌───────────────────────────────────┐
│   FRONTEND (Astro SSR)  │   │    BACKEND (Go Fiber)             │
│   Port: 4321            │   │    Port: 3000                     │
│   Node: 20-alpine       │   │    ┌─────────────────────────┐    │
│                         │   │    │  /api/auth/login        │    │
│  /api/auth/login (Proxy)│───┼───→│  /api/auth/register     │    │
│         │               │   │    │  /api/downloads         │    │
│         ├─ @astrojs/node│   │    └─────────────────────────┘    │
│         └─ server-destroy│  │                │                  │
└─────────────────────────┘   └────────────────┼──────────────────┘
                                               │
                                               ▼
                              ┌────────────────────────────────────┐
                              │   POSTGRESQL                       │
                              │   Port: 5432                       │
                              │   User: tealinux_user              │
                              │   Pass: dionmulmedabadi            │
                              │   DB: tealinux                     │
                              └────────────────────────────────────┘

================================================================================
FILE YANG DIUBAH
================================================================================

1. tealinux-fe/src/pages/api/auth/login.ts
   - Dari: Mock users hardcoded
   - Ke: Proxy ke backend Go

2. tealinux-fe/package.json
   - Tambah: @astrojs/node, server-destroy

3. tealinux-fe/astro.config.mjs
   - Tambah: Node.js adapter configuration

4. tealinux-fe/Dockerfile.prod
   - Dari: Simple single-stage copy pre-built dist
   - Ke: Multi-stage build dengan full dependencies

5. docker-compose.yml
   - Frontend: Tambah BACKEND_URL env var
   - Postgres: Hapus env_file yang tidak perlu

================================================================================
TROUBLESHOOTING
================================================================================

🔴 Error: "Cannot find package 'server-destroy'"
   → Run: docker compose build frontend --no-cache

🔴 Error: "401 Unauthorized" saat login
   → Pastikan backend sudah running
   → Cek logs: docker compose logs backend
   → Pastikan user sudah di-seed

🔴 Error: "Connection refused to tealinux-backend"
   → Pastikan backend di network yang sama: docker network inspect tealinux-network
   → Restart services: docker compose restart

🔴 Error: "password authentication failed"
   → Hapus volume postgres: docker volume rm web-tealinux_postgres_data
   → Rebuild: docker compose up -d

🔴 Setelah masuk admin panel: Database error atau user not found
   → Berarti masih pakai Astro JWT token (lama)
   → Clear cookies browser
   → Login ulang dengan admin@tealinux.org

================================================================================
VERIFIKASI DEPLOYMENT
================================================================================

# 1. Cek semua container running
docker compose ps
# Expected: 3 containers (postgres, backend, frontend) - all "Up"

# 2. Cek backend logs
docker compose logs backend | tail -20
# Expected: "Database connected & migrated" dan "Admin user seeded"

# 3. Cek frontend logs
docker compose logs frontend | tail -20
# Expected: Running on 0.0.0.0:4321

# 4. Test koneksi internal
docker compose exec frontend ping -c 2 tealinux-backend
# Expected: Success

# 5. Test API langsung
docker compose exec frontend wget -qO- http://tealinux-backend:3000/api/health || echo "Setup health endpoint"

# 6. Test login dari browser
# Buka: https://tealinuxos.org/login
# Input: admin@tealinux.org / admin123
# Expected: Redirect ke /admin atau / dengan user terautentikasi

================================================================================
