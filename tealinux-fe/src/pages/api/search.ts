// API endpoint for documentation search
import { getCollection } from 'astro:content';
import type { APIRoute } from 'astro';

export const GET: APIRoute = async ({ url }) => {
    const query = url.searchParams.get('q');

    if (!query || query.trim().length < 2) {
        return new Response(JSON.stringify([]), {
            status: 200,
            headers: { 'Content-Type': 'application/json' }
        });
    }

    const searchTerm = query.toLowerCase().trim();
    const results: any[] = [];

    try {
        const docs = await getCollection('docs');

        for (const doc of docs) {
            const { data, body, id } = doc;
            const title = data.title || '';
            const content = body || '';

            const titleMatch = title.toLowerCase().includes(searchTerm);
            const contentMatch = content.toLowerCase().includes(searchTerm);

            if (titleMatch || contentMatch) {
                const snippet = createSnippet(content, searchTerm);
                const highlightedSnippet = highlightMatches(snippet, searchTerm);
                const highlightedTitle = highlightMatches(title, searchTerm);

                results.push({
                    title: highlightedTitle,
                    url: `/docs/${id}`,
                    snippet: highlightedSnippet
                });
            }
        }

        // Sort by relevance
        results.sort((a, b) => {
            const aTitle = a.title.toLowerCase().includes(searchTerm);
            const bTitle = b.title.toLowerCase().includes(searchTerm);
            if (aTitle && !bTitle) return -1;
            if (!aTitle && bTitle) return 1;
            return 0;
        });

        return new Response(JSON.stringify(results.slice(0, 10)), {
            status: 200,
            headers: { 'Content-Type': 'application/json' }
        });
    } catch (error) {
        console.error('Search error:', error);
        return new Response(JSON.stringify([]), {
            status: 500,
            headers: { 'Content-Type': 'application/json' }
        });
    }
};

function createSnippet(content: string, searchTerm: string, maxLength: number = 150): string {
    const lowerContent = content.toLowerCase();
    const index = lowerContent.indexOf(searchTerm);

    if (index === -1) {
        return content.slice(0, maxLength).trim() + '...';
    }

    const start = Math.max(0, index - 60);
    const end = Math.min(content.length, index + searchTerm.length + 90);

    let snippet = content.slice(start, end).trim();

    if (start > 0) snippet = '...' + snippet;
    if (end < content.length) snippet = snippet + '...';

    snippet = snippet
        .replace(/#{1,6}\s/g, '')
        .replace(/\*\*/g, '')
        .replace(/\*/g, '')
        .replace(/\[([^\]]+)\]\([^\)]+\)/g, '$1')
        .replace(/`([^`]+)`/g, '$1')
        .replace(/\n+/g, ' ')
        .trim();

    return snippet;
}

function highlightMatches(text: string, searchTerm: string): string {
    if (!text || !searchTerm) return text;
    const regex = new RegExp(`(${escapeRegex(searchTerm)})`, 'gi');
    return text.replace(regex, '<span class="search-highlight">$1</span>');
}

function escapeRegex(str: string): string {
    return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}
