import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { 
  Server, 
  Cpu, 
  HardDrive, 
  Wifi, 
  WifiOff,
  Eye,
  Settings,
  RefreshCw,
  Clock,
  X,
  AlertTriangle
} from 'lucide-react';
import { agentService } from '../services/agentService';

// Define types locally to avoid import issues
interface Agent {
  id: string;
  name: string;
  hostname: string;
  os: string;
  status: 'online' | 'offline' | 'unknown';
  last_seen: string;
  cpu_usage: number;
  memory_usage: number;
  organization_id: string;
  version: string;
  ip_address?: string;
  metadata?: Record<string, any>;
}

interface AgentStats {
  total: number;
  online: number;
  offline: number;
  avgCpu: number;
  avgMemory: number;
}

// Real data for agents
const useAgentsData = () => {
  const [data, setData] = useState({
    agents: [] as Agent[],
    stats: {
      total: 0,
      online: 0,
      offline: 0,
      avgCpu: 0,
      avgMemory: 0
    } as AgentStats,
    isLoading: true,
    error: null as string | null
  });

  const loadData = async () => {
    try {
      setData(prev => ({ ...prev, isLoading: true, error: null }));
      
      // Fetch agents and stats in parallel
      const [agents, stats] = await Promise.all([
        agentService.getAgents(),
        agentService.getAgentStats()
      ]);
      
      setData({
        agents,
        stats,
        isLoading: false,
        error: null
      });
    } catch (error) {
      console.error('Failed to load agent data:', error);
      setData(prev => ({ 
        ...prev, 
        isLoading: false, 
        error: 'Failed to load agent data. Please try again.' 
      }));
    }
  };

  useEffect(() => {
    loadData();
    
    // Set up auto-refresh every 30 seconds
    const interval = setInterval(loadData, 30000);
    return () => clearInterval(interval);
  }, []);

  return { ...data, refresh: loadData };
};

const Agents: React.FC = () => {
  const { agents, stats, isLoading, error, refresh } = useAgentsData();
  const [filter, setFilter] = useState<'all' | 'online' | 'offline'>('all');

  const filteredAgents = agents.filter(agent => {
    if (filter === 'all') return true;
    return agent.status === filter;
  });

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online':
        return <Wifi className="h-4 w-4 text-green-600" />;
      case 'offline':
        return <WifiOff className="h-4 w-4 text-red-600" />;
      default:
        return <AlertTriangle className="h-4 w-4 text-yellow-600" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online':
        return 'bg-green-100 text-green-800 border-green-300';
      case 'offline':
        return 'bg-red-100 text-red-800 border-red-300';
      default:
        return 'bg-yellow-100 text-yellow-800 border-yellow-300';
    }
  };

  const formatLastSeen = (lastSeen: string) => {
    const date = new Date(lastSeen);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    
    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffMins < 1440) return `${Math.floor(diffMins / 60)}h ago`;
    return `${Math.floor(diffMins / 1440)}d ago`;
  };

  const restartAgent = async (agentId: string) => {
    try {
      await agentService.restartAgent(agentId);
      toast.success('Agent restart initiated');
      refresh();
    } catch (error) {
      console.error('Failed to restart agent:', error);
      toast.error('Failed to restart agent. Please try again.');
    }
  };

  const killAgent = async (agentId: string) => {
    if (!window.confirm('Are you sure you want to kill this agent? This action cannot be undone.')) {
      return;
    }
    
    try {
      await agentService.killAgent(agentId);
      toast.success('Agent kill initiated');
      refresh();
    } catch (error) {
      console.error('Failed to kill agent:', error);
      toast.error('Failed to kill agent. Please try again.');
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="animate-pulse space-y-6">
            <div className="h-8 bg-gray-200 rounded border-3 border-black"></div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              {Array.from({ length: 4 }).map((_, i) => (
                <div key={i} className="h-32 bg-gray-200 rounded border-3 border-black"></div>
              ))}
            </div>
            <div className="h-96 bg-gray-200 rounded border-3 border-black"></div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="p-6 bg-red-50 border-3 border-red-300 rounded shadow-neubrutalist-lg">
            <div className="flex items-center gap-3 mb-4">
              <AlertTriangle className="h-6 w-6 text-red-600" />
              <h2 className="text-xl font-bold text-red-800">Error Loading Agents</h2>
            </div>
            <p className="text-red-700 mb-4">{error}</p>
            <button 
              onClick={refresh}
              className="px-4 py-2 bg-red-100 text-red-800 border-2 border-red-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-red-200 transition-colors"
            >
              <RefreshCw className="h-4 w-4 mr-2 inline-block" />
              Retry
            </button>
          </div>
        </div>
      </div>
    );
  }

  // Empty state when no agents
  if (agents.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="text-center py-20">
            <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg max-w-2xl mx-auto">
              <div className="mb-6">
                <Server className="h-16 w-16 text-gray-400 mx-auto mb-4" />
                <h1 className="text-3xl font-bold text-black mb-2">No Agents Connected</h1>
                <p className="text-gray-600 text-lg">
                  Deploy the ZeroTrace agent on your systems to start monitoring
                </p>
              </div>
              
              <div className="space-y-4">
                <button 
                  onClick={refresh}
                  className="px-8 py-4 bg-orange-500 text-white border-3 border-black rounded shadow-neubrutalist-lg hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-neubrutalist-md transition-all duration-150 ease-in-out font-bold uppercase tracking-wider"
                >
                  <RefreshCw className="h-5 w-5 mr-2 inline-block" />
                  REFRESH
                </button>
                
                <div className="text-sm text-gray-500">
                  <p>Download the agent for your platform:</p>
                  <div className="flex justify-center gap-4 mt-2">
                    <span className="px-3 py-1 bg-gray-100 border-2 border-gray-300 rounded text-xs font-bold">macOS</span>
                    <span className="px-3 py-1 bg-gray-100 border-2 border-gray-300 rounded text-xs font-bold">Linux</span>
                    <span className="px-3 py-1 bg-gray-100 border-2 border-gray-300 rounded text-xs font-bold">Windows</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b-3 border-black shadow-lg">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold uppercase tracking-wider text-black">Agent Management</h1>
              <p className="text-gray-600 mt-1">Monitor and manage ZeroTrace agents</p>
            </div>
            <div className="flex items-center gap-4">
              <button 
                onClick={refresh}
                className="px-4 py-2 bg-gray-100 text-gray-800 border-2 border-gray-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-gray-200 transition-colors"
              >
                <RefreshCw className="h-4 w-4 mr-2 inline-block" />
                REFRESH
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto p-6 space-y-6">
        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {/* Total Agents */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg hover:shadow-neubrutalist-xl transition-all duration-150">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-blue-100 rounded border-2 border-black">
                <Server className="h-6 w-6 text-blue-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-black">{stats.total}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Total Agents</div>
              </div>
            </div>
          </div>

          {/* Online Agents */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg hover:shadow-neubrutalist-xl transition-all duration-150">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-green-100 rounded border-2 border-black">
                <Wifi className="h-6 w-6 text-green-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-green-600">{stats.online}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Online</div>
              </div>
            </div>
          </div>

          {/* Offline Agents */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg hover:shadow-neubrutalist-xl transition-all duration-150">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-red-100 rounded border-2 border-black">
                <WifiOff className="h-6 w-6 text-red-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-red-600">{stats.offline}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Offline</div>
              </div>
            </div>
          </div>

          {/* Avg CPU Usage */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg hover:shadow-neubrutalist-xl transition-all duration-150">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-orange-100 rounded border-2 border-black">
                <Cpu className="h-6 w-6 text-orange-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-orange-600">{stats.avgCpu.toFixed(1)}%</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Avg CPU</div>
              </div>
            </div>
          </div>
        </div>

        {/* Filter Buttons */}
        <div className="flex gap-4">
          <button
            onClick={() => setFilter('all')}
            className={`px-4 py-2 border-3 border-black rounded text-sm font-bold uppercase tracking-wider transition-all duration-150 ${
              filter === 'all' 
                ? 'bg-orange-100 text-orange-800 border-orange-300 shadow-neubrutalist-md' 
                : 'bg-white text-black hover:bg-gray-50 hover:shadow-neubrutalist-sm'
            }`}
          >
            ALL ({stats.total})
          </button>
          <button
            onClick={() => setFilter('online')}
            className={`px-4 py-2 border-3 border-black rounded text-sm font-bold uppercase tracking-wider transition-all duration-150 ${
              filter === 'online' 
                ? 'bg-orange-100 text-orange-800 border-orange-300 shadow-neubrutalist-md' 
                : 'bg-white text-black hover:bg-gray-50 hover:shadow-neubrutalist-sm'
            }`}
          >
            ONLINE ({stats.online})
          </button>
          <button
            onClick={() => setFilter('offline')}
            className={`px-4 py-2 border-3 border-black rounded text-sm font-bold uppercase tracking-wider transition-all duration-150 ${
              filter === 'offline' 
                ? 'bg-orange-100 text-orange-800 border-orange-300 shadow-neubrutalist-md' 
                : 'bg-white text-black hover:bg-gray-50 hover:shadow-neubrutalist-sm'
            }`}
          >
            OFFLINE ({stats.offline})
          </button>
        </div>

        {/* Agents Table */}
        <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-bold uppercase tracking-wider text-black">Agent Details</h2>
            <div className="text-gray-600 text-sm">
              Showing {filteredAgents.length} of {agents.length} agents
            </div>
          </div>
          
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b-2 border-black">
                  <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">Agent Name</th>
                  <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">Hostname</th>
                  <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">OS</th>
                  <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">Status</th>
                  <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">CPU Usage</th>
                  <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">Memory Usage</th>
                  <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">Last Seen</th>
                  <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredAgents.map((agent) => (
                  <tr key={agent.id} className="border-b border-gray-200 hover:bg-gray-50">
                    <td className="py-3 px-4 font-bold text-black">{agent.name}</td>
                    <td className="py-3 px-4 text-gray-600">{agent.hostname}</td>
                    <td className="py-3 px-4 text-gray-600">{agent.os}</td>
                    <td className="py-3 px-4">
                      <div className="flex items-center gap-2">
                        {getStatusIcon(agent.status)}
                        <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider border-2 ${getStatusColor(agent.status)}`}>
                          {agent.status}
                        </span>
                      </div>
                    </td>
                    <td className="py-3 px-4">
                      <div className="flex items-center gap-2">
                        <Cpu className="h-4 w-4 text-blue-600" />
                        <span className="font-medium">{agent.cpu_usage}%</span>
                      </div>
                    </td>
                    <td className="py-3 px-4">
                      <div className="flex items-center gap-2">
                        <HardDrive className="h-4 w-4 text-green-600" />
                        <span className="font-medium">{agent.memory_usage}%</span>
                      </div>
                    </td>
                    <td className="py-3 px-4">
                      <div className="flex items-center gap-2">
                        <Clock className="h-4 w-4 text-gray-600" />
                        <span className="text-sm text-gray-600">{formatLastSeen(agent.last_seen)}</span>
                      </div>
                    </td>
                    <td className="py-3 px-4">
                      <div className="flex items-center gap-2">
                        <button 
                          onClick={() => restartAgent(agent.id)}
                          className="p-2 hover:bg-blue-100 rounded border-2 border-black text-blue-600"
                          title="Restart Agent"
                        >
                          <RefreshCw className="h-4 w-4" />
                        </button>
                        <button 
                          onClick={() => killAgent(agent.id)}
                          className="p-2 hover:bg-red-100 rounded border-2 border-black text-red-600"
                          title="Kill Agent"
                        >
                          <X className="h-4 w-4" />
                        </button>
                        <Link 
                          to={`/agents/${agent.id}`}
                          className="p-2 hover:bg-gray-200 rounded border-2 border-black"
                          title="View Asset Details"
                        >
                          <Eye className="h-4 w-4" />
                        </Link>
                        <button className="p-2 hover:bg-gray-200 rounded border-2 border-black">
                          <Settings className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Agents;