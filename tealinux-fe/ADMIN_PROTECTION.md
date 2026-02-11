# Admin Route Protection

The `/admin` path is now protected with both server-side and client-side authentication.

## How It Works

### Server-Side Protection (Middleware)
- **File**: `src/middleware.ts`
- Intercepts all requests to `/admin/*` routes
- Checks for `tealinux_token` and `tealinux_user` cookies
- Verifies that the user has `role: 'admin'`
- Redirects unauthorized users:
  - No token → `/login`
  - Not admin → `/` (home page)

### Client-Side Protection (AdminLayout)
- **File**: `src/layouts/AdminLayout.astro`
- Double-checks authentication after page load
- Redirects if user is not authenticated or not admin
- Provides a seamless user experience

## Authentication Flow

1. **Login**: User logs in via `/login` page
2. **Save Credentials**: Auth data is saved to:
   - `localStorage` (for client-side access)
   - `cookies` (for server-side middleware)
3. **Access Admin**: When accessing `/admin`:
   - Middleware checks cookies first (server-side)
   - If authorized, page loads
   - AdminLayout double-checks (client-side)
4. **Logout**: Clears both localStorage and cookies

## Testing

To test the protection:

1. **Without login**: Navigate to `http://localhost:4321/admin`
   - Should redirect to `/login`

2. **With non-admin user**: Login with a regular user account
   - Should redirect to `/` (home page)

3. **With admin user**: Login with admin credentials
   - Should access the admin dashboard successfully

## Files Modified

- `src/middleware.ts` - New server-side protection
- `src/lib/auth.ts` - Updated to save/clear cookies
- `src/layouts/AdminLayout.astro` - Already had client-side protection
