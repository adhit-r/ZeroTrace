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
  }
};
