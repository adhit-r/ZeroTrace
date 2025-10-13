import { agentService } from './agentService';

// Export interfaces first
export interface TopologyNode {
  id: string;
  name: string;
  type: 'agent' | 'server' | 'workstation' | 'router' | 'switch';
  ip: string;
  hasVulns: boolean;
  vulnerabilityCount: number;
  criticalVulns: number;
  status: 'online' | 'offline' | 'unknown';
  lastSeen: string;
  os: string;
  hostname: string;
  riskScore: number;
  location?: {
    city: string;
    region: string;
    country: string;
  };
}

export interface TopologyLink {
  source: string;
  target: string;
  type: 'network' | 'dependency' | 'communication';
  bandwidth?: number;
  latency?: number;
}

export interface TopologyData {
  nodes: TopologyNode[];
  links: TopologyLink[];
  clusters: {
    id: string;
    name: string;
    type: string;
    nodes: string[];
  }[];
  metadata: {
    totalNodes: number;
    totalLinks: number;
    vulnerableNodes: number;
    criticalNodes: number;
    lastUpdated: string;
  };
}

class TopologyService {
  private baseUrl = '/api';

  async getTopologyData(): Promise<TopologyData> {
    try {
      // Fetch agents data
      const agents = await agentService.getAgents();
      const agentStats = await agentService.getAgentStats();

      // Convert agents to topology nodes
      const nodes: TopologyNode[] = agents.map(agent => {
        const hasVulns = (agent.metadata?.total_vulnerabilities || 0) > 0;
        const vulnerabilityCount = agent.metadata?.total_vulnerabilities || 0;
        const criticalVulns = agent.metadata?.critical_vulnerabilities || 0;
        
        // Determine node type based on OS and role
        let nodeType: TopologyNode['type'] = 'workstation';
        if (agent.os?.toLowerCase().includes('server') || agent.hostname?.toLowerCase().includes('server')) {
          nodeType = 'server';
        } else if (agent.hostname?.toLowerCase().includes('router') || agent.hostname?.toLowerCase().includes('gateway')) {
          nodeType = 'router';
        } else if (agent.hostname?.toLowerCase().includes('switch')) {
          nodeType = 'switch';
        }

        // Calculate risk score
        const riskScore = this.calculateRiskScore(vulnerabilityCount, criticalVulns, agent.status);

        return {
          id: agent.id,
          name: agent.hostname || agent.name || `Agent-${agent.id.slice(0, 8)}`,
          type: nodeType,
          ip: agent.ip_address || 'Unknown',
          hasVulns,
          vulnerabilityCount,
          criticalVulns,
          status: agent.status as 'online' | 'offline' | 'unknown',
          lastSeen: agent.last_seen,
          os: agent.os || 'Unknown',
          hostname: agent.hostname || 'Unknown',
          riskScore,
          location: agent.metadata?.location ? {
            city: agent.metadata.location.city || 'Unknown',
            region: agent.metadata.location.region || 'Unknown',
            country: agent.metadata.location.country || 'Unknown'
          } : undefined
        };
      });

      // Generate network links based on agent relationships
      const links = this.generateNetworkLinks(nodes);

      // Create clusters based on location and type
      const clusters = this.createClusters(nodes);

      // Calculate metadata
      const vulnerableNodes = nodes.filter(node => node.hasVulns).length;
      const criticalNodes = nodes.filter(node => node.criticalVulns > 0).length;

      return {
        nodes,
        links,
        clusters,
        metadata: {
          totalNodes: nodes.length,
          totalLinks: links.length,
          vulnerableNodes,
          criticalNodes,
          lastUpdated: new Date().toISOString()
        }
      };

    } catch (error) {
      console.error('Error fetching topology data:', error);
      throw new Error('Failed to fetch topology data');
    }
  }

  private calculateRiskScore(vulnerabilityCount: number, criticalVulns: number, status: string): number {
    let score = 0;
    
    // Base score from vulnerability count
    score += Math.min(vulnerabilityCount * 0.5, 5);
    
    // Critical vulnerabilities add more risk
    score += Math.min(criticalVulns * 2, 10);
    
    // Offline agents are higher risk
    if (status === 'offline') {
      score += 3;
    }
    
    // Cap at 10
    return Math.min(score, 10);
  }

  private generateNetworkLinks(nodes: TopologyNode[]): TopologyLink[] {
    const links: TopologyLink[] = [];
    
    // Create a simple network topology based on IP addresses and types
    const routers = nodes.filter(node => node.type === 'router');
    const switches = nodes.filter(node => node.type === 'switch');
    const servers = nodes.filter(node => node.type === 'server');
    const workstations = nodes.filter(node => node.type === 'workstation');

    // Connect routers to switches
    routers.forEach(router => {
      switches.forEach(switchNode => {
        links.push({
          source: router.id,
          target: switchNode.id,
          type: 'network',
          bandwidth: 1000, // 1Gbps
          latency: 1
        });
      });
    });

    // Connect switches to servers and workstations
    switches.forEach(switchNode => {
      [...servers, ...workstations].forEach(node => {
        links.push({
          source: switchNode.id,
          target: node.id,
          type: 'network',
          bandwidth: node.type === 'server' ? 1000 : 100, // Servers get 1Gbps, workstations 100Mbps
          latency: 2
        });
      });
    });

    // Add some inter-server communication links
    servers.forEach(server => {
      servers.forEach(otherServer => {
        if (server.id !== otherServer.id) {
          links.push({
            source: server.id,
            target: otherServer.id,
            type: 'communication',
            bandwidth: 1000,
            latency: 1
          });
        }
      });
    });

    return links;
  }

  private createClusters(nodes: TopologyNode[]): TopologyData['clusters'] {
    const clusters: TopologyData['clusters'] = [];
    
    // Group by location
    const locationGroups = nodes.reduce((groups, node) => {
      const location = node.location?.city || 'Unknown';
      if (!groups[location]) {
        groups[location] = [];
      }
      groups[location].push(node.id);
      return groups;
    }, {} as Record<string, string[]>);

    Object.entries(locationGroups).forEach(([location, nodeIds]) => {
      if (nodeIds.length > 1) {
        clusters.push({
          id: `location-${location.toLowerCase().replace(/\s+/g, '-')}`,
          name: `${location} Office`,
          type: 'location',
          nodes: nodeIds
        });
      }
    });

    // Group by type
    const typeGroups = nodes.reduce((groups, node) => {
      if (!groups[node.type]) {
        groups[node.type] = [];
      }
      groups[node.type].push(node.id);
      return groups;
    }, {} as Record<string, string[]>);

    Object.entries(typeGroups).forEach(([type, nodeIds]) => {
      if (nodeIds.length > 1) {
        clusters.push({
          id: `type-${type}`,
          name: `${type.charAt(0).toUpperCase() + type.slice(1)}s`,
          type: 'device_type',
          nodes: nodeIds
        });
      }
    });

    return clusters;
  }

  async getNodeDetails(nodeId: string): Promise<TopologyNode | null> {
    try {
      const topologyData = await this.getTopologyData();
      return topologyData.nodes.find(node => node.id === nodeId) || null;
    } catch (error) {
      console.error('Error fetching node details:', error);
      return null;
    }
  }

  async getNetworkHealth(): Promise<{
    totalNodes: number;
    onlineNodes: number;
    vulnerableNodes: number;
    criticalNodes: number;
    averageRiskScore: number;
  }> {
    try {
      const topologyData = await this.getTopologyData();
      
      const onlineNodes = topologyData.nodes.filter(node => node.status === 'online').length;
      const vulnerableNodes = topologyData.nodes.filter(node => node.hasVulns).length;
      const criticalNodes = topologyData.nodes.filter(node => node.criticalVulns > 0).length;
      const averageRiskScore = topologyData.nodes.reduce((sum, node) => sum + node.riskScore, 0) / topologyData.nodes.length;

      return {
        totalNodes: topologyData.nodes.length,
        onlineNodes,
        vulnerableNodes,
        criticalNodes,
        averageRiskScore: Math.round(averageRiskScore * 10) / 10
      };
    } catch (error) {
      console.error('Error fetching network health:', error);
      throw new Error('Failed to fetch network health');
    }
  }
}

// Export the service instance
export const topologyService = new TopologyService();