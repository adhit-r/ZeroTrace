import React from 'react';
import NetworkFlowVisualizer from '../components/network/NetworkFlowVisualizer';
import type { Node, Edge } from 'reactflow';

const NetworkTopology: React.FC = () => {
  const handleNodeClick = (node: Node) => {
    console.log('Node clicked:', node);
    // You can add modal or detail view here
  };

  const handleEdgeClick = (edge: Edge) => {
    console.log('Edge clicked:', edge);
  };

  return (
    <div className="h-screen w-full">
      <div className="h-full">
        <NetworkFlowVisualizer
          onNodeClick={handleNodeClick}
          onEdgeClick={handleEdgeClick}
        />
      </div>
    </div>
  );
};

export default NetworkTopology;

