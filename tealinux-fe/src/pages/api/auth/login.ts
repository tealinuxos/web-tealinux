import type { APIRoute } from 'astro';
import { generateTokenPair, type JWTPayload } from '../../../lib/jwt';

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

        // TODO: Replace with actual database check
        // This is a mock authentication - replace with real user verification
        const mockUsers = [
            {
                id: 1,
                email: 'admin@tealinuxos.org',
                password: 'admin123', // In production, use hashed passwords!
                name: 'Admin TeaLinux',
                role: 'admin'
            },
            {
                id: 2,
                email: 'user@tealinuxos.org',
                password: 'user123',
                name: 'Regular User',
                role: 'user'
            }
        ];

        const user = mockUsers.find(u => u.email === email && u.password === password);

        if (!user) {
            return new Response(JSON.stringify({
                success: false,
                error: 'Invalid email or password'
            }), {
                status: 401,
                headers: { 'Content-Type': 'application/json' }
            });
        }

        // Generate JWT tokens
        const payload: JWTPayload = {
            userId: user.id,
            email: user.email,
            name: user.name,
            role: user.role
        };

        const { accessToken, refreshToken } = await generateTokenPair(payload);

        // Set HttpOnly cookies (more secure than localStorage)
        cookies.set('tealinux_access_token', accessToken, {
            httpOnly: true,
            secure: import.meta.env.PROD, // Only HTTPS in production
            sameSite: 'lax',
            path: '/',
            maxAge: 60 * 60 * 24 * 7 // 7 days
        });

        cookies.set('tealinux_refresh_token', refreshToken, {
            httpOnly: true,
            secure: import.meta.env.PROD,
            sameSite: 'lax',
            path: '/',
            maxAge: 60 * 60 * 24 * 30 // 30 days
        });

        // Return user info (without password) and tokens
        return new Response(JSON.stringify({
            success: true,
            user: {
                id: user.id,
                email: user.email,
                name: user.name,
                role: user.role
            },
            tokens: {
                accessToken,
                refreshToken
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
