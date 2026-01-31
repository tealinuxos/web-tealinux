import { API_BASE_URL } from './config';
import { getToken } from './auth';

// Types
export interface DownloadStats {
    total_downloads: number;
    downloads_by_edition: Record<string, number>;
    downloads_today: number;
    downloads_this_week: number;
    downloads_this_month: number;
    recent_downloads: Download[];
}

export interface Download {
    id: number;
    edition: string;
    ip_address: string;
    user_agent: string;
    user_id: number | null;
    created_at: string;
}

export interface DailyDownload {
    date: string;
    count: number;
}

// Helper function for API calls
async function apiCall<T>(
    endpoint: string,
    options: RequestInit = {}
): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const token = getToken();

    const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        ...(options.headers as Record<string, string>),
    };

    if (token && !endpoint.includes('/auth/')) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
        ...options,
        headers,
    });

    if (!response.ok) {
        const error = await response.json().catch(() => ({ error: 'Request failed' }));
        throw new Error(error.error || `HTTP ${response.status}`);
    }

    return response.json();
}

// Authentication API
export async function register(name: string, email: string, password: string) {
    return apiCall<{ id: number; name: string; email: string; role: string; avatar: string }>(
        '/auth/register',
        {
            method: 'POST',
            body: JSON.stringify({ name, email, password }),
        }
    );
}

export async function login(email: string, password: string) {
    return apiCall<{
        access_token: string;
        refresh_token: string;
        user: { id: number; name: string; email: string; role: string; avatar: string };
    }>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
    });
}

export async function logout() {
    return apiCall('/api/logout', { method: 'POST' });
}

export async function refreshToken(refreshToken: string) {
    return apiCall<{ access_token: string; refresh_token: string }>(
        '/auth/refresh',
        {
            method: 'POST',
            body: JSON.stringify({ refresh_token: refreshToken }),
        }
    );
}

export async function getMe() {
    return apiCall<{ id: number; name: string; email: string; role: string; avatar: string }>(
        '/api/me'
    );
}

// Download Tracking API
export async function trackDownload(edition: string): Promise<{ message: string; id: number }> {
    return apiCall<{ message: string; id: number }>('/api/downloads/track', {
        method: 'POST',
        body: JSON.stringify({ edition }),
    });
}

// Admin API
export async function getDownloadStats(): Promise<DownloadStats> {
    return apiCall<DownloadStats>('/api/admin/downloads/stats');
}

export async function getDownloadHistory(days: number = 30): Promise<DailyDownload[]> {
    return apiCall<DailyDownload[]>(`/api/admin/downloads/history?days=${days}`);
}
