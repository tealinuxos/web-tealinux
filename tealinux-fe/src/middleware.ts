import { defineMiddleware } from 'astro:middleware';

/**
 * Decode JWT payload without signature verification.
 * This is safe for UI route protection because:
 * 1. All actual API calls are verified by the Go backend's JWT middleware
 * 2. This middleware only decides whether to show the admin UI or redirect to login
 * 3. Even if someone forges a token to see the admin page, they can't perform 
 *    any actions because every API call is verified server-side by Go
 */
function decodeJWTPayload(token: string): { id?: number; userId?: number; role?: string; email?: string; exp?: number } | null {
    try {
        const parts = token.split('.');
        if (parts.length !== 3) return null;

        // Base64url decode the payload
        let base64 = parts[1].replace(/-/g, '+').replace(/_/g, '/');
        while (base64.length % 4) base64 += '=';
        const binary = atob(base64);
        const bytes = new Uint8Array(binary.length);
        for (let i = 0; i < binary.length; i++) {
            bytes[i] = binary.charCodeAt(i);
        }
        const jsonStr = new TextDecoder().decode(bytes);
        return JSON.parse(jsonStr);
    } catch {
        return null;
    }
}

function extractBearerToken(authHeader: string): string | null {
    const parts = authHeader.split(' ');
    if (parts.length !== 2 || parts[0] !== 'Bearer') return null;
    return parts[1];
}

export const onRequest = defineMiddleware(async (context, next) => {
    const { url, cookies, redirect, request } = context;

    // Check if the request is for an admin page
    if (url.pathname.startsWith('/admin')) {
        // Try to get token from HttpOnly cookie first (preferred)
        let token = cookies.get('tealinux_access_token')?.value;

        // If not in cookie, try Authorization header
        if (!token) {
            const authHeader = request.headers.get('Authorization');
            token = authHeader ? extractBearerToken(authHeader) : null;
        }

        // No token found
        if (!token) {
            console.log('[Auth] No token found, redirecting to login');
            return redirect('/login?redirect=' + encodeURIComponent(url.pathname));
        }

        // Decode JWT payload (without signature verification - Go backend handles security)
        const payload = decodeJWTPayload(token);

        if (!payload) {
            console.log('[Auth] Failed to decode token');
            cookies.delete('tealinux_access_token', { path: '/' });
            cookies.delete('tealinux_refresh_token', { path: '/' });
            return redirect('/login?error=invalid');
        }

        // Check token expiration
        const now = Math.floor(Date.now() / 1000);
        if (payload.exp && payload.exp < now) {
            console.log('[Auth] Token expired');
            cookies.delete('tealinux_access_token', { path: '/' });
            cookies.delete('tealinux_refresh_token', { path: '/' });
            return redirect('/login?error=expired&redirect=' + encodeURIComponent(url.pathname));
        }

        // Check if user is admin
        if (payload.role !== 'admin') {
            console.log('[Auth] User is not admin, access denied');
            return redirect('/?error=forbidden');
        }

        // Attach user info to context
        context.locals.user = {
            userId: payload.userId || payload.id || 0,
            email: payload.email || '',
            role: payload.role || '',
            name: ''
        };

        console.log(`[Auth] Admin access granted for user ${payload.id || payload.userId}`);
    }

    return next();
});
