import React, { useState, useEffect } from 'react';
import { RefreshCw, Search, SlidersHorizontal, Share2, ZoomIn, ZoomOut, Maximize, Shield, Activity, AlertTriangle } from 'lucide-react';
import { topologyService } from '../services/topologyService';
import NetworkTopology from '../components/NetworkTopology';
// import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/card';

interface NetworkTopologyData {
  nodes: any[];
  links: any[];
  id: string;
  companyId: string;
  clusters: any[];
  lastUpdated: string;
}

interface NetworkHealth {
  onlineNodes: number;
  vulnerableNodes: number;
  averageRiskScore: number;
}

const Topology: React.FC = () => {
  const [data, setData] = useState<NetworkTopologyData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [health, setHealth] = useState<NetworkHealth | null>(null);
  // const topologyRef = useRef<any>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const topologyData = await topologyService.getTopologyData();
        setData({ ...topologyData, id: '1', companyId: '1', clusters: [], lastUpdated: new Date().toISOString() });
        
        // Mock health data
        setHealth({
          onlineNodes: topologyData.nodes.filter((n: any) => n.status === 'online').length,
          vulnerableNodes: topologyData.nodes.filter((n: any) => n.riskScore > 5).length,
          averageRiskScore: topologyData.nodes.reduce((acc: any, n: any) => acc + n.riskScore, 0) / topologyData.nodes.length || 0
        });

      } catch (err) {
        setError('Failed to load topology data. Please try again later.');
        console.error('Error loading topology:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const handleRefresh = () => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const topologyData = await topologyService.getTopologyData();
        setData({ ...topologyData, id: '1', companyId: '1', clusters: [], lastUpdated: new Date().toISOString() });
        
        // Mock health data
        setHealth({
          onlineNodes: topologyData.nodes.filter((n: any) => n.status === 'online').length,
          vulnerableNodes: topologyData.nodes.filter((n: any) => n.riskScore > 5).length,
          averageRiskScore: topologyData.nodes.reduce((acc: any, n: any) => acc + n.riskScore, 0) / topologyData.nodes.length || 0
        });

      } catch (err) {
        setError('Failed to load topology data. Please try again later.');
        console.error('Error loading topology:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  };

  const handleNodeClick = () => {
    // Handle node click
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
          {health && (
            <div className="flex gap-4">
              <div className="px-4 py-2 bg-white border-3 border-black rounded-lg shadow-neo-brutal-small">
                <div className="flex items-center gap-2">
                  <Activity className="h-4 w-4 text-green-600" />
                  <span className="font-bold text-green-600">{health.onlineNodes}</span>
                  <span className="text-sm text-gray-600">Online</span>
                </div>
              </div>
              <div className="px-4 py-2 bg-white border-3 border-black rounded-lg shadow-neo-brutal-small">
                <div className="flex items-center gap-2">
                  <AlertTriangle className="h-4 w-4 text-red-600" />
                  <span className="font-bold text-red-600">{health.vulnerableNodes}</span>
                  <span className="text-sm text-gray-600">Vulnerable</span>
                </div>
              </div>
              <div className="px-4 py-2 bg-white border-3 border-black rounded-lg shadow-neo-brutal-small">
                <div className="flex items-center gap-2">
                  <Shield className="h-4 w-4 text-blue-600" />
                  <span className="font-bold text-blue-600">{health.averageRiskScore.toFixed(1)}</span>
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
                <RefreshCw className={`h-5 w-5 ${loading ? 'animate-spin' : ''}`} />
                {loading ? 'Refreshing...' : 'Refresh'}
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
                     ) : data ? (
                  <NetworkTopology
                    data={data}
                    onNodeClick={handleNodeClick}
                  />
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

            {loading && (
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
