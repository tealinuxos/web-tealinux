import type { APIRoute } from 'astro';

export const POST: APIRoute = async ({ cookies }) => {
    try {
        // Clear all auth cookies
        cookies.delete('tealinux_access_token', { path: '/' });
        cookies.delete('tealinux_refresh_token', { path: '/' });

        return new Response(JSON.stringify({
            success: true,
            message: 'Logged out successfully'
        }), {
            status: 200,
            headers: { 'Content-Type': 'application/json' }
        });

    } catch (error) {
        console.error('[Logout API] Error:', error);
        return new Response(JSON.stringify({
            success: false,
            error: 'Internal server error'
        }), {
            status: 500,
            headers: { 'Content-Type': 'application/json' }
        });
    }
};
