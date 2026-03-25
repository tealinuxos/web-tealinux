// JWT utility functions using Web Crypto API (no external dependencies)

const JWT_SECRET = process.env.JWT_SECRET || import.meta.env.JWT_SECRET || 'tealinux-secret-key-change-in-production';
const JWT_EXPIRES_IN = 7 * 24 * 60 * 60; // 7 days in seconds
const REFRESH_EXPIRES_IN = 30 * 24 * 60 * 60; // 30 days in seconds

export interface JWTPayload {
    userId: number;
    email: string;
    role: string;
    name: string;
}

interface JWTHeader {
    alg: string;
    typ: string;
}

interface JWTFullPayload extends JWTPayload {
    iat: number;
    exp: number;
    iss: string;
}

// Base64url encode
function base64urlEncode(data: string): string {
    const encoder = new TextEncoder();
    const bytes = encoder.encode(data);
    let binary = '';
    for (const byte of bytes) {
        binary += String.fromCharCode(byte);
    }
    return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

// Base64url decode
function base64urlDecode(str: string): string {
    let base64 = str.replace(/-/g, '+').replace(/_/g, '/');
    while (base64.length % 4) base64 += '=';
    const binary = atob(base64);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
        bytes[i] = binary.charCodeAt(i);
    }
    return new TextDecoder().decode(bytes);
}

// Create HMAC-SHA256 signature
async function createSignature(input: string, secret: string): Promise<string> {
    const encoder = new TextEncoder();
    const key = await crypto.subtle.importKey(
        'raw',
        encoder.encode(secret),
        { name: 'HMAC', hash: 'SHA-256' },
        false,
        ['sign']
    );
    const signature = await crypto.subtle.sign('HMAC', key, encoder.encode(input));
    const bytes = new Uint8Array(signature);
    let binary = '';
    for (const byte of bytes) {
        binary += String.fromCharCode(byte);
    }
    return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

// Verify HMAC-SHA256 signature
async function verifySignature(input: string, signature: string, secret: string): Promise<boolean> {
    const expected = await createSignature(input, secret);
    return expected === signature;
}

/**
 * Generate JWT token
 */
export async function generateAccessToken(payload: JWTPayload): Promise<string> {
    const header: JWTHeader = { alg: 'HS256', typ: 'JWT' };
    const now = Math.floor(Date.now() / 1000);

    const fullPayload: JWTFullPayload = {
        ...payload,
        iat: now,
        exp: now + JWT_EXPIRES_IN,
        iss: 'tealinuxos'
    };

    const headerB64 = base64urlEncode(JSON.stringify(header));
    const payloadB64 = base64urlEncode(JSON.stringify(fullPayload));
    const signature = await createSignature(`${headerB64}.${payloadB64}`, JWT_SECRET);

    return `${headerB64}.${payloadB64}.${signature}`;
}

/**
 * Generate refresh token
 */
export async function generateRefreshToken(payload: JWTPayload): Promise<string> {
    const header: JWTHeader = { alg: 'HS256', typ: 'JWT' };
    const now = Math.floor(Date.now() / 1000);

    const fullPayload: JWTFullPayload = {
        ...payload,
        iat: now,
        exp: now + REFRESH_EXPIRES_IN,
        iss: 'tealinuxos'
    };

    const headerB64 = base64urlEncode(JSON.stringify(header));
    const payloadB64 = base64urlEncode(JSON.stringify(fullPayload));
    const signature = await createSignature(`${headerB64}.${payloadB64}`, JWT_SECRET);

    return `${headerB64}.${payloadB64}.${signature}`;
}

/**
 * Generate both tokens
 */
export async function generateTokenPair(payload: JWTPayload) {
    const [accessToken, refreshToken] = await Promise.all([
        generateAccessToken(payload),
        generateRefreshToken(payload)
    ]);
    return { accessToken, refreshToken };
}

/**
 * Verify and decode JWT token
 */
export async function verifyToken(token: string): Promise<JWTPayload> {
    const parts = token.split('.');
    if (parts.length !== 3) {
        throw new Error('INVALID_TOKEN');
    }

    const [headerB64, payloadB64, signature] = parts;

    // Verify signature
    const valid = await verifySignature(`${headerB64}.${payloadB64}`, signature, JWT_SECRET);
    if (!valid) {
        throw new Error('INVALID_TOKEN');
    }

    // Decode payload
    const payload: any = JSON.parse(base64urlDecode(payloadB64));

    // Check expiration
    const now = Math.floor(Date.now() / 1000);
    if (payload.exp && payload.exp < now) {
        throw new Error('TOKEN_EXPIRED');
    }

    // Handle both Go backend format (id, role) and Astro format (userId, email, role, name)
    return {
        userId: payload.userId || payload.id || 0,
        email: payload.email || '',
        role: payload.role || '',
        name: payload.name || ''
    };
}

/**
 * Check if user has admin role
 */
export function isAdminRole(payload: JWTPayload): boolean {
    return payload.role === 'admin';
}

/**
 * Extract Bearer token from Authorization header
 */
export function extractBearerToken(authHeader: string): string | null {
    const parts = authHeader.split(' ');
    if (parts.length !== 2 || parts[0] !== 'Bearer') return null;
    return parts[1];
}
