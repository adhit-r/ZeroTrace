import React, { useState, useEffect } from 'react';
import { 
  Server, 
  Activity, 
  Cpu, 
  HardDrive, 
  Wifi, 
  WifiOff,
  Eye,
  Settings,
  RefreshCw,
  Terminal
} from 'lucide-react';

// Mock data for agents
const useAgentsData = () => {
  const [data, setData] = useState({
    agents: [],
    stats: {
      total: 0,
      online: 0,
      offline: 0,
      avgCpu: 0,
      avgMemory: 0
    },
    isLoading: true
  });

  useEffect(() => {
    const loadData = async () => {
      // Simulate API call delay
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setData({
        agents: [
          {
            id: '1',
            name: 'AGENT-001',
            hostname: 'web-server-01',
            os: 'Ubuntu 22.04',
            status: 'online',
            lastSeen: new Date(),
            cpuUsage: 45.2,
            memoryUsage: 67.8,
            organization: 'ACME Corp',
            version: '1.0.0'
          },
          {
            id: '2',
            name: 'AGENT-002',
            hostname: 'db-server-01',
            os: 'CentOS 8',
            status: 'online',
            lastSeen: new Date(Date.now() - 5 * 60 * 1000), // 5 minutes ago
            cpuUsage: 23.1,
            memoryUsage: 89.2,
            organization: 'ACME Corp',
            version: '1.0.0'
          },
          {
            id: '3',
            name: 'AGENT-003',
            hostname: 'file-server-01',
            os: 'Windows Server 2019',
            status: 'offline',
            lastSeen: new Date(Date.now() - 30 * 60 * 1000), // 30 minutes ago
            cpuUsage: 0,
            memoryUsage: 0,
            organization: 'ACME Corp',
            version: '1.0.0'
          },
          {
            id: '4',
            name: 'AGENT-004',
            hostname: 'app-server-01',
            os: 'macOS 13.0',
            status: 'online',
            lastSeen: new Date(),
            cpuUsage: 12.5,
            memoryUsage: 34.7,
            organization: 'TechCorp',
            version: '1.0.0'
          }
        ],
        stats: {
          total: 4,
          online: 3,
          offline: 1,
          avgCpu: 20.2,
          avgMemory: 47.9
        },
        isLoading: false
      });
    };

    loadData();
  }, []);

  return data;
};

const Agents: React.FC = () => {
  const { agents, stats, isLoading } = useAgentsData();
  const [filter, setFilter] = useState('all');

  const filteredAgents = agents.filter(agent => {
    if (filter === 'all') return true;
    if (filter === 'online') return agent.status === 'online';
    if (filter === 'offline') return agent.status === 'offline';
    return true;
  });

  if (isLoading) {
    return (
      <div className="flex-center h-64">
        <div className="loading"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gold text-glow">AGENT MONITORING</h1>
          <p className="text-text-secondary mt-1">ACTIVE AGENT STATUS AND PERFORMANCE</p>
        </div>
        <div className="flex items-center space-x-3">
          <button className="btn btn-secondary btn-sm">
            <RefreshCw className="h-4 w-4 mr-2" />
            REFRESH
          </button>
          <button className="btn btn-primary btn-sm">
            <Terminal className="h-4 w-4 mr-2" />
            DEPLOY AGENT
          </button>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        {/* Total Agents */}
        <div className="card card-terminal glow-border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-text-secondary text-sm uppercase tracking-wide">TOTAL AGENTS</p>
              <p className="text-3xl font-bold text-text-primary">{stats.total}</p>
            </div>
            <div className="h-12 w-12 bg-medium-gray rounded-lg flex items-center justify-center border border-gold">
              <Server className="h-6 w-6 text-gold" />
            </div>
          </div>
        </div>

        {/* Online Agents */}
        <div className="card card-terminal glow-border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-text-secondary text-sm uppercase tracking-wide">ONLINE</p>
              <p className="text-3xl font-bold status-online">{stats.online}</p>
            </div>
            <div className="h-12 w-12 bg-medium-gray rounded-lg flex items-center justify-center border border-success">
              <Wifi className="h-6 w-6 text-success" />
            </div>
          </div>
        </div>

        {/* Offline Agents */}
        <div className="card card-terminal glow-border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-text-secondary text-sm uppercase tracking-wide">OFFLINE</p>
              <p className="text-3xl font-bold status-offline">{stats.offline}</p>
            </div>
            <div className="h-12 w-12 bg-medium-gray rounded-lg flex items-center justify-center border border-text-muted">
              <WifiOff className="h-6 w-6 text-text-muted" />
            </div>
          </div>
        </div>

        {/* Average CPU */}
        <div className="card card-terminal glow-border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-text-secondary text-sm uppercase tracking-wide">AVG CPU</p>
              <p className="text-3xl font-bold text-warning">{stats.avgCpu.toFixed(1)}%</p>
            </div>
            <div className="h-12 w-12 bg-medium-gray rounded-lg flex items-center justify-center border border-warning">
              <Cpu className="h-6 w-6 text-warning" />
            </div>
          </div>
        </div>

        {/* Average Memory */}
        <div className="card card-terminal glow-border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-text-secondary text-sm uppercase tracking-wide">AVG MEMORY</p>
              <p className="text-3xl font-bold text-info">{stats.avgMemory.toFixed(1)}%</p>
            </div>
            <div className="h-12 w-12 bg-medium-gray rounded-lg flex items-center justify-center border border-info">
              <HardDrive className="h-6 w-6 text-info" />
            </div>
          </div>
        </div>
      </div>

      {/* Filters */}
      <div className="flex items-center space-x-4">
        <button
          onClick={() => setFilter('all')}
          className={`btn btn-sm ${filter === 'all' ? 'btn-primary' : 'btn-secondary'}`}
        >
          ALL ({stats.total})
        </button>
        <button
          onClick={() => setFilter('online')}
          className={`btn btn-sm ${filter === 'online' ? 'btn-primary' : 'btn-secondary'}`}
        >
          ONLINE ({stats.online})
        </button>
        <button
          onClick={() => setFilter('offline')}
          className={`btn btn-sm ${filter === 'offline' ? 'btn-primary' : 'btn-secondary'}`}
        >
          OFFLINE ({stats.offline})
        </button>
      </div>

      {/* Agents Table */}
      <div className="card card-terminal">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-text-primary">AGENT DETAILS</h2>
          <div className="text-text-secondary text-sm">
            SHOWING {filteredAgents.length} OF {agents.length} AGENTS
          </div>
        </div>
        
        <div className="overflow-x-auto">
          <table className="table">
            <thead>
              <tr>
                <th>AGENT NAME</th>
                <th>HOSTNAME</th>
                <th>OS</th>
                <th>STATUS</th>
                <th>CPU USAGE</th>
                <th>MEMORY USAGE</th>
                <th>LAST SEEN</th>
                <th>ORGANIZATION</th>
                <th>ACTIONS</th>
              </tr>
            </thead>
            <tbody>
              {filteredAgents.map((agent) => (
                <tr key={agent.id}>
                  <td className="text-text-primary font-medium">{agent.name}</td>
                  <td className="text-text-secondary">{agent.hostname}</td>
                  <td className="text-text-secondary">{agent.os}</td>
                  <td>
                    <span className={`badge ${agent.status === 'online' ? 'badge-low' : 'badge-critical'}`}>
                      {agent.status.toUpperCase()}
                    </span>
                  </td>
                  <td>
                    <div className="flex items-center space-x-2">
                      <div className="progress-bar w-16">
                        <div 
                          className={`progress-fill ${agent.cpuUsage > 80 ? 'progress-danger' : agent.cpuUsage > 60 ? 'progress-warning' : 'progress-success'}`}
                          style={{ width: `${agent.cpuUsage}%` }}
                        ></div>
                      </div>
                      <span className="text-text-secondary text-sm">{agent.cpuUsage.toFixed(1)}%</span>
                    </div>
                  </td>
                  <td>
                    <div className="flex items-center space-x-2">
                      <div className="progress-bar w-16">
                        <div 
                          className={`progress-fill ${agent.memoryUsage > 80 ? 'progress-danger' : agent.memoryUsage > 60 ? 'progress-warning' : 'progress-success'}`}
                          style={{ width: `${agent.memoryUsage}%` }}
                        ></div>
                      </div>
                      <span className="text-text-secondary text-sm">{agent.memoryUsage.toFixed(1)}%</span>
                    </div>
                  </td>
                  <td className="text-text-secondary text-sm">
                    {agent.lastSeen.toLocaleTimeString()}
                  </td>
                  <td className="text-text-secondary">{agent.organization}</td>
                  <td>
                    <div className="flex items-center space-x-2">
                      <button className="btn btn-ghost btn-sm">
                        <Eye className="h-4 w-4" />
                      </button>
                      <button className="btn btn-ghost btn-sm">
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
  );
};

export default Agents;

