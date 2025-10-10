import React, { useState } from 'react';
import { MapPin, Building2, AlertTriangle, Shield, TrendingUp } from 'lucide-react';

interface HeatmapData {
  id: string;
  name: string;
  location: string;
  coordinates: { lat: number; lng: number };
  riskScore: number;
  criticalVulns: number;
  totalAssets: number;
  complianceScore: number;
  lastScan: string;
  status: 'active' | 'inactive' | 'maintenance';
}

interface RiskHeatmapProps {
  data: HeatmapData[];
  selectedBranchId?: string;
  onBranchSelect: (branchId: string) => void;
  className?: string;
}

const RiskHeatmap: React.FC<RiskHeatmapProps> = ({
  data,
  selectedBranchId,
  onBranchSelect,
  className = ''
}) => {
  const [viewMode, setViewMode] = useState<'map' | 'grid'>('map');
  const [riskFilter, setRiskFilter] = useState<'all' | 'critical' | 'high' | 'medium' | 'low'>('all');

  const getRiskColor = (score: number) => {
    if (score >= 80) return 'bg-red-500';
    if (score >= 60) return 'bg-orange-500';
    if (score >= 40) return 'bg-yellow-500';
    return 'bg-green-500';
  };

  const getRiskLabel = (score: number) => {
    if (score >= 80) return 'Critical';
    if (score >= 60) return 'High';
    if (score >= 40) return 'Medium';
    return 'Low';
  };

  const getRiskIntensity = (score: number) => {
    if (score >= 80) return 'shadow-red-500/50';
    if (score >= 60) return 'shadow-orange-500/50';
    if (score >= 40) return 'shadow-yellow-500/50';
    return 'shadow-green-500/50';
  };

  const filteredData = data.filter(branch => {
    if (riskFilter === 'all') return true;
    const score = branch.riskScore;
    switch (riskFilter) {
      case 'critical': return score >= 80;
      case 'high': return score >= 60 && score < 80;
      case 'medium': return score >= 40 && score < 60;
      case 'low': return score < 40;
      default: return true;
    }
  });

  const MapView = () => (
    <div className="relative h-96 bg-gray-100 border-3 border-black rounded overflow-hidden">
      {/* Simplified world map background */}
      <div className="absolute inset-0 bg-gradient-to-br from-blue-100 to-green-100">
        <div className="absolute top-1/4 left-1/4 w-32 h-20 bg-green-200 rounded opacity-50"></div>
        <div className="absolute top-1/3 right-1/3 w-24 h-16 bg-green-200 rounded opacity-50"></div>
        <div className="absolute bottom-1/4 left-1/3 w-28 h-18 bg-green-200 rounded opacity-50"></div>
      </div>
      
      {/* Branch markers */}
      {filteredData.map((branch) => (
        <button
          key={branch.id}
          onClick={() => onBranchSelect(branch.id)}
          className={`
            absolute transform -translate-x-1/2 -translate-y-1/2
            w-8 h-8 rounded-full border-3 border-black shadow-lg
            hover:scale-110 transition-transform duration-150
            ${getRiskColor(branch.riskScore)}
            ${selectedBranchId === branch.id ? 'ring-4 ring-orange-500' : ''}
          `}
          style={{
            left: `${(branch.coordinates.lng + 180) / 360 * 100}%`,
            top: `${(90 - branch.coordinates.lat) / 180 * 100}%`
          }}
          title={`${branch.name} - Risk: ${branch.riskScore}`}
        >
          <div className="w-full h-full flex items-center justify-center">
            <AlertTriangle className="w-4 h-4 text-white" />
          </div>
        </button>
      ))}
    </div>
  );

  const GridView = () => (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {filteredData.map((branch) => (
        <div
          key={branch.id}
          onClick={() => onBranchSelect(branch.id)}
          className={`
            p-4 rounded border-3 border-black bg-white shadow-lg
            hover:shadow-xl hover:translate-x-1 hover:translate-y-1
            transition-all duration-150 ease-in-out cursor-pointer
            ${selectedBranchId === branch.id ? 'ring-4 ring-orange-500 bg-orange-50' : ''}
          `}
        >
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-2">
              <Building2 className="w-5 h-5" />
              <h3 className="font-bold">{branch.name}</h3>
            </div>
            <div className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider text-white ${getRiskColor(branch.riskScore)}`}>
              {getRiskLabel(branch.riskScore)}
            </div>
          </div>
          
          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-gray-600">Risk Score:</span>
              <span className={`font-bold ${getRiskColor(branch.riskScore).replace('bg-', 'text-')}`}>
                {branch.riskScore}
              </span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-600">Critical Vulns:</span>
              <span className="font-bold text-red-600">{branch.criticalVulns}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-600">Assets:</span>
              <span className="font-bold">{branch.totalAssets}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-600">Compliance:</span>
              <span className="font-bold text-green-600">{branch.complianceScore}%</span>
            </div>
          </div>
          
          <div className="mt-3 pt-3 border-t-2 border-gray-200">
            <div className="flex items-center justify-between text-xs text-gray-500">
              <span>{branch.location}</span>
              <span>{new Date(branch.lastScan).toLocaleDateString()}</span>
            </div>
          </div>
        </div>
      ))}
    </div>
  );

  return (
    <div className={`p-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold uppercase tracking-wider">Risk Heatmap</h2>
          <p className="text-gray-600">Geographic risk distribution across branches</p>
        </div>
        <div className="flex items-center gap-3">
          {/* View Mode Toggle */}
          <div className="flex border-2 border-black rounded">
            <button
              onClick={() => setViewMode('map')}
              className={`px-3 py-1 text-sm font-bold uppercase tracking-wider transition-colors ${
                viewMode === 'map' 
                  ? 'bg-orange-100 text-orange-800 border-r-2 border-black' 
                  : 'bg-white text-gray-600 hover:bg-gray-50'
              }`}
            >
              Map
            </button>
            <button
              onClick={() => setViewMode('grid')}
              className={`px-3 py-1 text-sm font-bold uppercase tracking-wider transition-colors ${
                viewMode === 'grid' 
                  ? 'bg-orange-100 text-orange-800' 
                  : 'bg-white text-gray-600 hover:bg-gray-50'
              }`}
            >
              Grid
            </button>
          </div>
          
          {/* Risk Filter */}
          <select
            value={riskFilter}
            onChange={(e) => setRiskFilter(e.target.value as any)}
            className="px-3 py-1 border-2 border-black rounded text-sm font-bold uppercase tracking-wider focus:outline-none focus:border-orange-500"
          >
            <option value="all">All Risk Levels</option>
            <option value="critical">Critical</option>
            <option value="high">High</option>
            <option value="medium">Medium</option>
            <option value="low">Low</option>
          </select>
        </div>
      </div>

      {/* Legend */}
      <div className="mb-6 p-4 border-3 border-black bg-gray-50 rounded">
        <h3 className="font-bold mb-3">Risk Legend</h3>
        <div className="flex flex-wrap gap-4">
          <div className="flex items-center gap-2">
            <div className="w-4 h-4 bg-red-500 rounded border-2 border-black"></div>
            <span className="text-sm font-bold">Critical (80-100)</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-4 h-4 bg-orange-500 rounded border-2 border-black"></div>
            <span className="text-sm font-bold">High (60-79)</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-4 h-4 bg-yellow-500 rounded border-2 border-black"></div>
            <span className="text-sm font-bold">Medium (40-59)</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-4 h-4 bg-green-500 rounded border-2 border-black"></div>
            <span className="text-sm font-bold">Low (0-39)</span>
          </div>
        </div>
      </div>

      {/* Content */}
      {viewMode === 'map' ? <MapView /> : <GridView />}

      {/* Summary Stats */}
      <div className="mt-6 grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="p-4 border-3 border-black bg-white rounded shadow-lg">
          <div className="flex items-center gap-2 mb-2">
            <AlertTriangle className="w-5 h-5 text-red-600" />
            <span className="font-bold">Critical Branches</span>
          </div>
          <div className="text-2xl font-bold text-red-600">
            {data.filter(b => b.riskScore >= 80).length}
          </div>
        </div>
        <div className="p-4 border-3 border-black bg-white rounded shadow-lg">
          <div className="flex items-center gap-2 mb-2">
            <TrendingUp className="w-5 h-5 text-orange-600" />
            <span className="font-bold">Total Assets</span>
          </div>
          <div className="text-2xl font-bold text-orange-600">
            {data.reduce((sum, b) => sum + b.totalAssets, 0)}
          </div>
        </div>
        <div className="p-4 border-3 border-black bg-white rounded shadow-lg">
          <div className="flex items-center gap-2 mb-2">
            <Shield className="w-5 h-5 text-green-600" />
            <span className="font-bold">Avg Compliance</span>
          </div>
          <div className="text-2xl font-bold text-green-600">
            {data.length > 0 ? Math.round(data.reduce((sum, b) => sum + b.complianceScore, 0) / data.length) : 0}%
          </div>
        </div>
        <div className="p-4 border-3 border-black bg-white rounded shadow-lg">
          <div className="flex items-center gap-2 mb-2">
            <MapPin className="w-5 h-5 text-blue-600" />
            <span className="font-bold">Active Branches</span>
          </div>
          <div className="text-2xl font-bold text-blue-600">
            {data.filter(b => b.status === 'active').length}
          </div>
        </div>
      </div>
    </div>
  );
};

export default RiskHeatmap;
