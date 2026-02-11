// Authentication helper functions for client-side

const TOKEN_KEY = 'tealinux_token';
const REFRESH_TOKEN_KEY = 'tealinux_refresh_token';
const USER_KEY = 'tealinux_user';

export interface User {
    id: number;
    name: string;
    email: string;
    role: string;
    avatar: string;
}

export interface AuthTokens {
    access_token: string;
    refresh_token: string;
}

// Token management
export function getToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem(TOKEN_KEY, token);
}

export function getRefreshToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(REFRESH_TOKEN_KEY);
}

export function setRefreshToken(token: string): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem(REFRESH_TOKEN_KEY, token);
}

export function removeTokens(): void {
    if (typeof window === 'undefined') return;
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
}

// User management
export function getCurrentUser(): User | null {
    if (typeof window === 'undefined') return null;
    const userStr = localStorage.getItem(USER_KEY);
    if (!userStr) return null;
    try {
        return JSON.parse(userStr);
    } catch {
        return null;
    }
}

export function setCurrentUser(user: User): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem(USER_KEY, JSON.stringify(user));
}

// Authentication status
export function isAuthenticated(): boolean {
    return !!getToken();
}

// Check if user is admin
export function isAdmin(): boolean {
    const user = getCurrentUser();
    return user !== null && user.role === 'admin';
}

// Save auth data (tokens + user)
export function saveAuthData(tokens: AuthTokens, user: User): void {
    setToken(tokens.access_token);
    setRefreshToken(tokens.refresh_token);
    setCurrentUser(user);

    // Also save to cookies for server-side middleware access
    if (typeof document !== 'undefined') {
        // Set cookies with 7 days expiry
        const expires = new Date();
        expires.setDate(expires.getDate() + 7);

        document.cookie = `tealinux_token=${tokens.access_token}; path=/; expires=${expires.toUTCString()}; SameSite=Lax`;
        document.cookie = `tealinux_user=${encodeURIComponent(JSON.stringify(user))}; path=/; expires=${expires.toUTCString()}; SameSite=Lax`;
    }
}

// Clear all auth data
export function clearAuthData(): void {
    removeTokens();

    // Also clear cookies
    if (typeof document !== 'undefined') {
        document.cookie = 'tealinux_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT';
        document.cookie = 'tealinux_user=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT';
    }
}
