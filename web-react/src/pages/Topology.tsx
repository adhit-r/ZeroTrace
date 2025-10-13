import React, { useState, useEffect } from 'react';
import { RefreshCw, Search, SlidersHorizontal, Share2, ZoomIn, ZoomOut, Maximize, AlertTriangle, Shield, Activity, Wifi, WifiOff } from 'lucide-react';
import { topologyService } from '../services/topologyService';
import type { TopologyData, TopologyNode } from '../services/topologyService';

const Topology: React.FC = () => {
  const [topologyData, setTopologyData] = useState<TopologyData | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [networkHealth, setNetworkHealth] = useState<any>(null);

  useEffect(() => {
    loadTopologyData();
  }, []);

  const loadTopologyData = async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      const [topology, health] = await Promise.all([
        topologyService.getTopologyData(),
        topologyService.getNetworkHealth()
      ]);
      
      setTopologyData(topology);
      setNetworkHealth(health);
    } catch (err) {
      setError('Failed to load topology data');
      console.error('Error loading topology:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRefresh = () => {
    loadTopologyData();
  };

  const getNodeColor = (node: TopologyNode) => {
    if (node.hasVulns) return 'bg-red-500 border-red-700';
    if (node.status === 'offline') return 'bg-gray-400 border-gray-600';
    
    switch (node.type) {
      case 'router': return 'bg-blue-500 border-blue-700';
      case 'switch': return 'bg-yellow-500 border-yellow-700';
      case 'server': return 'bg-purple-500 border-purple-700';
      case 'workstation': return 'bg-green-500 border-green-700';
      case 'agent': return 'bg-orange-500 border-orange-700';
      default: return 'bg-gray-500 border-gray-700';
    }
  };

  const getRiskColor = (riskScore: number) => {
    if (riskScore >= 8) return 'text-red-600 bg-red-100';
    if (riskScore >= 6) return 'text-orange-600 bg-orange-100';
    if (riskScore >= 4) return 'text-yellow-600 bg-yellow-100';
    return 'text-green-600 bg-green-100';
  };

  return (
    <div className="p-4 md:p-6 lg:p-8 font-sans bg-gray-100 min-h-screen">
      {/* Header */}
      <header className="mb-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-black text-black uppercase tracking-wider">Network Topology</h1>
            <p className="text-sm text-gray-600 font-bold">Visualize your asset connections and security posture.</p>
          </div>
          {networkHealth && (
            <div className="flex gap-4">
              <div className="px-4 py-2 bg-white border-3 border-black rounded-lg shadow-neo-brutal-small">
                <div className="flex items-center gap-2">
                  <Activity className="h-4 w-4 text-green-600" />
                  <span className="font-bold text-green-600">{networkHealth.onlineNodes}</span>
                  <span className="text-sm text-gray-600">Online</span>
                </div>
              </div>
              <div className="px-4 py-2 bg-white border-3 border-black rounded-lg shadow-neo-brutal-small">
                <div className="flex items-center gap-2">
                  <AlertTriangle className="h-4 w-4 text-red-600" />
                  <span className="font-bold text-red-600">{networkHealth.vulnerableNodes}</span>
                  <span className="text-sm text-gray-600">Vulnerable</span>
                </div>
              </div>
              <div className="px-4 py-2 bg-white border-3 border-black rounded-lg shadow-neo-brutal-small">
                <div className="flex items-center gap-2">
                  <Shield className="h-4 w-4 text-blue-600" />
                  <span className="font-bold text-blue-600">{networkHealth.averageRiskScore.toFixed(1)}</span>
                  <span className="text-sm text-gray-600">Avg Risk</span>
                </div>
              </div>
            </div>
          )}
        </div>
      </header>

      {/* Main Content Area */}
      <div className="flex flex-col lg:flex-row gap-6">

        {/* Control Panel */}
        <aside className="lg:w-1/4">
          <div className="bg-white p-4 border-3 border-black rounded-lg shadow-neo-brutal">
            <h2 className="text-lg font-black text-black uppercase mb-4">Controls</h2>
            <div className="space-y-4">
              <button
                onClick={handleRefresh}
                className="w-full flex items-center justify-center gap-2 px-4 py-3 bg-blue-500 text-white font-bold uppercase tracking-wide rounded-lg border-3 border-black shadow-neo-brutal-small hover:shadow-neo-brutal-small-hover transition-all"
              >
                <RefreshCw className={`h-5 w-5 ${isLoading ? 'animate-spin' : ''}`} />
                {isLoading ? 'Refreshing...' : 'Refresh'}
              </button>
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                <input
                  type="text"
                  placeholder="Search assets..."
                  className="w-full pl-10 pr-4 py-3 bg-gray-100 text-black font-bold rounded-lg border-3 border-black focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              <button className="w-full flex items-center justify-center gap-2 px-4 py-3 bg-gray-200 text-black font-bold uppercase tracking-wide rounded-lg border-3 border-black shadow-neo-brutal-small hover:shadow-neo-brutal-small-hover transition-all">
                <SlidersHorizontal className="h-5 w-5" />
                Filters
              </button>
            </div>
          </div>
        </aside>

        {/* Topology Visualization */}
        <main className="flex-1">
          <div className="relative bg-white h-[600px] lg:h-full border-3 border-black rounded-lg shadow-neo-brutal p-4">
            {/* Toolbar */}
            <div className="absolute top-4 right-4 flex gap-2">
              <button className="p-2 bg-gray-100 rounded-lg border-3 border-black shadow-neo-brutal-small"><ZoomIn className="h-5 w-5" /></button>
              <button className="p-2 bg-gray-100 rounded-lg border-3 border-black shadow-neo-brutal-small"><ZoomOut className="h-5 w-5" /></button>
              <button className="p-2 bg-gray-100 rounded-lg border-3 border-black shadow-neo-brutal-small"><Maximize className="h-5 w-5" /></button>
              <button className="p-2 bg-gray-100 rounded-lg border-3 border-black shadow-neo-brutal-small"><Share2 className="h-5 w-5" /></button>
            </div>
            
                   {/* Network Visualization */}
                   <div className="w-full h-full relative">
                     {error ? (
                       <div className="flex items-center justify-center h-full">
                         <div className="text-center">
                           <AlertTriangle className="h-12 w-12 text-red-500 mx-auto mb-4" />
                           <p className="text-red-600 font-bold">{error}</p>
                           <button 
                             onClick={handleRefresh}
                             className="mt-4 px-4 py-2 bg-red-500 text-white border-3 border-black rounded-lg shadow-neo-brutal-small hover:shadow-neo-brutal-small-hover transition-all"
                           >
                             Retry
                           </button>
                         </div>
                       </div>
                     ) : topologyData ? (
                       <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 p-4">
                         {topologyData.nodes.map((node) => (
                           <div
                             key={node.id}
                             className={`p-4 rounded-lg border-3 shadow-neo-brutal-small hover:shadow-neo-brutal-small-hover transition-all cursor-pointer ${getNodeColor(node)}`}
                           >
                             <div className="flex items-center justify-between mb-2">
                               <div className="flex items-center gap-2">
                                 {node.status === 'online' ? (
                                   <Wifi className="h-4 w-4 text-green-600" />
                                 ) : (
                                   <WifiOff className="h-4 w-4 text-gray-400" />
                                 )}
                                 <span className="font-bold text-white text-sm">{node.name}</span>
                               </div>
                               <div className={`px-2 py-1 rounded text-xs font-bold ${getRiskColor(node.riskScore)}`}>
                                 {node.riskScore.toFixed(1)}
                               </div>
                             </div>
                             <div className="text-xs text-white/80 space-y-1">
                               <div>IP: {node.ip}</div>
                               <div>OS: {node.os}</div>
                               {node.hasVulns && (
                                 <div className="flex items-center gap-1">
                                   <AlertTriangle className="h-3 w-3 text-red-300" />
                                   <span>{node.vulnerabilityCount} vulns</span>
                                   {node.criticalVulns > 0 && (
                                     <span className="text-red-300 font-bold">({node.criticalVulns} critical)</span>
                                   )}
                                 </div>
                               )}
                               {node.location && (
                                 <div className="text-white/60">
                                   {node.location.city}, {node.location.country}
                                 </div>
                               )}
                             </div>
                           </div>
                         ))}
                       </div>
                     ) : (
                       <div className="flex items-center justify-center h-full">
                         <div className="text-center">
                           <RefreshCw className="h-8 w-8 animate-spin mx-auto mb-4 text-orange-500" />
                           <p className="text-gray-600 font-bold">Loading network topology...</p>
                         </div>
                       </div>
                     )}
                   </div>

            {/* Legend */}
            <div className="absolute bottom-4 left-4 bg-white p-3 border-3 border-black rounded-lg shadow-neo-brutal-small">
              <h3 className="font-black uppercase mb-2">Legend</h3>
              <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-sm">
                <div className="flex items-center gap-2"><div className="w-3 h-3 rounded-full bg-blue-500 border-2 border-black"></div>Router</div>
                <div className="flex items-center gap-2"><div className="w-3 h-3 rounded-full bg-yellow-500 border-2 border-black"></div>Switch</div>
                <div className="flex items-center gap-2"><div className="w-3 h-3 rounded-full bg-purple-500 border-2 border-black"></div>Server</div>
                <div className="flex items-center gap-2"><div className="w-3 h-3 rounded-full bg-green-500 border-2 border-black"></div>Workstation</div>
                <div className="flex items-center gap-2"><div className="w-3 h-3 rounded-full bg-orange-500 border-2 border-black"></div>Agent</div>
                <div className="flex items-center gap-2"><div className="w-3 h-3 rounded-full bg-red-500 border-2 border-black"></div>Vulnerable</div>
                <div className="flex items-center gap-2"><div className="w-3 h-3 rounded-full bg-gray-400 border-2 border-black"></div>Offline</div>
              </div>
            </div>

            {isLoading && (
              <div className="absolute inset-0 bg-white bg-opacity-80 flex items-center justify-center">
                <p className="text-black font-bold text-lg">Loading Topology...</p>
              </div>
            )}
          </div>
        </main>

      </div>
    </div>
  );
};

export default Topology;
