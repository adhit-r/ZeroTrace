/**
 * Application configuration
 * Centralized config to avoid hardcoded values
 */

export const config = {
  api: {
    baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
    timeout: 30000,
  },
  network: {
    defaultScanTargets: import.meta.env.VITE_DEFAULT_SCAN_TARGETS?.split(',') || ['192.168.1.0/24'],
    defaultScanType: 'tcp',
    defaultTimeout: 30,
    defaultConcurrency: 10,
  },
  organization: {
    // Will be fetched from auth context or API
    defaultId: null as string | null,
  },
  ui: {
    refreshInterval: 30000, // 30 seconds
    staleTime: 5 * 60 * 1000, // 5 minutes
  },
} as const;

export default config;

