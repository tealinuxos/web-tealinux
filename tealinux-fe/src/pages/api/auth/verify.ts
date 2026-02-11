import type { APIRoute } from 'astro';
import { verifyToken } from '../../../lib/jwt';

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

        // Verify token
        const payload = await verifyToken(token);

        return new Response(JSON.stringify({
            success: true,
            authenticated: true,
            user: {
                userId: payload.userId,
                email: payload.email,
                name: payload.name,
                role: payload.role
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
