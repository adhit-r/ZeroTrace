import React, { useState } from 'react';
import NetworkFlowVisualizer from '../components/network/NetworkFlowVisualizer';
import type { Node, Edge } from 'reactflow';
import { Network, Info } from 'lucide-react';

const Topology: React.FC = () => {
  const [selectedNode, setSelectedNode] = useState<Node | null>(null);

  const handleNodeClick = (node: Node) => {
    setSelectedNode(node);
    console.log('Node clicked:', node);
  };

  const handleEdgeClick = (edge: Edge) => {
    console.log('Edge clicked:', edge);
  };

  return (
    <div className="h-screen w-full flex flex-col bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b px-6 py-4 flex items-center justify-between shadow-sm">
        <div className="flex items-center gap-3">
          <Network className="w-6 h-6 text-blue-500" />
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Network Topology</h1>
            <p className="text-sm text-gray-600">Interactive n8n-style network device visualization</p>
          </div>
        </div>
        {selectedNode && (
          <div className="flex items-center gap-2 text-sm text-gray-600 bg-blue-50 px-4 py-2 rounded-lg">
            <Info className="w-4 h-4 text-blue-500" />
            <span className="font-medium">Selected: {selectedNode.data.label}</span>
          </div>
        )}
      </div>

      {/* Main Visualizer */}
      <div className="flex-1 relative">
        <NetworkFlowVisualizer
          onNodeClick={handleNodeClick}
          onEdgeClick={handleEdgeClick}
        />
      </div>
    </div>
  );
};

export default Topology;
