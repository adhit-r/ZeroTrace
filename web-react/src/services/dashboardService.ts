import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export interface Asset {
  total: number;
  vulnerable: number;
  critical: number;
  high: number;
  medium: number;
  low: number;
  lastScan: string | null;
}

export interface Vulnerability {
  id: number;
  name: string;
  severity: string;
  asset: string;
  status: string;
}

export interface RecentActivity {
  id: number;
  type: string;
  message: string;
  time: string;
}

export interface TopVulnerableAsset {
  name: string;
  vulnerabilities: number;
  critical: number;
}

export interface DashboardData {
  assets: Asset;
  vulnerabilities: Vulnerability[];
  recentActivity: RecentActivity[];
  topVulnerableAssets: TopVulnerableAsset[];
  scanStatus: string;
}

// Re-export for compatibility
export type { Asset, Vulnerability, RecentActivity, TopVulnerableAsset };

export const dashboardService = {
  async getDashboardData(): Promise<DashboardData> {
    try {
      // Fetch dashboard overview from API
      const response = await api.get('/api/dashboard/overview');
      const apiData = response.data.data;
      
      // Transform API data to match frontend structure
      return {
        assets: apiData.assets || this.getDefaultData().assets,
        vulnerabilities: [], // API returns vulnerability counts, not array
        recentActivity: [],
        topVulnerableAssets: [],
        scanStatus: apiData.agents?.online > 0 ? 'active' : 'idle'
      };
    } catch (error) {
      console.error('Failed to fetch dashboard data:', error);
      return this.getDefaultData();
    }
  },

  getDefaultData(): DashboardData {
    return {
      assets: {
        total: 0,
        vulnerable: 0,
        critical: 0,
        high: 0,
        medium: 0,
        low: 0,
        lastScan: null
      },
      vulnerabilities: [],
      recentActivity: [],
      topVulnerableAssets: [],
      scanStatus: 'idle'
    };
  }
};
