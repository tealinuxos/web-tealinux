// Documentation utilities for reading and organizing markdown files

export interface DocMetadata {
    title: string;
    date?: string;
    navigation?: boolean;
}

export interface DocSection {
    title: string;
    slug: string;
    items: DocItem[];
}

export interface DocItem {
    title: string;
    slug: string;
    order: number;
}

/**
 * Convert file path to URL slug
 * Example: "1.Welcome To TealinuxOS/1.What is TeaLinux .md" -> "welcome-to-tealinuxos/what-is-tealinux"
 */
export function pathToSlug(path: string): string {
    return path
        .replace(/^\d+\./, '') // Remove leading numbers
        .replace(/\.md$/, '') // Remove .md extension
        .replace(/\s+/g, '-') // Replace spaces with hyphens
        .toLowerCase()
        .replace(/[^a-z0-9\/-]/g, ''); // Remove special characters
}

/**
 * Extract order number from filename
 * Example: "80.About Page.md" -> 80
 */
export function extractOrder(filename: string): number {
    const match = filename.match(/^(\d+)\./);
    return match ? parseInt(match[1], 10) : 999;
}

/**
 * Get navigation structure for documentation
 */
export function getDocNavigation(): DocSection[] {
    return [
        {
            title: 'Welcome to TeaLinuxOS',
            slug: 'welcome-to-tealinuxos',
            items: [
                { title: 'What is TeaLinux', slug: 'welcome-to-tealinuxos/what-is-tealinux', order: 1 },
                { title: "What's New", slug: 'welcome-to-tealinuxos/whats-new', order: 2 },
            ]
        },
        {
            title: 'Installation',
            slug: 'installation',
            items: [
                { title: 'Requirements', slug: 'installation/requirements', order: 1 },
                { title: 'Which Edition to Choose', slug: 'installation/which-edition-to-choose', order: 2 },
                { title: 'Boot TeaLinuxOS', slug: 'installation/boot-tealinuxos', order: 3 },
                { title: 'Boot TeaLinuxOS in VMs', slug: 'installation/boot-tealinuxos-in-vms', order: 4 },
                { title: 'Install Method', slug: 'installation/install-method', order: 5 },
            ]
        },
        {
            title: 'Installation Process',
            slug: 'installation-process',
            items: [
                { title: 'About Page', slug: 'about-page', order: 80 },
                { title: 'Setup Localization', slug: 'setup-localization', order: 81 },
                { title: 'Single Boot', slug: 'single-boot', order: 83 },
                { title: 'Dual Boot', slug: 'dual-boot', order: 84 },
                { title: 'Manual Partition', slug: 'manual-partition', order: 85 },
                { title: 'Create User', slug: 'create-user', order: 86 },
                { title: 'Summary', slug: 'summary', order: 87 },
                { title: 'Installation Stage', slug: 'installation-stage', order: 88 },
            ]
        },
    ];
}
