import React, { useEffect, useRef, useState, useCallback } from 'react';
import * as d3 from 'd3';

interface TopologyNode {
  id: string;
  type: 'agent' | 'asset' | 'amass_discovery';
  assetId?: string;
  agentId?: string;
  name: string;
  ipAddress: string;
  location: string;
  clusterId?: string;
  riskScore: number;
  status: 'active' | 'inactive' | 'discovered';
  metadata: Record<string, any>;
}

interface TopologyLink {
  source: string;
  target: string;
  type: 'network' | 'scan' | 'external';
  strength: number;
  protocol?: string;
  port?: number;
}

interface NetworkTopology {
  id: string;
  companyId: string;
  nodes: TopologyNode[];
  links: TopologyLink[];
  clusters: any[];
  lastUpdated: string;
}

interface NetworkTopologyProps {
  data: NetworkTopology;
  onNodeClick?: (node: TopologyNode) => void;
  onRefresh?: () => void;
  viewMode?: 'network' | 'floor' | 'geographic' | 'cluster';
  agentFilter?: 'all' | 'vulnerability' | 'endpoint' | 'network' | 'amass';
  connectionFilter?: 'all' | 'connected' | 'disconnected' | 'intermittent';
}

const NetworkTopology: React.FC<NetworkTopologyProps> = ({
  data,
  onNodeClick,
  onRefresh,
  viewMode = 'network',
  agentFilter = 'all',
  connectionFilter = 'all'
}) => {
  const svgRef = useRef<SVGSVGElement>(null);
  const [tooltip, setTooltip] = useState<{ show: boolean; content: string; x: number; y: number }>({
    show: false,
    content: '',
    x: 0,
    y: 0
  });

  const getNodeColor = useCallback((node: TopologyNode) => {
    const colors = {
      agent: '#64ffda',
      asset: '#ff6b6b',
      amass_discovery: '#45b7d1'
    };
    return colors[node.type] || '#ffa500';
  }, []);

  const getNodeSize = useCallback((node: TopologyNode) => {
    const baseSizes = {
      agent: 15,
      asset: 12,
      amass_discovery: 10
    };
    const riskMultiplier = Math.min(node.riskScore * 0.3, 10);
    return (baseSizes[node.type] || 10) + riskMultiplier;
  }, []);

  const getLinkColor = useCallback((link: TopologyLink) => {
    const colors = {
      scan: '#64ffda',
      network: '#4ecdc4',
      external: '#45b7d1'
    };
    return colors[link.type] || '#666';
  }, []);

  const showTooltip = useCallback((event: React.MouseEvent, node: TopologyNode) => {
    const content = `
      <strong>${node.name}</strong><br>
      Type: ${node.type}<br>
      Status: ${node.status}<br>
      IP: ${node.ipAddress}<br>
      Risk Score: ${node.riskScore.toFixed(1)}<br>
      Location: ${node.location}<br>
      Last Updated: ${new Date().toLocaleString()}
    `;

    setTooltip({
      show: true,
      content,
      x: event.pageX + 10,
      y: event.pageY - 10
    });
  }, []);

  const hideTooltip = useCallback(() => {
    setTooltip(prev => ({ ...prev, show: false }));
  }, []);

  const handleNodeClick = useCallback((node: TopologyNode) => {
    if (onNodeClick) {
      onNodeClick(node);
    }
  }, [onNodeClick]);

  useEffect(() => {
    if (!svgRef.current || !data) return;

    const svg = d3.select(svgRef.current);
    const width = window.innerWidth;
    const height = window.innerHeight;

    // Clear previous content
    svg.selectAll("*").remove();

    // Add glow filter
    const defs = svg.append("defs");
    const glowFilter = defs.append("filter").attr("id", "glow");
    glowFilter.append("feGaussianBlur").attr("stdDeviation", "3").attr("result", "coloredBlur");
    const feMerge = glowFilter.append("feMerge");
    feMerge.append("feMergeNode").attr("in", "coloredBlur");
    feMerge.append("feMergeNode").attr("in", "SourceGraphic");

    // Create force simulation
    const simulation = d3.forceSimulation(data.nodes as any)
      .force("link", d3.forceLink(data.links).id((d: any) => d.id).distance(100))
      .force("charge", d3.forceManyBody().strength(-300))
      .force("center", d3.forceCenter(width / 2, height / 2))
      .force("collision", d3.forceCollide().radius(30));

    // Draw links
    const link = svg.append("g")
      .selectAll("line")
      .data(data.links)
      .enter().append("line")
      .attr("class", "link")
      .attr("stroke", (d: any) => getLinkColor(d))
      .attr("stroke-width", (d: any) => Math.sqrt(d.strength) * 2)
      .style("stroke-opacity", 0.6)
      .style("transition", "stroke-width 0.3s ease");

    // Draw nodes
    const node = svg.append("g")
      .selectAll("circle")
      .data(data.nodes)
      .enter().append("circle")
      .attr("class", "node")
      .attr("r", (d: any) => getNodeSize(d))
      .attr("fill", (d: any) => getNodeColor(d))
      .attr("stroke", "#fff")
      .attr("stroke-width", 1.5)
      .attr("filter", "url(#glow)")
      .style("cursor", "pointer")
      .style("transition", "all 0.3s ease")
      .on("mouseover", function(event, d) {
        d3.select(this).style("stroke-width", "3px");
        showTooltip(event, d);
      })
      .on("mouseout", function() {
        d3.select(this).style("stroke-width", "1.5px");
        hideTooltip();
      })
      .on("click", function(d) {
        handleNodeClick(d);
      })
      .call(d3.drag()
        .on("start", dragstarted)
        .on("drag", dragged)
        .on("end", dragended) as any);

    // Add labels
    const labels = svg.append("g")
      .selectAll("text")
      .data(data.nodes)
      .enter().append("text")
      .attr("class", "node-label")
      .text((d: any) => d.name)
      .attr("dy", 25)
      .style("font-size", "10px")
      .style("fill", "#fff")
      .style("text-anchor", "middle")
      .style("pointer-events", "none");

    // Update positions on simulation tick
    simulation.on("tick", () => {
      link
        .attr("x1", (d: any) => d.source.x)
        .attr("y1", (d: any) => d.source.y)
        .attr("x2", (d: any) => d.target.x)
        .attr("y2", (d: any) => d.target.y);

      node
        .attr("cx", (d: any) => d.x)
        .attr("cy", (d: any) => d.y);

      labels
        .attr("x", (d: any) => d.x)
        .attr("y", (d: any) => d.y);
    });

    function dragstarted(event: any, d: any) {
      if (!event.active) simulation.alphaTarget(0.3).restart();
      d.fx = d.x;
      d.fy = d.y;
    }

    function dragged(event: any, d: any) {
      d.fx = event.x;
      d.fy = event.y;
    }

    function dragended(event: any, d: any) {
      if (!event.active) simulation.alphaTarget(0);
      d.fx = null;
      d.fy = null;
    }

    // Handle window resize
    const handleResize = () => {
      const newWidth = window.innerWidth;
      const newHeight = window.innerHeight;
      svg.attr("width", newWidth).attr("height", newHeight);
      simulation.force("center", d3.forceCenter(newWidth / 2, newHeight / 2));
      simulation.restart();
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      simulation.stop();
    };
  }, [data, getNodeColor, getNodeSize, getLinkColor, showTooltip, hideTooltip, handleNodeClick]);

  return (
    <div className="relative w-full h-full">
      <svg
        ref={svgRef}
        width={window.innerWidth}
        height={window.innerHeight}
        className="absolute top-0 left-0"
      />
      
      {/* Tooltip */}
      {tooltip.show && (
        <div
          className="absolute z-50 bg-black bg-opacity-90 text-white p-3 rounded text-xs border border-white border-opacity-20 max-w-xs"
          style={{
            left: tooltip.x,
            top: tooltip.y,
            pointerEvents: 'none'
          }}
          dangerouslySetInnerHTML={{ __html: tooltip.content }}
        />
      )}

      {/* Controls */}
      <div className="absolute top-4 right-4 z-10 bg-black bg-opacity-90 p-4 rounded-lg border border-white border-opacity-20">
        <div className="space-y-3">
          <div>
            <label className="block text-cyan-400 text-xs font-semibold uppercase mb-1">
              View Mode
            </label>
            <select 
              className="w-full p-2 bg-white bg-opacity-10 border border-white border-opacity-20 rounded text-white text-sm"
              value={viewMode}
              onChange={(e) => console.log('View mode changed:', e.target.value)}
            >
              <option value="network">Network Topology</option>
              <option value="floor">Floor Plan View</option>
              <option value="geographic">Geographic View</option>
              <option value="cluster">Agent Clusters</option>
            </select>
          </div>

          <div>
            <label className="block text-cyan-400 text-xs font-semibold uppercase mb-1">
              Filter by Agent Type
            </label>
            <select 
              className="w-full p-2 bg-white bg-opacity-10 border border-white border-opacity-20 rounded text-white text-sm"
              value={agentFilter}
              onChange={(e) => console.log('Agent filter changed:', e.target.value)}
            >
              <option value="all">All Agents</option>
              <option value="vulnerability">Vulnerability Scanners</option>
              <option value="endpoint">Endpoint Agents</option>
              <option value="network">Network Discovery</option>
              <option value="amass">OWASP Amass</option>
            </select>
          </div>

          <div>
            <label className="block text-cyan-400 text-xs font-semibold uppercase mb-1">
              Connection Status
            </label>
            <select 
              className="w-full p-2 bg-white bg-opacity-10 border border-white border-opacity-20 rounded text-white text-sm"
              value={connectionFilter}
              onChange={(e) => console.log('Connection filter changed:', e.target.value)}
            >
              <option value="all">All Connections</option>
              <option value="connected">Connected Only</option>
              <option value="disconnected">Disconnected Only</option>
              <option value="intermittent">Intermittent</option>
            </select>
          </div>

          <div className="space-y-2">
            <button
              onClick={onRefresh}
              className="w-full p-2 bg-gradient-to-r from-cyan-400 to-teal-400 text-black font-semibold rounded hover:transform hover:-translate-y-0.5 hover:shadow-lg transition-all duration-300"
            >
               Refresh Data
            </button>
            
            <button
              onClick={() => {
                const exportData = {
                  timestamp: new Date().toISOString(),
                  topology: data
                };
                const blob = new Blob([JSON.stringify(exportData, null, 2)], {type: "application/json"});
                const url = URL.createObjectURL(blob);
                const a = document.createElement("a");
                a.href = url;
                a.download = `network_topology_${Date.now()}.json`;
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
                URL.revokeObjectURL(url);
              }}
              className="w-full p-2 bg-gradient-to-r from-blue-400 to-purple-400 text-white font-semibold rounded hover:transform hover:-translate-y-0.5 hover:shadow-lg transition-all duration-300"
            >
               Export Topology
            </button>
          </div>
        </div>
      </div>

      {/* Stats */}
      <div className="absolute bottom-4 left-4 z-10 bg-black bg-opacity-90 p-4 rounded-lg border border-white border-opacity-20">
        <div className="space-y-2 text-xs">
          <div className="flex justify-between">
            <span className="text-gray-300">Total Agents:</span>
            <span className="text-cyan-400 font-semibold">{data.nodes.filter(n => n.type === 'agent').length}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-300">Active Connections:</span>
            <span className="text-cyan-400 font-semibold">{data.links.length}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-300">Discovered Assets:</span>
            <span className="text-cyan-400 font-semibold">{data.nodes.filter(n => n.type === 'amass_discovery').length}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-300">Risk Score (Avg):</span>
            <span className="text-cyan-400 font-semibold">
              {(data.nodes.reduce((sum, n) => sum + n.riskScore, 0) / data.nodes.length).toFixed(1)}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-300">Last Update:</span>
            <span className="text-cyan-400 font-semibold">{new Date().toLocaleTimeString()}</span>
          </div>
        </div>
      </div>

      {/* Legend */}
      <div className="absolute bottom-4 right-4 z-10 bg-black bg-opacity-90 p-4 rounded-lg border border-white border-opacity-20">
        <div className="space-y-2 text-xs">
          <div className="flex items-center">
            <div className="w-3 h-3 rounded-full bg-cyan-400 mr-2"></div>
            <span>Vulnerability Agent</span>
          </div>
          <div className="flex items-center">
            <div className="w-3 h-3 rounded-full bg-red-400 mr-2"></div>
            <span>Endpoint Agent</span>
          </div>
          <div className="flex items-center">
            <div className="w-3 h-3 rounded-full bg-teal-400 mr-2"></div>
            <span>Network Device</span>
          </div>
          <div className="flex items-center">
            <div className="w-3 h-3 rounded-full bg-blue-400 mr-2"></div>
            <span>OWASP Amass Discovery</span>
          </div>
          <div className="flex items-center">
            <div className="w-3 h-3 rounded-full bg-orange-400 mr-2"></div>
            <span>Critical Asset</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NetworkTopology;
