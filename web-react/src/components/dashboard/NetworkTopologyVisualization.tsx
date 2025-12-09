import React, { useState, useEffect, useRef, useCallback } from 'react';
import {
  Network,
  RefreshCw,
  Filter,
  ZoomIn,
  ZoomOut,
  AlertTriangle,
  Search,
  Download,
  RotateCcw
} from 'lucide-react';
import * as d3 from 'd3';

interface NetworkNode {
  id: string;
  label: string;
  type: 'host' | 'service' | 'vulnerability' | 'network' | 'firewall' | 'router';
  x: number;
  y: number;
  color: string;
  size: number;
  data: any;
  status: 'online' | 'offline' | 'unknown';
  riskScore: number;
  vulnerabilities: number;
  services: number;
}

interface NetworkEdge {
  source: string;
  target: string;
  type: 'connection' | 'vulnerability' | 'service' | 'firewall_rule' | 'routing';
  color: string;
  width: number;
  data: any;
  status: 'active' | 'blocked' | 'monitored';
}

interface NetworkTopology {
  nodes: NetworkNode[];
  edges: NetworkEdge[];
}

interface NetworkTopologyVisualizationProps {
  className?: string;
  data?: NetworkTopology;
  onNodeClick?: (node: NetworkNode) => void;
  onEdgeClick?: (edge: NetworkEdge) => void;
}

const NetworkTopologyVisualization: React.FC<NetworkTopologyVisualizationProps> = ({
  className = '',
  data,
  onNodeClick,
  onEdgeClick
}) => {
  const svgRef = useRef<SVGSVGElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedNode, setSelectedNode] = useState<NetworkNode | null>(null);
  const [selectedEdge, setSelectedEdge] = useState<NetworkEdge | null>(null);
  const [zoom, setZoom] = useState(1);
  const [, setPan] = useState({ x: 0, y: 0 });
  const [layout, setLayout] = useState<'force' | 'hierarchical' | 'circular'>('force');
  const [filter, setFilter] = useState<string>('all');
  const [searchTerm, setSearchTerm] = useState('');

  // D3 simulation and zoom behavior
  const simulationRef = useRef<d3.Simulation<NetworkNode, NetworkEdge> | null>(null);
  const zoomBehaviorRef = useRef<d3.ZoomBehavior<SVGSVGElement, unknown> | null>(null);

  useEffect(() => {
    const loadNetworkData = async () => {
      try {
        setIsLoading(true);
        setError(null);

        if (!data) {
          // Fetch network topology data
          const response = await fetch('/api/v2/network/topology');
          if (!response.ok) {
            throw new Error('Failed to fetch network topology data');
          }
          const topologyData = await response.json();
          renderTopology(topologyData);
        } else {
          renderTopology(data);
        }
      } catch (err) {
        console.error('Failed to load network topology:', err);
        setError('Failed to load network topology data');
      } finally {
        setIsLoading(false);
      }
    };

    loadNetworkData();
  }, [data]);

  const renderTopology = useCallback((topologyData: NetworkTopology) => {
    if (!svgRef.current || !containerRef.current) return;

    const svg = d3.select(svgRef.current);
    const container = containerRef.current;
    const { width, height } = container.getBoundingClientRect();

    // Clear previous content
    svg.selectAll('*').remove();

    // Set up SVG
    svg
      .attr('width', width)
      .attr('height', height)
      .attr('viewBox', `0 0 ${width} ${height}`);

    // Create main group for zoom/pan
    const g = svg.append('g');

    // Set up zoom behavior
    const zoomBehavior = d3.zoom<SVGSVGElement, unknown>()
      .scaleExtent([0.1, 4])
      .on('zoom', (event) => {
        const { transform } = event;
        g.attr('transform', transform);
        setZoom(transform.k);
        setPan({ x: transform.x, y: transform.y });
      });

    svg.call(zoomBehavior);
    zoomBehaviorRef.current = zoomBehavior;

    // Filter data based on current filter
    const filteredNodes = topologyData.nodes.filter(node => {
      if (filter === 'all') return true;
      if (filter === 'hosts') return node.type === 'host';
      if (filter === 'services') return node.type === 'service';
      if (filter === 'vulnerabilities') return node.type === 'vulnerability';
      if (filter === 'network') return node.type === 'network' || node.type === 'firewall' || node.type === 'router';
      return true;
    });

    const filteredEdges = topologyData.edges.filter(edge => {
      const sourceNode = filteredNodes.find(n => n.id === edge.source);
      const targetNode = filteredNodes.find(n => n.id === edge.target);
      return sourceNode && targetNode;
    });

    // Set up force simulation
    const simulation = d3.forceSimulation<NetworkNode>(filteredNodes)
      .force('link', d3.forceLink<NetworkNode, NetworkEdge>(filteredEdges)
        .id(d => d.id)
        .distance(100)
        .strength(0.1)
      )
      .force('charge', d3.forceManyBody().strength(-300))
      .force('center', d3.forceCenter(width / 2, height / 2))
      .force('collision', d3.forceCollide().radius(30));

    simulationRef.current = simulation;

    // Create edges
    const edges = g.append('g')
      .attr('class', 'edges')
      .selectAll('line')
      .data(filteredEdges)
      .enter()
      .append('line')
      .attr('stroke', d => d.color)
      .attr('stroke-width', d => d.width)
      .attr('stroke-opacity', 0.6)
      .attr('stroke-dasharray', d => d.type === 'firewall_rule' ? '5,5' : 'none')
      .on('click', (event, d) => {
        event.stopPropagation();
        setSelectedEdge(d);
        onEdgeClick?.(d);
      })
      .on('mouseover', function (event, d) {
        d3.select(this)
          .attr('stroke-width', d.width * 2)
          .attr('stroke-opacity', 1);

        // Show tooltip
        const tooltip = d3.select('body').append('div')
          .attr('class', 'network-tooltip')
          .style('position', 'absolute')
          .style('background', 'rgba(0, 0, 0, 0.8)')
          .style('color', 'white')
          .style('padding', '8px')
          .style('border-radius', '4px')
          .style('font-size', '12px')
          .style('pointer-events', 'none')
          .style('z-index', '1000');

        tooltip.html(`
          <div><strong>${d.type.toUpperCase()}</strong></div>
          <div>Source: ${d.source}</div>
          <div>Target: ${d.target}</div>
          <div>Status: ${d.status}</div>
        `);

        tooltip.style('left', (event.pageX + 10) + 'px')
          .style('top', (event.pageY - 10) + 'px');
      })
      .on('mouseout', function () {
        d3.select(this)
          .attr('stroke-width', (d: any) => d.width)
          .attr('stroke-opacity', 0.6);

        d3.selectAll('.network-tooltip').remove();
      });

    // Create nodes
    const nodes = g.append('g')
      .attr('class', 'nodes')
      .selectAll('g')
      .data(filteredNodes)
      .enter()
      .append('g')
      .attr('class', 'node')
      .on('click', (event, d) => {
        event.stopPropagation();
        setSelectedNode(d);
        onNodeClick?.(d);
      })
      .on('mouseover', function (event, d) {
        // Highlight connected edges
        edges.attr('stroke-opacity', edge =>
          edge.source === d.id || edge.target === d.id ? 1 : 0.1
        );

        // Show tooltip
        const tooltip = d3.select('body').append('div')
          .attr('class', 'network-tooltip')
          .style('position', 'absolute')
          .style('background', 'rgba(0, 0, 0, 0.8)')
          .style('color', 'white')
          .style('padding', '8px')
          .style('border-radius', '4px')
          .style('font-size', '12px')
          .style('pointer-events', 'none')
          .style('z-index', '1000');

        tooltip.html(`
          <div><strong>${d.label}</strong></div>
          <div>Type: ${d.type}</div>
          <div>Status: ${d.status}</div>
          <div>Risk Score: ${d.riskScore}</div>
          <div>Vulnerabilities: ${d.vulnerabilities}</div>
          <div>Services: ${d.services}</div>
        `);

        tooltip.style('left', (event.pageX + 10) + 'px')
          .style('top', (event.pageY - 10) + 'px');
      })
      .on('mouseout', function () {
        // Reset edge opacity
        edges.attr('stroke-opacity', 0.6);

        d3.selectAll('.network-tooltip').remove();
      });

    // Add node circles
    nodes.append('circle')
      .attr('r', d => d.size)
      .attr('fill', d => d.color)
      .attr('stroke', '#000')
      .attr('stroke-width', 2)
      .attr('opacity', 0.8);

    // Add node icons
    nodes.append('text')
      .attr('text-anchor', 'middle')
      .attr('dy', '0.35em')
      .attr('font-size', '12px')
      .attr('font-weight', 'bold')
      .attr('fill', 'white')
      .text(d => {
        switch (d.type) {
          case 'host': return 'H';
          case 'service': return 'S';
          case 'vulnerability': return 'V';
          case 'network': return 'N';
          case 'firewall': return 'F';
          case 'router': return 'R';
          default: return '●';
        }
      });

    // Add node labels
    nodes.append('text')
      .attr('text-anchor', 'middle')
      .attr('dy', d => d.size + 15)
      .attr('font-size', '10px')
      .attr('font-weight', 'bold')
      .attr('fill', '#000')
      .text(d => d.label.length > 10 ? d.label.substring(0, 10) + '...' : d.label);

    // Update positions on simulation tick
    simulation.on('tick', () => {
      edges
        .attr('x1', d => ((d.source as unknown) as NetworkNode).x!)
        .attr('y1', d => ((d.source as unknown) as NetworkNode).y!)
        .attr('x2', d => ((d.target as unknown) as NetworkNode).x!)
        .attr('y2', d => ((d.target as unknown) as NetworkNode).y!);

      nodes
        .attr('transform', d => `translate(${d.x},${d.y})`);
    });

    // Add legend
    const legend = g.append('g')
      .attr('class', 'legend')
      .attr('transform', `translate(20, 20)`);

    const legendItems = [
      { type: 'host', label: 'Host', color: '#3B82F6' },
      { type: 'service', label: 'Service', color: '#10B981' },
      { type: 'vulnerability', label: 'Vulnerability', color: '#EF4444' },
      { type: 'network', label: 'Network', color: '#8B5CF6' },
      { type: 'firewall', label: 'Firewall', color: '#F59E0B' },
      { type: 'router', label: 'Router', color: '#06B6D4' }
    ];

    legendItems.forEach((item, index) => {
      const legendItem = legend.append('g')
        .attr('transform', `translate(0, ${index * 25})`);

      legendItem.append('circle')
        .attr('r', 8)
        .attr('fill', item.color)
        .attr('stroke', '#000')
        .attr('stroke-width', 1);

      legendItem.append('text')
        .attr('x', 20)
        .attr('y', 5)
        .attr('font-size', '12px')
        .attr('font-weight', 'bold')
        .text(item.label);
    });

  }, [data, filter, onNodeClick, onEdgeClick]);

  const handleZoomIn = () => {
    if (zoomBehaviorRef.current && svgRef.current) {
      d3.select(svgRef.current).transition().call(
        zoomBehaviorRef.current.scaleBy, 1.5
      );
    }
  };

  const handleZoomOut = () => {
    if (zoomBehaviorRef.current && svgRef.current) {
      d3.select(svgRef.current).transition().call(
        zoomBehaviorRef.current.scaleBy, 1 / 1.5
      );
    }
  };

  const handleReset = () => {
    if (zoomBehaviorRef.current && svgRef.current) {
      d3.select(svgRef.current).transition().call(
        zoomBehaviorRef.current.transform,
        d3.zoomIdentity
      );
    }
  };

  const handleLayoutChange = (newLayout: 'force' | 'hierarchical' | 'circular') => {
    setLayout(newLayout);
    // Re-render with new layout
    if (data) {
      renderTopology(data);
    }
  };

  if (isLoading) {
    return (
      <div className={`p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal ${className}`}>
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded border-2 border-black"></div>
          <div className="h-96 bg-gray-200 rounded border-2 border-black"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal ${className}`}>
        <div className="text-center py-8">
          <AlertTriangle className="h-12 w-12 text-red-500 mx-auto mb-4" />
          <h3 className="text-lg font-bold text-red-800 mb-2">Error Loading Network Topology</h3>
          <p className="text-red-600 mb-4">{error}</p>
          <button
            onClick={() => window.location.reload()}
            className="px-4 py-2 bg-red-100 text-red-800 border-2 border-red-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-red-200 transition-colors"
          >
            <RefreshCw className="h-4 w-4 mr-2 inline-block" />
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className={`bg-white border-3 border-black rounded-lg shadow-neo-brutal ${className}`}>
      {/* Header */}
      <div className="p-6 border-b-3 border-black">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="p-3 bg-blue-100 rounded border-2 border-black">
              <Network className="h-8 w-8 text-blue-600" />
            </div>
            <div>
              <h1 className="text-3xl font-black uppercase tracking-wider text-black">
                Network Topology Visualization
              </h1>
              <p className="text-gray-600">Interactive network discovery and vulnerability mapping</p>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <button className="px-4 py-2 bg-gray-100 text-gray-800 border-2 border-gray-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-gray-200 transition-colors">
              <Filter className="h-4 w-4 mr-2 inline-block" />
              Filter
            </button>
            <button className="px-4 py-2 bg-orange-100 text-orange-800 border-2 border-orange-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-orange-200 transition-colors">
              <Download className="h-4 w-4 mr-2 inline-block" />
              Export
            </button>
          </div>
        </div>
      </div>

      {/* Controls */}
      <div className="p-4 bg-gray-50 border-b-3 border-black">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            {/* Layout Selection */}
            <div className="flex items-center gap-2">
              <label className="text-sm font-bold uppercase tracking-wider text-black">Layout:</label>
              <select
                value={layout}
                onChange={(e) => handleLayoutChange(e.target.value as any)}
                className="px-3 py-1 border-2 border-black rounded text-sm font-bold"
              >
                <option value="force">Force Directed</option>
                <option value="hierarchical">Hierarchical</option>
                <option value="circular">Circular</option>
              </select>
            </div>

            {/* Filter Selection */}
            <div className="flex items-center gap-2">
              <label className="text-sm font-bold uppercase tracking-wider text-black">Filter:</label>
              <select
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
                className="px-3 py-1 border-2 border-black rounded text-sm font-bold"
              >
                <option value="all">All</option>
                <option value="hosts">Hosts</option>
                <option value="services">Services</option>
                <option value="vulnerabilities">Vulnerabilities</option>
                <option value="network">Network Devices</option>
              </select>
            </div>

            {/* Search */}
            <div className="flex items-center gap-2">
              <Search className="h-4 w-4 text-gray-600" />
              <input
                type="text"
                placeholder="Search nodes..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="px-3 py-1 border-2 border-black rounded text-sm"
              />
            </div>
          </div>

          <div className="flex items-center gap-2">
            {/* Zoom Controls */}
            <button
              onClick={handleZoomIn}
              className="p-2 bg-gray-100 text-gray-800 border-2 border-gray-300 rounded hover:bg-gray-200 transition-colors"
            >
              <ZoomIn className="h-4 w-4" />
            </button>
            <button
              onClick={handleZoomOut}
              className="p-2 bg-gray-100 text-gray-800 border-2 border-gray-300 rounded hover:bg-gray-200 transition-colors"
            >
              <ZoomOut className="h-4 w-4" />
            </button>
            <button
              onClick={handleReset}
              className="p-2 bg-gray-100 text-gray-800 border-2 border-gray-300 rounded hover:bg-gray-200 transition-colors"
            >
              <RotateCcw className="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>

      {/* Visualization */}
      <div className="p-6">
        <div
          ref={containerRef}
          className="w-full h-96 bg-gray-50 border-2 border-black rounded relative"
        >
          <svg
            ref={svgRef}
            className="w-full h-full"
            style={{ cursor: 'grab' }}
          />

          {/* Zoom indicator */}
          <div className="absolute top-4 right-4 bg-white border-2 border-black rounded p-2">
            <div className="text-sm font-bold text-black">
              Zoom: {Math.round(zoom * 100)}%
            </div>
          </div>
        </div>
      </div>

      {/* Selected Node/Edge Details */}
      {(selectedNode || selectedEdge) && (
        <div className="p-6 border-t-3 border-black">
          <h3 className="text-lg font-bold uppercase tracking-wider text-black mb-4">
            {selectedNode ? 'Node Details' : 'Edge Details'}
          </h3>
          {selectedNode && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <h4 className="font-bold text-black mb-2">Basic Information</h4>
                <div className="space-y-2 text-sm">
                  <div><strong>Label:</strong> {selectedNode.label}</div>
                  <div><strong>Type:</strong> {selectedNode.type}</div>
                  <div><strong>Status:</strong> {selectedNode.status}</div>
                  <div><strong>Risk Score:</strong> {selectedNode.riskScore}</div>
                </div>
              </div>
              <div>
                <h4 className="font-bold text-black mb-2">Security Metrics</h4>
                <div className="space-y-2 text-sm">
                  <div><strong>Vulnerabilities:</strong> {selectedNode.vulnerabilities}</div>
                  <div><strong>Services:</strong> {selectedNode.services}</div>
                  <div><strong>Position:</strong> ({Math.round(selectedNode.x)}, {Math.round(selectedNode.y)})</div>
                </div>
              </div>
            </div>
          )}
          {selectedEdge && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <h4 className="font-bold text-black mb-2">Connection Information</h4>
                <div className="space-y-2 text-sm">
                  <div><strong>Type:</strong> {selectedEdge.type}</div>
                  <div><strong>Source:</strong> {selectedEdge.source}</div>
                  <div><strong>Target:</strong> {selectedEdge.target}</div>
                  <div><strong>Status:</strong> {selectedEdge.status}</div>
                </div>
              </div>
              <div>
                <h4 className="font-bold text-black mb-2">Visual Properties</h4>
                <div className="space-y-2 text-sm">
                  <div><strong>Color:</strong> {selectedEdge.color}</div>
                  <div><strong>Width:</strong> {selectedEdge.width}px</div>
                </div>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default NetworkTopologyVisualization;
