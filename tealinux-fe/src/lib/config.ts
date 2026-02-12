// Use relative paths - will be proxied by Astro SSR or routed by Traefik
export const API_BASE_URL = import.meta.env.PUBLIC_API_URL || '';
