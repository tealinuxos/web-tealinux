import type { APIRoute } from 'astro';
import { getCollection } from 'astro:content';

export const prerender = true;

export const GET: APIRoute = async () => {
    const docs = await getCollection('docs');

    const searchIndex = docs.map(doc => ({
        title: doc.data.title,
        id: doc.id,
        body: doc.body || '',
        category: doc.data.category || '',
    }));

    return new Response(JSON.stringify(searchIndex), {
        headers: {
            'Content-Type': 'application/json'
        }
    });
};
