import React, { useState } from 'react';
import AttackPathVisualizer from '../components/network/AttackPathVisualizer';
import { Target, Info } from 'lucide-react';

const AttackPaths: React.FC = () => {
  const [selectedPath, setSelectedPath] = useState<any>(null);

  const handlePathSelect = (path: any) => {
    setSelectedPath(path);
    console.log('Attack path selected:', path);
  };

  return (
    <div className="h-screen w-full flex flex-col bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b px-6 py-4 flex items-center justify-between shadow-sm">
        <div className="flex items-center gap-3">
          <Target className="w-6 h-6 text-red-500" />
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Attack Path Analysis</h1>
            <p className="text-sm text-gray-600">NodeZero-style exploit chain visualization</p>
          </div>
        </div>
        {selectedPath && (
          <div className="flex items-center gap-2 text-sm text-gray-600 bg-red-50 px-4 py-2 rounded-lg">
            <Info className="w-4 h-4 text-red-500" />
            <span className="font-medium">Path: {selectedPath.name}</span>
          </div>
        )}
      </div>

      {/* Main Visualizer */}
      <div className="flex-1 relative">
        <AttackPathVisualizer onPathSelect={handlePathSelect} />
      </div>
    </div>
  );
};

export default AttackPaths;

