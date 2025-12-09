import { api } from './api';

export interface Agent {
  id: string;
  name: string;
  hostname: string;
  os: string;
  status: 'online' | 'offline' | 'unknown';
  last_seen: string;
  cpu_usage: number;
  memory_usage: number;
  organization_id: string;
  version: string;
  ip_address?: string;
  metadata?: Record<string, any>;
  mac_address?: string;
  os_version?: string;
  serial_number?: string;
  city?: string;
  region?: string;
  country?: string;
  timezone?: string;
  risk_score?: number;
  // Extended fields
  os_name?: string;
  os_build?: string;
  kernel_version?: string;
  cpu_model?: string;
  cpu_cores?: number;
  memory_total_gb?: number;
  storage_total_gb?: number;
  gpu_model?: string;
  platform?: string;
  tags?: string; // string based on split usage
}

export interface AgentStats {
  total: number;
  online: number;
  offline: number;
  avgCpu: number;
  avgMemory: number;
}

export interface AgentResponse {
  success: boolean;
  data: Agent[];
  message: string;
  timestamp: string;
}

export interface AgentStatsResponse {
  success: boolean;
  data: AgentStats;
  message: string;
  timestamp: string;
}

export const agentService = {
  async getAgents(): Promise<Agent[]> {
    try {
      console.log('Fetching agents from /api/agents/');
      const response = await api.get<AgentResponse>('/api/agents/');
      console.log('Agents response:', response.data);
      return response.data.data || [];
    } catch (error) {
      console.error('Failed to fetch agents:', error);
      return [];
    }
  },

  async getOnlineAgents(): Promise<Agent[]> {
    try {
      const response = await api.get<AgentResponse>('/api/agents/online');
      return response.data.data || [];
    } catch (error) {
      console.error('Failed to fetch online agents:', error);
      return [];
    }
  },

  async getAgentStats(): Promise<AgentStats> {
    try {
      console.log('Fetching agent stats from /api/agents/stats/public');
      const response = await api.get<AgentStatsResponse>('/api/agents/stats/public');
      console.log('Stats response:', response.data);
      return response.data.data || {
        total: 0,
        online: 0,
        offline: 0,
        avgCpu: 0,
        avgMemory: 0
      };
    } catch (error) {
      console.error('Failed to fetch agent stats:', error);
      return {
        total: 0,
        online: 0,
        offline: 0,
        avgCpu: 0,
        avgMemory: 0
      };
    }
  },

  async getAgent(id: string): Promise<Agent | null> {
    try {
      const response = await api.get<{ success: boolean; data: Agent }>(`/api/agents/${id}`);
      return response.data.data || null;
    } catch (error) {
      console.error('Failed to fetch agent:', error);
      return null;
    }
  },

  async restartAgent(id: string): Promise<boolean> {
    try {
      const response = await api.post(`/api/agents/${id}/restart`);
      return response.status === 200 || response.status === 204;
    } catch (error) {
      console.error('Failed to restart agent:', error);
      throw error;
    }
  },

  async killAgent(id: string): Promise<boolean> {
    try {
      const response = await api.delete(`/api/agents/${id}`);
      return response.status === 200 || response.status === 204;
    } catch (error) {
      console.error('Failed to kill agent:', error);
      throw error;
    }
  },

  async getDashboardOverview(): Promise<any> {
    try {
      const response = await api.get('/api/dashboard/overview');
      const data = response.data.data;

      // Transform API response to match Dashboard component expected format
      return {
        total_assets: data.assets?.total || 0,
        vulnerable_assets: data.assets?.vulnerable || 0,
        critical_vulnerabilities: data.vulnerabilities?.critical || 0,
        high_vulnerabilities: data.vulnerabilities?.high || 0,
        medium_vulnerabilities: data.vulnerabilities?.medium || 0,
        low_vulnerabilities: data.vulnerabilities?.low || 0,
        total_vulnerabilities: data.vulnerabilities?.total || 0,
        last_scan: data.assets?.lastScan || null,
        top_vulnerable_assets: data.top_vulnerable_assets || [],
        recent_scans: data.recent_scans || [],
        agents_online: data.agents?.online || 0,
        agents_total: data.agents?.total || 0,
        applications_total: data.applications?.total || 0
      };
    } catch (error) {
      console.error('Failed to fetch dashboard overview:', error);
      throw error;
    }
  },

  async getVulnerabilities(): Promise<any[]> {
    try {
      const response = await api.get('/api/vulnerabilities/');
      return response.data.data || [];
    } catch (error) {
      console.error('Failed to fetch vulnerabilities:', error);
      return [];
    }
  }
};
