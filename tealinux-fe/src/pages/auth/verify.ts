import type { APIRoute } from 'astro';

/**
 * Decode JWT payload without signature verification.
 * Security is handled by the Go backend on actual API calls.
 */
function decodeJWTPayload(token: string): any | null {
    try {
        const parts = token.split('.');
        if (parts.length !== 3) return null;
        let base64 = parts[1].replace(/-/g, '+').replace(/_/g, '/');
        while (base64.length % 4) base64 += '=';
        const binary = atob(base64);
        const bytes = new Uint8Array(binary.length);
        for (let i = 0; i < binary.length; i++) {
            bytes[i] = binary.charCodeAt(i);
        }
        return JSON.parse(new TextDecoder().decode(bytes));
    } catch {
        return null;
    }
}

export const GET: APIRoute = async ({ cookies }) => {
    try {
        const token = cookies.get('tealinux_access_token')?.value;

        if (!token) {
            return new Response(JSON.stringify({
                success: false,
                authenticated: false,
                error: 'No token found'
            }), {
                status: 401,
                headers: { 'Content-Type': 'application/json' }
            });
        }

        const payload = decodeJWTPayload(token);

        if (!payload) {
            return new Response(JSON.stringify({
                success: false,
                authenticated: false,
                error: 'Invalid token format'
            }), {
                status: 401,
                headers: { 'Content-Type': 'application/json' }
            });
        }

        // Check expiration
        const now = Math.floor(Date.now() / 1000);
        if (payload.exp && payload.exp < now) {
            return new Response(JSON.stringify({
                success: false,
                authenticated: false,
                error: 'Token expired'
            }), {
                status: 401,
                headers: { 'Content-Type': 'application/json' }
            });
        }

        return new Response(JSON.stringify({
            success: true,
            authenticated: true,
            user: {
                userId: payload.userId || payload.id || 0,
                email: payload.email || '',
                name: payload.name || 'Admin',
                role: payload.role || ''
            }
        }), {
            status: 200,
            headers: { 'Content-Type': 'application/json' }
        });

    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Unknown error';

        return new Response(JSON.stringify({
            success: false,
            authenticated: false,
            error: errorMessage
        }), {
            status: 401,
            headers: { 'Content-Type': 'application/json' }
        });
    }
};
