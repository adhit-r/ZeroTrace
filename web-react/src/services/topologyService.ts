import { api } from './api';
import type { Agent } from '../types/api';

const fetchAgentData = async (): Promise<Agent[]> => {
  try {
    const response = await api.get('/api/agents/');
    return response.data?.data || [];
  } catch (error) {
    console.error("Failed to fetch agent data:", error);
    return [];
  }
};

export const topologyService = {
  async getTopologyData() {
    try {
      const agents = await fetchAgentData();
      if (!agents) {
        return { nodes: [], links: [] };
      }
      
      const nodes = agents.map(agent => ({
        id: agent.id,
        name: agent.hostname,
        type: 'agent',
        status: agent.status,
        riskScore: (agent as any).risk_score || 0,
        ipAddress: (agent as any).ip_address,
        location: `${(agent as any).city}, ${(agent as any).country}`,
        clusterId: (agent as any).branch_id || 'default',
        metadata: agent.metadata || {}
      }));
      
      const links = agents.flatMap(agent => 
        (agent.metadata?.connections || []).map((conn: any) => ({
          source: agent.id,
          target: conn.target_agent_id,
          type: 'network',
          strength: conn.strength || 1
        }))
      );
      
      return { nodes, links };
    } catch (error) {
      console.error('Error fetching topology data:', error);
      return { nodes: [], links: [] };
    }
  }
};