// Search utility for documentation content
// Uses API endpoint to avoid client-side astro:content issues

export interface SearchResult {
    title: string;
    url: string;
    snippet: string;
}

/**
 * Search through documentation content via API
 * @param query - Search query string
 * @returns Array of search results with highlighted snippets
 */
export async function searchDocs(query: string): Promise<SearchResult[]> {
    if (!query || query.trim().length < 2) {
        return [];
    }

    try {
        const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`);

        if (!response.ok) {
            throw new Error('Search request failed');
        }

        const results: SearchResult[] = await response.json();
        return results;
    } catch (error) {
        console.error('Search error:', error);
        return [];
    }
}
