import { API_BASE_URL } from './config';
import { getToken } from './auth';

// Types
export interface Category {
    id: string;
    name: string;
    slug: string;
    description: string;
    order: number;
    created_at: string;
    updated_at: string;
}

export interface Topic {
    id: string;
    title: string;
    slug: string;
    user_id: number;
    user: {
        id: number;
        name: string;
        email: string;
        avatar: string;
    };
    category_id: string;
    category: Category;
    views: number;
    is_pinned: boolean;
    is_locked: boolean;
    created_at: string;
    updated_at: string;
    posts?: Post[];
    tags?: Tag[];
}

export interface Post {
    id: string;
    topic_id: string;
    user_id: number;
    user: {
        id: number;
        name: string;
        email: string;
        avatar: string;
    };
    content: string;
    reply_to_id?: string;
    created_at: string;
    updated_at: string;
    likes?: Like[];
}

export interface Like {
    id: string;
    post_id: string;
    user_id: number;
    created_at: string;
}

export interface Tag {
    id: string;
    name: string;
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

// Categories API
export async function getCategories(): Promise<Category[]> {
    return apiCall<Category[]>('/categories');
}

export async function getCategory(id: string): Promise<Category> {
    return apiCall<Category>(`/categories/${id}`);
}

// Topics API
export async function getTopics(categoryId?: string): Promise<Topic[]> {
    const query = categoryId ? `?category_id=${categoryId}` : '';
    return apiCall<Topic[]>(`/topics${query}`);
}

export async function getTopic(id: string): Promise<Topic> {
    return apiCall<Topic>(`/topics/${id}`);
}

export async function createTopic(data: {
    title: string;
    category_id: string;
    content: string;
    tags?: string[];
}): Promise<{ topic: Topic; post: Post }> {
    return apiCall<{ topic: Topic; post: Post }>('/api/topics', {
        method: 'POST',
        body: JSON.stringify(data),
    });
}

export async function updateTopic(id: string, title: string): Promise<Topic> {
    return apiCall<Topic>(`/api/topics/${id}`, {
        method: 'PUT',
        body: JSON.stringify({ title }),
    });
}

export async function deleteTopic(id: string): Promise<{ message: string }> {
    return apiCall<{ message: string }>(`/api/topics/${id}`, {
        method: 'DELETE',
    });
}

// Posts API
export async function getTopicPosts(topicId: string): Promise<Post[]> {
    return apiCall<Post[]>(`/topics/${topicId}/posts`);
}

export async function createPost(
    topicId: string,
    content: string,
    replyToId?: string
): Promise<Post> {
    return apiCall<Post>(`/api/topics/${topicId}/posts`, {
        method: 'POST',
        body: JSON.stringify({ content, reply_to_id: replyToId }),
    });
}

export async function updatePost(id: string, content: string): Promise<Post> {
    return apiCall<Post>(`/api/posts/${id}`, {
        method: 'PUT',
        body: JSON.stringify({ content }),
    });
}

export async function deletePost(id: string): Promise<{ message: string }> {
    return apiCall<{ message: string }>(`/api/posts/${id}`, {
        method: 'DELETE',
    });
}

// Likes API
export async function likePost(postId: string): Promise<{ message: string }> {
    return apiCall<{ message: string }>(`/api/posts/${postId}/like`, {
        method: 'POST',
    });
}

export async function unlikePost(postId: string): Promise<{ message: string }> {
    return apiCall<{ message: string }>(`/api/posts/${postId}/like`, {
        method: 'DELETE',
    });
}

// Search API
export async function search(query: string): Promise<{ topics: Topic[] }> {
    return apiCall<{ topics: Topic[] }>(`/search?q=${encodeURIComponent(query)}`);
}
