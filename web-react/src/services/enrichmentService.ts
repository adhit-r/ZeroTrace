
import axios from 'axios';

// Direct connection to Enrichment Service (running on port 8000 by default)
// In production, this should be proxied via the main API gateway
const ENRICHMENT_URL = 'http://localhost:8000';

export interface SoftwarePackage {
    name: string;
    version: string;
    vendor?: string;
    cpe?: string | null;
    vulnerabilities?: Vulnerability[];
}

export interface Vulnerability {
    id: string;
    cve_id: string;
    description: string;
    severity: 'critical' | 'high' | 'medium' | 'low' | 'unknown';
    cvss_score: number;
    source: string;
    known_exploited?: boolean;
}

export interface ScanResult {
    vulnerabilities: Array<Vulnerability & { software: string; agent: string }>;
    scannedCount: number;
}

export const enrichmentService = {
    /**
     * Enrich a batch of software packages to detect vulnerabilities
     */
    async enrichSoftware(software: SoftwarePackage[]): Promise<SoftwarePackage[]> {
        try {
            const response = await axios.post<{ data: SoftwarePackage[] }>(
                `${ENRICHMENT_URL}/enrich/software`,
                software,
                {
                    headers: {
                        'Content-Type': 'application/json'
                    }
                }
            );
            return response.data.data || [];
        } catch (error) {
            console.error('Enrichment failed:', error);
            throw error;
        }
    },

    /**
     * Filter software based on priority categories and scan scope
     */
    filterSoftware(software: any[], scope: 'quick' | 'standard' | 'priority' | 'full', categories: string[]): any[] {
        const PRIORITY_PATTERNS: Record<string, string[]> = {
            database: ['mysql', 'postgres', 'mongodb', 'redis', 'sqlite', 'mariadb', 'oracle', 'sqlserver', 'cassandra', 'elasticsearch'],
            webserver: ['nginx', 'apache', 'httpd', 'tomcat', 'iis', 'lighttpd', 'caddy', 'haproxy', 'envoy'],
            runtime: ['node', 'python', 'java', 'ruby', 'php', 'dotnet', 'go', 'perl', 'rust', 'swift'],
            crypto: ['openssl', 'libressl', 'gnutls', 'boringssl', 'cryptlib', 'libsodium'],
            framework: ['django', 'flask', 'express', 'rails', 'spring', 'laravel', 'react', 'vue', 'angular', 'jquery'],
            browser: ['chrome', 'firefox', 'safari', 'edge', 'chromium', 'electron', 'webkit']
        };

        // Combine keywords from selected categories
        let priorityKeywords: string[] = [];
        categories.forEach(cat => {
            if (PRIORITY_PATTERNS[cat]) {
                priorityKeywords = priorityKeywords.concat(PRIORITY_PATTERNS[cat]);
            }
        });

        // Deduplicate software
        const uniqueSoftware = new Map<string, any>();
        software.forEach(sw => {
            const key = `${sw.name}:${sw.version}`;
            if (!uniqueSoftware.has(key)) {
                uniqueSoftware.set(key, sw);
            }
        });

        // Categorize
        const priorityApps: any[] = [];
        const otherApps: any[] = [];

        uniqueSoftware.forEach(sw => {
            const name = (sw.name || '').toLowerCase();
            const isPriority = priorityKeywords.some(kw => name.includes(kw));
            if (isPriority) {
                priorityApps.push(sw);
            } else {
                otherApps.push(sw);
            }
        });

        // Apply Scope
        switch (scope) {
            case 'quick':
                return priorityApps.slice(0, 50);
            case 'standard':
                return [...priorityApps, ...otherApps.slice(0, 200 - priorityApps.length)].slice(0, 200);
            case 'priority':
                return priorityApps;
            case 'full':
                return [...priorityApps, ...otherApps].slice(0, 1000);
            default:
                return priorityApps.slice(0, 100);
        }
    }
};
