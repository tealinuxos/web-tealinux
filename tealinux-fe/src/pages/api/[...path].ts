import type { APIRoute } from 'astro';

const BACKEND_URL = import.meta.env.BACKEND_URL || 'http://tealinux-backend:3000';

export const ALL: APIRoute = async ({ request, params }) => {
    const { path } = params;

    // Construct the backend URL
    // The 'path' param from [...path].ts will contain the rest of the URL
    const targetUrl = `${BACKEND_URL}/api/${path}`;

    console.log(`[API Proxy] Forwarding ${request.method} /api/${path} -> ${targetUrl}`);

    try {
        // Clone the request to forward it
        const headers = new Headers(request.headers);

        // Remove host header to avoid issues with backend
        headers.delete('host');

        const backendResponse = await fetch(targetUrl, {
            method: request.method,
            headers: headers,
            body: request.method !== 'GET' && request.method !== 'HEAD' ? await request.arrayBuffer() : undefined,
            redirect: 'manual'
        });

        // Return the backend's response as-is
        return new Response(backendResponse.body, {
            status: backendResponse.status,
            statusText: backendResponse.statusText,
            headers: backendResponse.headers
        });

    } catch (error) {
        console.error(`[API Proxy] Error forwarding to ${targetUrl}:`, error);
        return new Response(JSON.stringify({
            error: 'Backend connection failed',
            details: error instanceof Error ? error.message : String(error)
        }), {
            status: 502,
            headers: { 'Content-Type': 'application/json' }
        });
    }
};
