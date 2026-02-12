import type { APIRoute } from 'astro';

const BACKEND_URL = import.meta.env.BACKEND_URL || 'http://tealinux-backend:3000';

export const POST: APIRoute = async ({ request, cookies }) => {
    try {
        const body = await request.json();
        const { email, password } = body;

        // Validate input
        if (!email || !password) {
            return new Response(JSON.stringify({
                success: false,
                error: 'Email and password are required'
            }), {
                status: 400,
                headers: { 'Content-Type': 'application/json' }
            });
        }

        // Forward request to Go backend
        const backendResponse = await fetch(`${BACKEND_URL}/api/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        });

        const data = await backendResponse.json();

        if (!backendResponse.ok) {
            return new Response(JSON.stringify({
                success: false,
                error: data.error || 'Invalid email or password'
            }), {
                status: backendResponse.status,
                headers: { 'Content-Type': 'application/json' }
            });
        }

        // Backend returns: { access_token, refresh_token, user: { id, name, email, role, avatar } }
        const { access_token, refresh_token, user } = data;

        // Set HttpOnly cookies (more secure than localStorage)
        cookies.set('tealinux_access_token', access_token, {
            httpOnly: true,
            secure: import.meta.env.PROD,
            sameSite: 'lax',
            path: '/',
            maxAge: 60 * 60 * 24 * 7 // 7 days
        });

        cookies.set('tealinux_refresh_token', refresh_token, {
            httpOnly: true,
            secure: import.meta.env.PROD,
            sameSite: 'lax',
            path: '/',
            maxAge: 60 * 60 * 24 * 30 // 30 days
        });

        // Return user info and tokens
        return new Response(JSON.stringify({
            success: true,
            user: {
                id: user.id,
                email: user.email,
                name: user.name,
                role: user.role,
                avatar: user.avatar
            },
            tokens: {
                accessToken: access_token,
                refreshToken: refresh_token
            }
        }), {
            status: 200,
            headers: { 'Content-Type': 'application/json' }
        });

    } catch (error) {
        console.error('[Login API] Error:', error);
        return new Response(JSON.stringify({
            success: false,
            error: 'Internal server error'
        }), {
            status: 500,
            headers: { 'Content-Type': 'application/json' }
        });
    }
};
