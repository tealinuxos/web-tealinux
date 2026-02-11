# Sistem Autentikasi JWT - TeaLinuxOS

Sistem autentikasi lengkap menggunakan JSON Web Token (JWT) untuk melindungi halaman admin.

## 📋 Fitur

- ✅ Login dengan JWT (Access Token + Refresh Token)
- ✅ HttpOnly Cookies untuk keamanan maksimal
- ✅ Server-side middleware protection
- ✅ Client-side verification
- ✅ Role-based access control (Admin only)
- ✅ Token expiration handling
- ✅ Proper error handling (401, 403)
- ✅ Secure logout mechanism

## 🔐 Flow Autentikasi

```
┌─────────────┐
│   User      │
└──────┬──────┘
       │
       │ 1. POST /api/auth/login
       │    { email, password }
       ▼
┌─────────────────────┐
│  Login API          │
│  - Validate creds   │
│  - Generate JWT     │
│  - Set HttpOnly     │
│    cookies          │
└──────┬──────────────┘
       │
       │ 2. JWT Tokens
       │    (in cookies)
       ▼
┌─────────────────────┐
│  Client             │
│  - Redirect to      │
│    /admin or /      │
└──────┬──────────────┘
       │
       │ 3. Access /admin
       ▼
┌─────────────────────┐
│  Middleware         │
│  - Extract token    │
│  - Verify JWT       │
│  - Check role       │
└──────┬──────────────┘
       │
       │ 4a. Valid + Admin
       ▼
┌─────────────────────┐
│  Admin Dashboard    │
└─────────────────────┘
       │
       │ 4b. Invalid/Expired
       ▼
┌─────────────────────┐
│  Redirect to Login  │
└─────────────────────┘
```

## 📁 Struktur File

```
src/
├── lib/
│   ├── jwt.ts                    # JWT utilities (generate, verify)
│   └── auth.ts                   # Client-side auth helpers
├── middleware.ts                 # Server-side route protection
├── env.d.ts                      # TypeScript definitions
├── pages/
│   ├── login.astro              # Login page
│   ├── admin/
│   │   └── index.astro          # Protected admin page
│   └── api/
│       └── auth/
│           ├── login.ts         # POST /api/auth/login
│           ├── logout.ts        # POST /api/auth/logout
│           └── verify.ts        # GET /api/auth/verify
└── layouts/
    └── AdminLayout.astro        # Admin layout with auth check
```

## 🔑 API Endpoints

### 1. Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "admin@tealinuxos.org",
  "password": "admin123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "user": {
    "id": 1,
    "email": "admin@tealinuxos.org",
    "name": "Admin TeaLinux",
    "role": "admin"
  },
  "tokens": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

**Cookies Set:**
- `tealinux_access_token` (HttpOnly, 7 days)
- `tealinux_refresh_token` (HttpOnly, 30 days)

**Error Responses:**
- `400` - Missing email/password
- `401` - Invalid credentials

### 2. Logout
```http
POST /api/auth/logout
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

**Effect:** Clears all auth cookies

### 3. Verify Token
```http
GET /api/auth/verify
```

**Response (200 OK):**
```json
{
  "success": true,
  "authenticated": true,
  "user": {
    "userId": 1,
    "email": "admin@tealinuxos.org",
    "name": "Admin TeaLinux",
    "role": "admin"
  }
}
```

**Error Response (401):**
```json
{
  "success": false,
  "authenticated": false,
  "error": "TOKEN_EXPIRED"
}
```

## 🛡️ Middleware Protection

File: `src/middleware.ts`

**Protects:** All `/admin/*` routes

**Checks:**
1. Token exists (HttpOnly cookie or Authorization header)
2. Token is valid (not expired, correct signature)
3. User has admin role

**Redirects:**
- No token → `/login?redirect=/admin`
- Token expired → `/login?error=expired&redirect=/admin`
- Invalid token → `/login?error=invalid`
- Not admin → `/?error=forbidden`

## 👤 Test Accounts

### Admin Account
```
Email: admin@tealinuxos.org
Password: admin123
Role: admin
```

### Regular User
```
Email: user@tealinuxos.org
Password: user123
Role: user
```

## 🔒 Security Features

### 1. HttpOnly Cookies
- Tokens stored in HttpOnly cookies
- Not accessible via JavaScript
- Prevents XSS attacks

### 2. JWT Verification
- Signature verification with secret key
- Expiration checking
- Issuer and audience validation

### 3. Role-Based Access
- Admin role required for `/admin` routes
- Enforced on both server and client side

### 4. HTTPS in Production
- Secure flag on cookies in production
- Prevents man-in-the-middle attacks

### 5. SameSite Protection
- SameSite=Lax on cookies
- Prevents CSRF attacks

## 🚀 Usage Examples

### Login Flow
```typescript
// Client-side login
const response = await fetch('/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'admin@tealinuxos.org',
    password: 'admin123'
  })
});

const result = await response.json();

if (result.success) {
  // Cookies are automatically set
  // Redirect to admin
  window.location.href = '/admin';
}
```

### Logout Flow
```typescript
// Client-side logout
await fetch('/api/auth/logout', { method: 'POST' });
window.location.href = '/login';
```

### Check Authentication
```typescript
// Verify current session
const response = await fetch('/api/auth/verify');
const result = await response.json();

if (result.authenticated && result.user.role === 'admin') {
  // User is authenticated admin
}
```

## ⚙️ Configuration

### JWT Secret
Set in environment variable (production):
```env
JWT_SECRET=your-super-secret-key-here
```

Default (development):
```typescript
const JWT_SECRET = 'tealinux-secret-key-change-in-production';
```

### Token Expiration
```typescript
const JWT_EXPIRES_IN = '7d';          // Access token: 7 days
const REFRESH_TOKEN_EXPIRES_IN = '30d'; // Refresh token: 30 days
```

## 🐛 Error Handling

### Token Expired
```
URL: /login?error=expired&redirect=/admin
Message: "Your session has expired. Please login again."
```

### Invalid Token
```
URL: /login?error=invalid
Message: "Invalid authentication token. Please login again."
```

### Unauthorized
```
URL: /login?error=unauthorized
Message: "Unauthorized access. Please login."
```

### Forbidden (Not Admin)
```
URL: /?error=forbidden
Message: "You do not have permission to access that page."
```

## 📝 TODO / Production Checklist

- [ ] Install `jsonwebtoken` package: `bun add jsonwebtoken @types/jsonwebtoken`
- [ ] Replace mock user database with real database
- [ ] Use bcrypt for password hashing
- [ ] Set strong JWT_SECRET in environment variables
- [ ] Enable HTTPS in production
- [ ] Implement refresh token rotation
- [ ] Add rate limiting on login endpoint
- [ ] Add CAPTCHA for brute force protection
- [ ] Implement password reset flow
- [ ] Add audit logging for admin actions
- [ ] Set up token blacklist for logout
- [ ] Implement 2FA (Two-Factor Authentication)

## 🔧 Development

### Install Dependencies
```bash
cd tealinux-fe
bun add jsonwebtoken @types/jsonwebtoken cookie
```

### Run Development Server
```bash
bun run dev
```

### Test Login
1. Navigate to `http://localhost:4321/login`
2. Use test credentials (see Test Accounts section)
3. Should redirect to `/admin` for admin user
4. Should redirect to `/` for regular user

### Test Protection
1. Try accessing `http://localhost:4321/admin` without login
2. Should redirect to `/login?redirect=/admin`
3. After login, should redirect back to `/admin`

## 📚 References

- [JWT.io](https://jwt.io/) - JWT debugger and documentation
- [Astro Middleware](https://docs.astro.build/en/guides/middleware/)
- [HttpOnly Cookies](https://owasp.org/www-community/HttpOnly)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
