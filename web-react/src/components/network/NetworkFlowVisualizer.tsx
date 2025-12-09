import React, { useCallback, useEffect, useState, useMemo } from 'react';
import ReactFlow, {
  type Node,
  type Edge,
  Controls,
  Background,
  MiniMap,
  addEdge,
  useNodesState,
  useEdgesState,
  Panel,
  MarkerType,
  type NodeTypes,
} from 'reactflow';

// Connection type definition (not exported in reactflow v11)
type Connection = {
  source: string | null;
  target: string | null;
  sourceHandle: string | null;
  targetHandle: string | null;
};
import 'reactflow/dist/style.css';
import {
  Server,
  Router,
  Wifi,
  Smartphone,
  Shield,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Network,
  Database,
  Search,
  Download,
  RefreshCw,
} from 'lucide-react';
import { networkScanService } from '../../services/networkScanService';

// Custom node types for different device types
const DeviceNode = ({ data, selected }: { data: any; selected: boolean }) => {
  const getIcon = () => {
    switch (data.deviceType) {
      case 'switch':
        return <Network className="w-6 h-6" />;
      case 'router':
        return <Router className="w-6 h-6" />;
      case 'iot':
        return <Wifi className="w-6 h-6" />;
      case 'phone':
        return <Smartphone className="w-6 h-6" />;
      case 'server':
        return <Server className="w-6 h-6" />;
      default:
        return <Network className="w-6 h-6" />;
    }
  };

  const getStatusColor = () => {
    if (data.riskScore >= 70) return 'bg-red-500';
    if (data.riskScore >= 40) return 'bg-yellow-500';
    return 'bg-green-500';
  };

  return (
    <div
      className={`px-4 py-3 rounded-lg border-2 shadow-lg transition-all ${selected ? 'border-blue-500 shadow-xl' : 'border-gray-300'
        } ${getStatusColor()} bg-white min-w-[200px]`}
    >
      <div className="flex items-center gap-2 mb-2">
        <div className="text-gray-700">{getIcon()}</div>
        <div className="flex-1">
          <div className="font-semibold text-sm text-gray-900">{data.label}</div>
          <div className="text-xs text-gray-600">{data.deviceType}</div>
        </div>
        {data.status === 'online' ? (
          <CheckCircle className="w-4 h-4 text-green-500" />
        ) : (
          <XCircle className="w-4 h-4 text-red-500" />
        )}
      </div>

      <div className="mt-2 space-y-1 text-xs">
        <div className="flex justify-between">
          <span className="text-gray-600">IP:</span>
          <span className="font-mono text-gray-900">{data.ipAddress}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-600">Risk:</span>
          <span className={`font-semibold ${data.riskScore >= 70 ? 'text-red-600' :
            data.riskScore >= 40 ? 'text-yellow-600' : 'text-green-600'
            }`}>
            {data.riskScore.toFixed(1)}
          </span>
        </div>
        {data.vulnerabilities > 0 && (
          <div className="flex items-center gap-1 text-red-600">
            <AlertTriangle className="w-3 h-3" />
            <span>{data.vulnerabilities} vulns</span>
          </div>
        )}
        {data.openPorts && (
          <div className="text-gray-600">
            {data.openPorts.length} open ports
          </div>
        )}
      </div>
    </div>
  );
};

const VulnerabilityNode = ({ data, selected }: { data: any; selected: boolean }) => {
  const getSeverityColor = () => {
    switch (data.severity) {
      case 'critical':
        return 'bg-red-600 border-red-700';
      case 'high':
        return 'bg-orange-500 border-orange-600';
      case 'medium':
        return 'bg-yellow-500 border-yellow-600';
      case 'low':
        return 'bg-blue-500 border-blue-600';
      default:
        return 'bg-gray-500 border-gray-600';
    }
  };

  return (
    <div
      className={`px-3 py-2 rounded-lg border-2 shadow-md transition-all ${selected ? 'border-blue-500 shadow-xl' : ''
        } ${getSeverityColor()} text-white min-w-[180px]`}
    >
      <div className="flex items-center gap-2 mb-1">
        <Shield className="w-4 h-4" />
        <div className="font-semibold text-sm">{data.label}</div>
      </div>
      <div className="text-xs opacity-90">
        <div>{data.description}</div>
        {data.cve && <div className="mt-1 font-mono">CVE: {data.cve}</div>}
      </div>
    </div>
  );
};

const ServiceNode = ({ data, selected }: { data: any; selected: boolean }) => {
  return (
    <div
      className={`px-3 py-2 rounded-lg border-2 shadow-md transition-all ${selected ? 'border-blue-500 shadow-xl' : 'border-purple-300'
        } bg-purple-50 min-w-[150px]`}
    >
      <div className="flex items-center gap-2">
        <Database className="w-4 h-4 text-purple-600" />
        <div>
          <div className="font-semibold text-sm text-gray-900">{data.label}</div>
          <div className="text-xs text-gray-600">Port: {data.port}</div>
        </div>
      </div>
    </div>
  );
};

// Define custom node types
const nodeTypes: NodeTypes = {
  device: DeviceNode,
  vulnerability: VulnerabilityNode,
  service: ServiceNode,
};

interface NetworkFlowVisualizerProps {
  className?: string;
  onNodeClick?: (node: Node) => void;
  onEdgeClick?: (edge: Edge) => void;
}



const NetworkFlowVisualizer: React.FC<NetworkFlowVisualizerProps> = ({
  className = '',
  onNodeClick,
  onEdgeClick,
}) => {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterType, setFilterType] = useState<string>('all');
  const [selectedNode, setSelectedNode] = useState<Node | null>(null);

  // Fetch network scan results
  useEffect(() => {
    const fetchNetworkData = async () => {
      try {
        setIsLoading(true);
        setError(null);

        // Get network topology data
        const topologyData = await networkScanService.getNetworkTopology();

        // Convert to ReactFlow nodes and edges
        const deviceNodes: Node[] = topologyData.nodes
          .filter((n) => n.type === 'device')
          .map((node, index) => ({
            id: node.id,
            type: 'device',
            position: {
              x: (index % 5) * 250 + 100,
              y: Math.floor(index / 5) * 200 + 100,
            },
            data: {
              label: node.label,
              deviceType: node.deviceType,
              ipAddress: node.ipAddress,
              os: node.os,
              riskScore: node.riskScore || 0,
              vulnerabilities: node.vulnerabilities || 0,
              openPorts: node.openPorts || [],
              status: node.status || 'online',
            },
          }));

        const vulnNodes: Node[] = topologyData.nodes
          .filter((n) => n.type === 'vulnerability')
          .map((node, index) => ({
            id: node.id,
            type: 'vulnerability',
            position: {
              x: (index % 4) * 220 + 50,
              y: Math.floor(deviceNodes.length / 5) * 200 + 300 + Math.floor(index / 4) * 150,
            },
            data: {
              label: node.label,
              severity: node.severity,
              description: node.description,
              cve: node.cve,
              host: node.ipAddress,
              port: node.port,
            },
          }));

        const serviceNodes: Node[] = topologyData.nodes
          .filter((n) => n.type === 'service')
          .map((node, index) => ({
            id: node.id,
            type: 'service',
            position: {
              x: (index % 6) * 180 + 100,
              y: Math.floor(deviceNodes.length / 5) * 200 + 500 + Math.floor(index / 6) * 120,
            },
            data: {
              label: node.label,
              port: node.port,
              protocol: node.protocol,
              host: node.ipAddress,
            },
          }));

        // Convert edges
        const newEdges: Edge[] = topologyData.edges.map((edge) => {
          // const sourceNode = topologyData.nodes.find((n) => n.id === edge.source);
          const targetNode = topologyData.nodes.find((n) => n.id === edge.target);

          let strokeColor = '#8b5cf6';
          if (targetNode?.type === 'vulnerability') {
            const severity = targetNode.severity;
            strokeColor =
              severity === 'critical'
                ? '#ef4444'
                : severity === 'high'
                  ? '#f97316'
                  : '#eab308';
          }

          return {
            id: edge.source + '-' + edge.target,
            source: edge.source,
            target: edge.target,
            type: 'smoothstep',
            animated: true,
            style: { stroke: strokeColor, strokeWidth: 2 },
            markerEnd: {
              type: MarkerType.ArrowClosed,
              color: strokeColor,
            },
          };
        });

        const allNodes = [...deviceNodes, ...vulnNodes, ...serviceNodes];
        setNodes(allNodes);
        setEdges(newEdges);
      } catch (err: any) {
        console.error('Failed to fetch network data:', err);
        setError(err.message || 'Failed to load network data');
      } finally {
        setIsLoading(false);
      }
    };

    fetchNetworkData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Filter nodes based on search and filter
  const filteredNodes = useMemo(() => {
    return nodes.filter((node) => {
      const matchesSearch =
        searchTerm === '' ||
        node.data.label?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        node.data.ipAddress?.toLowerCase().includes(searchTerm.toLowerCase());

      const matchesFilter =
        filterType === 'all' ||
        (filterType === 'devices' && node.type === 'device') ||
        (filterType === 'vulnerabilities' && node.type === 'vulnerability') ||
        (filterType === 'services' && node.type === 'service');

      return matchesSearch && matchesFilter;
    });
  }, [nodes, searchTerm, filterType]);

  // Update edges when nodes are filtered
  const filteredEdges = useMemo(() => {
    const filteredNodeIds = new Set(filteredNodes.map((n) => n.id));
    return edges.filter(
      (edge) => filteredNodeIds.has(edge.source) && filteredNodeIds.has(edge.target)
    );
  }, [edges, filteredNodes]);

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  );

  const handleNodeClick = useCallback(
    (_event: React.MouseEvent, node: Node) => {
      setSelectedNode(node);
      if (onNodeClick) {
        onNodeClick(node);
      }
    },
    [onNodeClick]
  );

  const handleEdgeClick = useCallback(
    (_event: React.MouseEvent, edge: Edge) => {
      if (onEdgeClick) {
        onEdgeClick(edge);
      }
    },
    [onEdgeClick]
  );

  const exportData = useCallback(() => {
    const data = {
      nodes: filteredNodes,
      edges: filteredEdges,
      exportedAt: new Date().toISOString(),
    };
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `network-topology-${Date.now()}.json`;
    a.click();
    URL.revokeObjectURL(url);
  }, [filteredNodes, filteredEdges]);

  if (isLoading) {
    return (
      <div className={`flex items-center justify-center h-full ${className}`}>
        <div className="text-center">
          <RefreshCw className="w-8 h-8 animate-spin mx-auto mb-4 text-blue-500" />
          <p className="text-gray-600">Loading network topology...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`flex items-center justify-center h-full ${className}`}>
        <div className="text-center">
          <AlertTriangle className="w-8 h-8 mx-auto mb-4 text-red-500" />
          <p className="text-red-600">{error}</p>
        </div>
      </div>
    );
  }

  // Show empty state if no nodes
  if (nodes.length === 0 && !isLoading) {
    return (
      <div className={`flex items-center justify-center h-full ${className}`}>
        <div className="text-center max-w-md">
          <Network className="w-16 h-16 mx-auto mb-4 text-gray-400" />
          <h3 className="text-xl font-semibold text-gray-900 mb-2">No Network Scan Results</h3>
          <p className="text-gray-600 mb-4">
            Network scan results will appear here once the agent completes a scan.
          </p>
          <div className="text-sm text-gray-500 space-y-2">
            <p>• Network scans run every 6 hours by default</p>
            <p>• First scan starts 30 seconds after agent launch</p>
            <p>• Check agent logs for scan progress</p>
          </div>
          <button
            onClick={() => window.location.reload()}
            className="mt-6 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            <RefreshCw className="w-4 h-4 inline mr-2" />
            Refresh
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className={`w-full h-full relative ${className}`}>
      <ReactFlow
        nodes={filteredNodes}
        edges={filteredEdges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onNodeClick={handleNodeClick}
        onEdgeClick={handleEdgeClick}
        nodeTypes={nodeTypes}
        fitView
        attributionPosition="bottom-left"
      >
        <Controls />
        <Background />
        <MiniMap
          nodeColor={(node) => {
            if (node.type === 'device') {
              const risk = node.data?.riskScore || 0;
              return risk >= 70 ? '#ef4444' : risk >= 40 ? '#f97316' : '#10b981';
            }
            if (node.type === 'vulnerability') {
              const severity = node.data?.severity || 'low';
              return severity === 'critical'
                ? '#ef4444'
                : severity === 'high'
                  ? '#f97316'
                  : severity === 'medium'
                    ? '#eab308'
                    : '#3b82f6';
            }
            return '#8b5cf6';
          }}
        />

        <Panel position="top-left" className="bg-white p-4 rounded-lg shadow-lg m-4">
          <div className="space-y-3">
            <div className="flex items-center gap-2 mb-3">
              <Network className="w-5 h-5 text-blue-500" />
              <h3 className="font-semibold text-lg">Network Topology</h3>
            </div>

            <div className="relative">
              <Search className="absolute left-2 top-2.5 w-4 h-4 text-gray-400" />
              <input
                type="text"
                placeholder="Search devices..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-8 pr-3 py-2 border border-gray-300 rounded-md text-sm w-full focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-xs font-medium text-gray-700 mb-1">Filter by Type</label>
              <select
                value={filterType}
                onChange={(e) => setFilterType(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="all">All</option>
                <option value="devices">Devices</option>
                <option value="vulnerabilities">Vulnerabilities</option>
                <option value="services">Services</option>
              </select>
            </div>

            <div className="pt-2 border-t">
              <div className="text-xs text-gray-600 space-y-1">
                <div>Nodes: {filteredNodes.length}</div>
                <div>Connections: {filteredEdges.length}</div>
              </div>
            </div>
          </div>
        </Panel>

        <Panel position="top-right" className="bg-white p-2 rounded-lg shadow-lg m-4">
          <div className="flex gap-2">
            <button
              onClick={exportData}
              className="p-2 hover:bg-gray-100 rounded transition-colors"
              title="Export"
            >
              <Download className="w-4 h-4" />
            </button>
            <button
              onClick={() => window.location.reload()}
              className="p-2 hover:bg-gray-100 rounded transition-colors"
              title="Refresh"
            >
              <RefreshCw className="w-4 h-4" />
            </button>
          </div>
        </Panel>

        {selectedNode && (
          <Panel position="bottom-right" className="bg-white p-4 rounded-lg shadow-lg m-4 max-w-md border border-gray-200">
            <div className="flex items-center justify-between mb-3">
              <h4 className="font-semibold text-gray-900">Node Details</h4>
              <button
                onClick={() => setSelectedNode(null)}
                className="text-gray-400 hover:text-gray-600 transition-colors"
                title="Close"
              >
                <XCircle className="w-4 h-4" />
              </button>
            </div>
            <div className="space-y-2 text-sm">
              <div>
                <span className="font-medium">Type:</span> {selectedNode.type}
              </div>
              <div>
                <span className="font-medium">Label:</span> {selectedNode.data.label}
              </div>
              {selectedNode.data.ipAddress && (
                <div>
                  <span className="font-medium">IP:</span> {selectedNode.data.ipAddress}
                </div>
              )}
              {selectedNode.data.riskScore !== undefined && (
                <div>
                  <span className="font-medium">Risk Score:</span> {selectedNode.data.riskScore.toFixed(1)}
                </div>
              )}
              {selectedNode.data.severity && (
                <div>
                  <span className="font-medium">Severity:</span>{' '}
                  <span className="capitalize">{selectedNode.data.severity}</span>
                </div>
              )}
            </div>
          </Panel>
        )}
      </ReactFlow>
    </div>
  );
};

export default NetworkFlowVisualizer;

