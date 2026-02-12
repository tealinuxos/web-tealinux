import type { APIRoute } from 'astro';

const BACKEND_URL = import.meta.env.BACKEND_URL || 'http://tealinux-backend:3000';

export const POST: APIRoute = async ({ request }) => {
    try {
        const body = await request.json();
        const { name, email, password } = body;

        // Validate input
        if (!name || !email || !password) {
            return new Response(JSON.stringify({
                error: 'Name, email and password are required'
            }), {
                status: 400,
                headers: { 'Content-Type': 'application/json' }
            });
        }

        // Forward request to Go backend
        const backendResponse = await fetch(`${BACKEND_URL}/api/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name, email, password })
        });

        const data = await backendResponse.json();

        if (!backendResponse.ok) {
            return new Response(JSON.stringify({
                error: data.error || 'Registration failed'
            }), {
                status: backendResponse.status,
                headers: { 'Content-Type': 'application/json' }
            });
        }

        // Backend returns: { id, name, email, role, avatar }
        return new Response(JSON.stringify(data), {
            status: 200,
            headers: { 'Content-Type': 'application/json' }
        });

    } catch (error) {
        console.error('[Register API] Error:', error);
        return new Response(JSON.stringify({
            error: 'Internal server error'
        }), {
            status: 500,
            headers: { 'Content-Type': 'application/json' }
        });
    }
};
