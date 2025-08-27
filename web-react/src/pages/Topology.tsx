import React, { useState, useEffect } from 'react';
import NetworkTopology from '../components/NetworkTopology';

// Mock data for demonstration
const mockTopologyData = {
  id: "topology-1",
  companyId: "company-123",
  nodes: [
    // Vulnerability Management Agents
    {
      id: "vuln-agent-1",
      type: "agent" as const,
      agentId: "agent-001",
      name: "Vuln Scanner DC1",
      ipAddress: "192.168.1.10",
      location: "HQ - Floor 1",
      riskScore: 15,
      status: "active" as const,
      metadata: { vulnerabilities: 15, lastSeen: new Date() }
    },
    {
      id: "vuln-agent-2",
      type: "agent" as const,
      agentId: "agent-002",
      name: "Vuln Scanner DC2",
      ipAddress: "192.168.2.10",
      location: "HQ - Floor 2",
      riskScore: 8,
      status: "active" as const,
      metadata: { vulnerabilities: 8, lastSeen: new Date() }
    },
    
    // Endpoint Agents
    {
      id: "endpoint-1",
      type: "asset" as const,
      assetId: "asset-001",
      name: "Workstation-001",
      ipAddress: "192.168.1.101",
      location: "HQ - Floor 1",
      riskScore: 3,
      status: "active" as const,
      metadata: { os: "Windows 11", vulnerabilities: 3 }
    },
    {
      id: "endpoint-2",
      type: "asset" as const,
      assetId: "asset-002",
      name: "Server-DB",
      ipAddress: "192.168.1.50",
      location: "HQ - Basement",
      riskScore: 12,
      status: "active" as const,
      metadata: { os: "Ubuntu 22.04", vulnerabilities: 12 }
    },
    
    // Network Devices
    {
      id: "switch-1",
      type: "asset" as const,
      assetId: "asset-003",
      name: "Core Switch",
      ipAddress: "192.168.1.1",
      location: "HQ - Basement",
      riskScore: 0,
      status: "active" as const,
      metadata: { deviceType: "network", vulnerabilities: 0 }
    },
    {
      id: "router-1",
      type: "asset" as const,
      assetId: "asset-004",
      name: "Border Router",
      ipAddress: "10.0.0.1",
      location: "HQ - Basement",
      riskScore: 2,
      status: "active" as const,
      metadata: { deviceType: "network", vulnerabilities: 2 }
    },
    
    // OWASP Amass Discovered Assets
    {
      id: "amass-1",
      type: "amass_discovery" as const,
      name: "External API",
      ipAddress: "203.0.113.10",
      location: "Cloud - External",
      riskScore: 5,
      status: "discovered" as const,
      metadata: { domain: "api.company.com", source: "amass" }
    },
    {
      id: "amass-2",
      type: "amass_discovery" as const,
      name: "Mail Server",
      ipAddress: "203.0.113.20",
      location: "Cloud - External",
      riskScore: 1,
      status: "discovered" as const,
      metadata: { domain: "mail.company.com", source: "amass" }
    },
  ],
  links: [
    { source: "vuln-agent-1", target: "endpoint-1", type: "scan" as const, strength: 1 },
    { source: "vuln-agent-1", target: "switch-1", type: "scan" as const, strength: 1 },
    { source: "vuln-agent-2", target: "endpoint-2", type: "scan" as const, strength: 1 },
    { source: "switch-1", target: "endpoint-1", type: "network" as const, strength: 2 },
    { source: "switch-1", target: "endpoint-2", type: "network" as const, strength: 2 },
    { source: "switch-1", target: "router-1", type: "network" as const, strength: 3 },
    { source: "router-1", target: "amass-1", type: "external" as const, strength: 1 },
    { source: "router-1", target: "amass-2", type: "external" as const, strength: 1 },
  ],
  clusters: [
    {
      id: "cluster-1",
      name: "HQ Network",
      type: "subnet",
      nodeIds: ["vuln-agent-1", "endpoint-1", "switch-1"],
      riskScore: 8.5,
      description: "Main headquarters network"
    },
    {
      id: "cluster-2",
      name: "External Assets",
      type: "geographic",
      nodeIds: ["amass-1", "amass-2"],
      riskScore: 3.0,
      description: "Externally discovered assets"
    }
  ],
  lastUpdated: new Date().toISOString()
};

const Topology: React.FC = () => {
  const [topologyData, setTopologyData] = useState(mockTopologyData);
  const [viewMode, setViewMode] = useState<'network' | 'floor' | 'geographic' | 'cluster'>('network');
  const [agentFilter, setAgentFilter] = useState<'all' | 'vulnerability' | 'endpoint' | 'network' | 'amass'>('all');
  const [connectionFilter, setConnectionFilter] = useState<'all' | 'connected' | 'disconnected' | 'intermittent'>('all');
  const [isLoading, setIsLoading] = useState(false);

  const handleRefresh = async () => {
    setIsLoading(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // Update mock data with some randomness
      const updatedData = {
        ...topologyData,
        nodes: topologyData.nodes.map(node => ({
          ...node,
          riskScore: node.riskScore + (Math.random() > 0.7 ? Math.floor(Math.random() * 3) : 0),
          metadata: {
            ...node.metadata,
            lastSeen: new Date()
          }
        })),
        lastUpdated: new Date().toISOString()
      };
      
      setTopologyData(updatedData);
    } catch (error) {
      console.error('Error refreshing topology data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleNodeClick = (node: any) => {
    console.log('Node clicked:', node);
    // Here you could open a detailed view, trigger a rescan, etc.
    alert(`Clicked on ${node.name} (${node.type})`);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-blue-900 to-gray-900">
      {/* Header */}
      <div className="bg-black bg-opacity-80 backdrop-blur-sm border-b border-white border-opacity-10 p-4">
        <h1 className="text-2xl font-bold text-cyan-400 mb-1">
          Network Asset Discovery & Topology Visualizer
        </h1>
        <p className="text-gray-300 text-sm">
          Real-time visualization of vulnerability management agents and discovered assets
        </p>
      </div>

      {/* Main Content */}
      <div className="relative h-screen">
        {isLoading && (
          <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white bg-opacity-10 backdrop-blur-sm p-6 rounded-lg border border-white border-opacity-20">
              <div className="flex items-center space-x-3">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-cyan-400"></div>
                <span className="text-white">Refreshing topology data...</span>
              </div>
            </div>
          </div>
        )}

        <NetworkTopology
          data={topologyData}
          onNodeClick={handleNodeClick}
          onRefresh={handleRefresh}
          viewMode={viewMode}
          agentFilter={agentFilter}
          connectionFilter={connectionFilter}
        />
      </div>
    </div>
  );
};

export default Topology;
