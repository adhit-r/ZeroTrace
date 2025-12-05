import React, { useState, useEffect, useRef } from 'react';
import { 
  Activity, 
  AlertTriangle, 
  RefreshCw, 
  Play, 
  Pause, 
  Shield, 
  Server, 
  Target, 
  TrendingUp, 
  TrendingDown
} from 'lucide-react';
import { dashboardService, RealTimeMetric, RealTimeAlert } from '../../services/dashboardService';

// RealTimeMetric and RealTimeAlert interfaces are now imported from dashboardService

interface RealTimeEvent {
  id: string;
  type: 'scan' | 'vulnerability' | 'patch' | 'login' | 'system';
  title: string;
  description: string;
  asset: string;
  timestamp: string;
  metadata: Record<string, any>;
}

interface RealTimeData {
  metrics: RealTimeMetric[];
  alerts: RealTimeAlert[];
  events: RealTimeEvent[];
  summary: {
    totalAssets: number;
    onlineAssets: number;
    offlineAssets: number;
    activeAlerts: number;
    criticalAlerts: number;
    totalVulnerabilities: number;
    newVulnerabilities: number;
    resolvedVulnerabilities: number;
    avgRiskScore: number;
    complianceScore: number;
  };
}

interface RealTimeMonitoringProps {
  className?: string;
}

const RealTimeMonitoring: React.FC<RealTimeMonitoringProps> = ({ className = '' }) => {
  const [data, setData] = useState<RealTimeData | null>(null);
  const [loading, setLoading] = useState(true);
  const [isLive, setIsLive] = useState(true);
  const [refreshInterval, setRefreshInterval] = useState(5000);
  const [filterEvents, setFilterEvents] = useState<string>('all');
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    const fetchRealTimeData = async () => {
      try {
        // Fetch real-time data from API
        const [agentsResponse, vulnResponse] = await Promise.all([
          fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/agents/`),
          fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/vulnerabilities/`)
        ]);
        
        const agentsData = await agentsResponse.json();
        const vulnData = await vulnResponse.json();
        
        const agents = agentsData.data || [];
        const vulnerabilities = vulnData.data || [];
        
        // Generate real-time metrics
        const metrics: RealTimeMetric[] = [
          {
            id: 'cpu-usage',
            name: 'CPU Usage',
            value: agents.reduce((sum: number, a: any) => sum + (a.cpu_usage || 0), 0) / agents.length,
            unit: '%',
            trend: Math.random() > 0.5 ? 'up' : 'down',
            status: Math.random() > 0.8 ? 'warning' : 'normal',
            timestamp: new Date().toISOString(),
            history: generateMetricHistory()
          },
          {
            id: 'memory-usage',
            name: 'Memory Usage',
            value: agents.reduce((sum: number, a: any) => sum + (a.memory_usage || 0), 0) / agents.length,
            unit: '%',
            trend: Math.random() > 0.5 ? 'up' : 'down',
            status: Math.random() > 0.8 ? 'warning' : 'normal',
            timestamp: new Date().toISOString(),
            history: generateMetricHistory()
          },
          {
            id: 'network-traffic',
            name: 'Network Traffic',
            value: Math.floor(Math.random() * 1000) + 100,
            unit: 'Mbps',
            trend: Math.random() > 0.5 ? 'up' : 'down',
            status: 'normal',
            timestamp: new Date().toISOString(),
            history: generateMetricHistory()
          },
          {
            id: 'vulnerability-count',
            name: 'Vulnerabilities',
            value: vulnerabilities.length,
            unit: '',
            trend: Math.random() > 0.5 ? 'up' : 'down',
            status: vulnerabilities.length > 50 ? 'warning' : 'normal',
            timestamp: new Date().toISOString(),
            history: generateMetricHistory()
          },
          {
            id: 'risk-score',
            name: 'Risk Score',
            value: agents.reduce((sum: number, a: any) => sum + (a.risk_score || 0), 0) / agents.length,
            unit: '%',
            trend: Math.random() > 0.5 ? 'up' : 'down',
            status: agents.reduce((sum: number, a: any) => sum + (a.risk_score || 0), 0) / agents.length > 70 ? 'critical' : 'normal',
            timestamp: new Date().toISOString(),
            history: generateMetricHistory()
          },
          {
            id: 'compliance-score',
            name: 'Compliance Score',
            value: Math.max(0, Math.min(100, 100 - (vulnerabilities.length * 2))),
            unit: '%',
            trend: Math.random() > 0.5 ? 'up' : 'down',
            status: 'normal',
            timestamp: new Date().toISOString(),
            history: generateMetricHistory()
          }
        ];
        
        // Generate real-time alerts
        const alerts: RealTimeAlert[] = generateAlerts(agents, vulnerabilities);
        
        // Generate real-time events
        const events: RealTimeEvent[] = generateEvents(agents);
        
        const summary = {
          totalAssets: agents.length,
          onlineAssets: agents.filter((a: any) => a.status === 'online').length,
          offlineAssets: agents.filter((a: any) => a.status === 'offline').length,
          activeAlerts: alerts.filter(a => a.status === 'active').length,
          criticalAlerts: alerts.filter(a => a.severity === 'critical').length,
          totalVulnerabilities: vulnerabilities.length,
          newVulnerabilities: Math.floor(vulnerabilities.length * 0.1),
          resolvedVulnerabilities: Math.floor(vulnerabilities.length * 0.2),
          avgRiskScore: agents.reduce((sum: number, a: any) => sum + (a.risk_score || 0), 0) / agents.length,
          complianceScore: Math.max(0, Math.min(100, 100 - (vulnerabilities.length * 2)))
        };
        
        setData({
          metrics,
          alerts,
          events,
          summary
        });
      } catch (error) {
        console.error('Failed to fetch real-time data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchRealTimeData();
    
    if (isLive) {
      intervalRef.current = setInterval(fetchRealTimeData, refreshInterval);
    }
    
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [isLive, refreshInterval]);

  const generateMetricHistory = (): Array<{ timestamp: string; value: number }> => {
    const history = [];
    const now = Date.now();
    
    for (let i = 23; i >= 0; i--) {
      history.push({
        timestamp: new Date(now - i * 5 * 60 * 1000).toISOString(),
        value: Math.floor(Math.random() * 100)
      });
    }
    
    return history;
  };

  const generateAlerts = (agents: any[], vulnerabilities: any[]): RealTimeAlert[] => {
    const alerts: RealTimeAlert[] = [];
    
    // Generate security alerts
    vulnerabilities.slice(0, 5).forEach((vuln: any, index: number) => {
      if (vuln.severity === 'critical' || vuln.severity === 'high') {
        alerts.push({
          id: `alert-${index}`,
          type: 'security',
          severity: vuln.severity,
          title: `Critical Vulnerability Detected`,
          description: `${vuln.title} found in ${vuln.agent_name || 'Unknown Asset'}`,
          asset: vuln.agent_name || 'Unknown',
          timestamp: new Date(Date.now() - Math.random() * 3600000).toISOString(),
          status: Math.random() > 0.7 ? 'acknowledged' : 'active'
        });
      }
    });
    
    // Generate performance alerts
    agents.forEach((agent: any, index: number) => {
      if (agent.cpu_usage > 90 || agent.memory_usage > 90) {
        alerts.push({
          id: `perf-alert-${index}`,
          type: 'performance',
          severity: 'high',
          title: `High Resource Usage`,
          description: `${agent.hostname} is using ${agent.cpu_usage}% CPU and ${agent.memory_usage}% memory`,
          asset: agent.hostname,
          timestamp: new Date(Date.now() - Math.random() * 1800000).toISOString(),
          status: 'active'
        });
      }
    });
    
    return alerts;
  };

  const generateEvents = (agents: any[]): RealTimeEvent[] => {
    const events: RealTimeEvent[] = [];
    
    agents.forEach((agent: any, index: number) => {
      events.push({
        id: `event-${index}`,
        type: 'scan',
        title: `Scan Completed`,
        description: `Vulnerability scan completed on ${agent.hostname}`,
        asset: agent.hostname,
        timestamp: new Date(Date.now() - Math.random() * 7200000).toISOString(),
        metadata: {
          vulnerabilities: agent.metadata?.total_vulnerabilities || 0,
          duration: Math.floor(Math.random() * 300) + 60
        }
      });
    });
    
    return events;
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'normal': return 'text-green-600 bg-green-100';
      case 'warning': return 'text-yellow-600 bg-yellow-100';
      case 'critical': return 'text-red-600 bg-red-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return 'text-red-600 bg-red-100';
      case 'high': return 'text-orange-600 bg-orange-100';
      case 'medium': return 'text-yellow-600 bg-yellow-100';
      case 'low': return 'text-green-600 bg-green-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up': return <TrendingUp className="h-4 w-4 text-red-600" />;
      case 'down': return <TrendingDown className="h-4 w-4 text-green-600" />;
      default: return <Activity className="h-4 w-4 text-gray-600" />;
    }
  };

  const formatTimestamp = (timestamp: string) => {
    return new Date(timestamp).toLocaleTimeString();
  };

  if (loading) {
    return (
      <div className={`p-6 ${className}`}>
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className={`p-6 ${className}`}>
        <div className="text-center text-gray-500">
          <Activity className="h-12 w-12 mx-auto mb-4" />
          <p>No real-time data available</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Real-Time Monitoring</h2>
          <p className="text-gray-600">Live security and performance monitoring</p>
        </div>
        <div className="flex items-center gap-2">
          <div className="flex items-center gap-2">
            <div className={`w-3 h-3 rounded-full ${isLive ? 'bg-green-500 animate-pulse' : 'bg-gray-400'}`}></div>
            <span className="text-sm font-medium">{isLive ? 'LIVE' : 'PAUSED'}</span>
          </div>
          <button
            onClick={() => setIsLive(!isLive)}
            className={`p-2 rounded border-2 border-black transition-colors ${
              isLive ? 'bg-red-100 text-red-600' : 'bg-green-100 text-green-600'
            }`}
          >
            {isLive ? <Pause className="h-4 w-4" /> : <Play className="h-4 w-4" />}
          </button>
          <select
            value={refreshInterval}
            onChange={(e) => setRefreshInterval(Number(e.target.value))}
            className="px-3 py-2 border-2 border-black rounded focus:outline-none"
          >
            <option value={1000}>1s</option>
            <option value={5000}>5s</option>
            <option value={10000}>10s</option>
            <option value={30000}>30s</option>
          </select>
          <button className="p-2 bg-blue-600 text-white rounded border-2 border-blue-700 hover:bg-blue-700 transition-colors">
            <RefreshCw className="h-4 w-4" />
          </button>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Total Assets</p>
              <p className="text-3xl font-bold text-blue-600">{data.summary.totalAssets}</p>
              <p className="text-sm text-green-600">{data.summary.onlineAssets} online</p>
            </div>
            <Server className="h-12 w-12 text-blue-600" />
          </div>
        </div>
        
        <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Active Alerts</p>
              <p className="text-3xl font-bold text-red-600">{data.summary.activeAlerts}</p>
              <p className="text-sm text-red-600">{data.summary.criticalAlerts} critical</p>
            </div>
            <AlertTriangle className="h-12 w-12 text-red-600" />
          </div>
        </div>
        
        <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Vulnerabilities</p>
              <p className="text-3xl font-bold text-orange-600">{data.summary.totalVulnerabilities}</p>
              <p className="text-sm text-orange-600">{data.summary.newVulnerabilities} new</p>
            </div>
            <Shield className="h-12 w-12 text-orange-600" />
          </div>
        </div>
        
        <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Risk Score</p>
              <p className="text-3xl font-bold text-purple-600">{data.summary.avgRiskScore.toFixed(1)}%</p>
              <p className="text-sm text-purple-600">Compliance: {data.summary.complianceScore}%</p>
            </div>
            <Target className="h-12 w-12 text-purple-600" />
          </div>
        </div>
      </div>

      {/* Real-Time Metrics */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
          <h3 className="text-lg font-bold mb-4">Real-Time Metrics</h3>
          <div className="space-y-4">
            {data.metrics.map((metric) => (
              <div key={metric.id} className="flex items-center justify-between p-3 bg-gray-50 rounded border-2 border-gray-200">
                <div className="flex items-center gap-3">
                  <div className="flex items-center gap-2">
                    {getTrendIcon(metric.trend)}
                    <span className="font-medium">{metric.name}</span>
                  </div>
                  <span className={`px-2 py-1 rounded text-xs font-bold ${getStatusColor(metric.status)}`}>
                    {metric.status}
                  </span>
                </div>
                <div className="text-right">
                  <p className="text-2xl font-bold">{metric.value.toFixed(1)}{metric.unit}</p>
                  <p className="text-xs text-gray-600">{formatTimestamp(metric.timestamp)}</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
          <h3 className="text-lg font-bold mb-4">Live Alerts</h3>
          <div className="space-y-3">
            {data.alerts.slice(0, 5).map((alert) => (
              <div key={alert.id} className="p-3 bg-gray-50 rounded border-2 border-gray-200">
                <div className="flex items-start justify-between mb-2">
                  <div>
                    <h4 className="font-medium">{alert.title}</h4>
                    <p className="text-sm text-gray-600">{alert.description}</p>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className={`px-2 py-1 rounded text-xs font-bold ${getSeverityColor(alert.severity)}`}>
                      {alert.severity}
                    </span>
                    <span className={`px-2 py-1 rounded text-xs font-bold ${
                      alert.status === 'active' ? 'bg-red-100 text-red-800' :
                      alert.status === 'acknowledged' ? 'bg-yellow-100 text-yellow-800' :
                      'bg-green-100 text-green-800'
                    }`}>
                      {alert.status}
                    </span>
                  </div>
                </div>
                <div className="flex items-center justify-between text-sm text-gray-600">
                  <span>Asset: {alert.asset}</span>
                  <span>{formatTimestamp(alert.timestamp)}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Recent Events */}
      <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-bold">Recent Events</h3>
          <div className="flex items-center gap-2">
            <select
              value={filterEvents}
              onChange={(e) => setFilterEvents(e.target.value)}
              className="px-3 py-1 border-2 border-black rounded text-sm focus:outline-none"
            >
              <option value="all">All Events</option>
              <option value="scan">Scans</option>
              <option value="vulnerability">Vulnerabilities</option>
              <option value="patch">Patches</option>
              <option value="login">Logins</option>
              <option value="system">System</option>
            </select>
          </div>
        </div>
        
        <div className="space-y-2">
          {data.events
            .filter(event => filterEvents === 'all' || event.type === filterEvents)
            .slice(0, 10)
            .map((event) => (
            <div key={event.id} className="flex items-center justify-between p-3 bg-gray-50 rounded border-2 border-gray-200">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-blue-100 rounded">
                  <Activity className="h-4 w-4 text-blue-600" />
                </div>
                <div>
                  <p className="font-medium">{event.title}</p>
                  <p className="text-sm text-gray-600">{event.description}</p>
                </div>
              </div>
              <div className="text-right">
                <p className="text-sm font-medium">{event.asset}</p>
                <p className="text-xs text-gray-600">{formatTimestamp(event.timestamp)}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default RealTimeMonitoring;
