import { defineMiddleware } from 'astro:middleware';
import { verifyToken, isAdminRole, extractBearerToken } from './lib/jwt';

export const onRequest = defineMiddleware(async (context, next) => {
    const { url, cookies, redirect, request } = context;

    // Check if the request is for an admin page
    if (url.pathname.startsWith('/admin')) {
        try {
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

            // Verify JWT token
            const payload = await verifyToken(token);

            // Check if user is admin
            if (!isAdminRole(payload)) {
                console.log('[Auth] User is not admin, access denied');
                return redirect('/?error=forbidden');
            }

            // Token is valid and user is admin, attach user info to context
            context.locals.user = payload;

            console.log(`[Auth] Admin access granted for ${payload.email}`);

        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : 'Unknown error';
            console.error('[Auth] Token verification failed:', errorMessage);

            // Clear invalid token
            cookies.delete('tealinux_access_token', { path: '/' });
            cookies.delete('tealinux_refresh_token', { path: '/' });

            // Handle specific errors
            if (errorMessage === 'TOKEN_EXPIRED') {
                return redirect('/login?error=expired&redirect=' + encodeURIComponent(url.pathname));
            } else if (errorMessage === 'INVALID_TOKEN') {
                return redirect('/login?error=invalid');
            }

            return redirect('/login?error=unauthorized');
        }
    }

    return next();
});
