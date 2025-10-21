import React, { useState, useEffect } from 'react';
import { 
  BarChart3, 
  Activity, 
  Shield, 
  AlertTriangle, 
  CheckCircle, 
  Server, 
  Network, 
  Layers, 
  Zap,
  RefreshCw,
  Cpu,
  MemoryStick,
  HardDrive,
  Download
} from 'lucide-react';
import { dashboardService, DashboardMetrics } from '../../services/dashboardService';

// DashboardMetrics interface is now imported from dashboardService

interface InnovativeDashboardProps {
  className?: string;
}

const InnovativeDashboard: React.FC<InnovativeDashboardProps> = ({ className = '' }) => {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [timeRange, setTimeRange] = useState<'1h' | '24h' | '7d' | '30d'>('24h');
  const [viewMode, setViewMode] = useState<'overview' | 'security' | 'performance' | 'compliance'>('overview');

  useEffect(() => {
    const fetchMetrics = async () => {
      setLoading(true);
      try {
        const metrics = await dashboardService.getDashboardMetrics();
        setMetrics(metrics);
      } catch (error) {
        console.error('Failed to fetch dashboard metrics:', error);
      } finally {
        setLoading(false);
      }
    };
    
    fetchMetrics();
    
    // Auto-refresh every 30 seconds
    const interval = setInterval(fetchMetrics, 30000);
    return () => clearInterval(interval);
  }, [timeRange]);

  if (loading) {
    return (
      <div className={`p-6 ${className}`}>
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  if (!metrics) {
    return (
      <div className={`p-6 ${className}`}>
        <div className="text-center text-gray-500">
          <BarChart3 className="h-12 w-12 mx-auto mb-4" />
          <p>No metrics data available</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Innovative Dashboard</h2>
          <p className="text-gray-600">Real-time security and performance insights</p>
        </div>
        <div className="flex items-center gap-2">
          <select
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value as any)}
            className="px-4 py-2 border-2 border-black rounded focus:outline-none"
          >
            <option value="1h">Last Hour</option>
            <option value="24h">Last 24 Hours</option>
            <option value="7d">Last 7 Days</option>
            <option value="30d">Last 30 Days</option>
          </select>
          <button className="p-2 bg-blue-600 text-white rounded border-2 border-blue-700 hover:bg-blue-700 transition-colors">
            <RefreshCw className="h-4 w-4" />
          </button>
        </div>
      </div>

      {/* View Mode Tabs */}
      <div className="flex space-x-1 bg-gray-100 p-1 rounded-lg">
        {[
          { id: 'overview', label: 'Overview', icon: BarChart3 },
          { id: 'security', label: 'Security', icon: Shield },
          { id: 'performance', label: 'Performance', icon: Activity },
          { id: 'compliance', label: 'Compliance', icon: CheckCircle }
        ].map((tab) => (
          <button
            key={tab.id}
            onClick={() => setViewMode(tab.id as any)}
            className={`flex items-center gap-2 px-4 py-2 rounded-md transition-colors ${
              viewMode === tab.id
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            <tab.icon className="h-4 w-4" />
            {tab.label}
          </button>
        ))}
      </div>

      {/* Overview Mode */}
      {viewMode === 'overview' && (
        <div className="space-y-6">
          {/* Key Metrics Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg hover:shadow-xl transition-shadow">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Total Assets</p>
                  <p className="text-3xl font-bold text-blue-600">{metrics.assets.total}</p>
                  <p className="text-sm text-green-600">↑ {metrics.assets.online} online</p>
                </div>
                <Server className="h-12 w-12 text-blue-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg hover:shadow-xl transition-shadow">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Vulnerabilities</p>
                  <p className="text-3xl font-bold text-red-600">{metrics.vulnerabilities.total}</p>
                  <p className="text-sm text-red-600">{metrics.vulnerabilities.critical} critical</p>
                </div>
                <AlertTriangle className="h-12 w-12 text-red-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg hover:shadow-xl transition-shadow">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Compliance Score</p>
                  <p className="text-3xl font-bold text-green-600">{metrics.compliance.score}%</p>
                  <p className="text-sm text-green-600">{metrics.compliance.compliant} compliant</p>
                </div>
                <CheckCircle className="h-12 w-12 text-green-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg hover:shadow-xl transition-shadow">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Risk Score</p>
                  <p className="text-3xl font-bold text-orange-600">{metrics.security.riskScore}%</p>
                  <p className="text-sm text-orange-600">{metrics.security.patchesAvailable} patches</p>
                </div>
                <Shield className="h-12 w-12 text-orange-600" />
              </div>
            </div>
          </div>

          {/* Asset Risk Distribution */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <h3 className="text-lg font-bold mb-4">Asset Risk Distribution</h3>
              <div className="space-y-4">
                {[
                  { level: 'Critical', count: metrics.assets.critical, color: 'bg-red-600' },
                  { level: 'High', count: metrics.assets.high, color: 'bg-orange-600' },
                  { level: 'Medium', count: metrics.assets.medium, color: 'bg-yellow-600' },
                  { level: 'Low', count: metrics.assets.low, color: 'bg-green-600' }
                ].map((item) => (
                  <div key={item.level} className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className={`w-4 h-4 rounded ${item.color}`}></div>
                      <span className="font-medium">{item.level}</span>
                    </div>
                    <div className="flex items-center gap-3">
                      <div className="w-32 bg-gray-200 rounded-full h-2">
                        <div 
                          className={`h-2 rounded-full ${item.color}`}
                          style={{ width: `${(item.count / metrics.assets.total) * 100}%` }}
                        ></div>
                      </div>
                      <span className="font-bold w-8 text-right">{item.count}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <h3 className="text-lg font-bold mb-4">Vulnerability Breakdown</h3>
              <div className="space-y-4">
                {[
                  { severity: 'Critical', count: metrics.vulnerabilities.critical, color: 'bg-red-600' },
                  { severity: 'High', count: metrics.vulnerabilities.high, color: 'bg-orange-600' },
                  { severity: 'Medium', count: metrics.vulnerabilities.medium, color: 'bg-yellow-600' },
                  { severity: 'Low', count: metrics.vulnerabilities.low, color: 'bg-green-600' }
                ].map((item) => (
                  <div key={item.severity} className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className={`w-4 h-4 rounded ${item.color}`}></div>
                      <span className="font-medium">{item.severity}</span>
                    </div>
                    <div className="flex items-center gap-3">
                      <div className="w-32 bg-gray-200 rounded-full h-2">
                        <div 
                          className={`h-2 rounded-full ${item.color}`}
                          style={{ width: `${(item.count / metrics.vulnerabilities.total) * 100}%` }}
                        ></div>
                      </div>
                      <span className="font-bold w-8 text-right">{item.count}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Performance Metrics */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
                <Cpu className="h-5 w-5" />
                CPU Performance
              </h3>
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-gray-600">Average Usage</span>
                  <span className="font-bold">{metrics.performance.avgCpu.toFixed(1)}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-3">
                  <div 
                    className="bg-blue-600 h-3 rounded-full"
                    style={{ width: `${metrics.performance.avgCpu}%` }}
                  ></div>
                </div>
              </div>
            </div>

            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
                <MemoryStick className="h-5 w-5" />
                Memory Performance
              </h3>
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-gray-600">Average Usage</span>
                  <span className="font-bold">{metrics.performance.avgMemory.toFixed(1)}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-3">
                  <div 
                    className="bg-green-600 h-3 rounded-full"
                    style={{ width: `${metrics.performance.avgMemory}%` }}
                  ></div>
                </div>
              </div>
            </div>

            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
                <HardDrive className="h-5 w-5" />
                Storage Performance
              </h3>
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-gray-600">Average Usage</span>
                  <span className="font-bold">{metrics.performance.avgStorage}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-3">
                  <div 
                    className="bg-purple-600 h-3 rounded-full"
                    style={{ width: `${metrics.performance.avgStorage}%` }}
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Security Mode */}
      {viewMode === 'security' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Active Alerts</p>
                  <p className="text-3xl font-bold text-red-600">{metrics.security.alertsActive}</p>
                </div>
                <AlertTriangle className="h-12 w-12 text-red-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Patches Available</p>
                  <p className="text-3xl font-bold text-orange-600">{metrics.security.patchesAvailable}</p>
                </div>
                <Download className="h-12 w-12 text-orange-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Scans Completed</p>
                  <p className="text-3xl font-bold text-blue-600">{metrics.security.scansCompleted}</p>
                </div>
                <Zap className="h-12 w-12 text-blue-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Risk Score</p>
                  <p className="text-3xl font-bold text-orange-600">{metrics.security.riskScore}%</p>
                </div>
                <Shield className="h-12 w-12 text-orange-600" />
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Performance Mode */}
      {viewMode === 'performance' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Avg CPU Usage</p>
                  <p className="text-3xl font-bold text-blue-600">{metrics.performance.avgCpu.toFixed(1)}%</p>
                </div>
                <Cpu className="h-12 w-12 text-blue-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Avg Memory Usage</p>
                  <p className="text-3xl font-bold text-green-600">{metrics.performance.avgMemory.toFixed(1)}%</p>
                </div>
                <MemoryStick className="h-12 w-12 text-green-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Storage Usage</p>
                  <p className="text-3xl font-bold text-purple-600">{metrics.performance.avgStorage}%</p>
                </div>
                <HardDrive className="h-12 w-12 text-purple-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Network Utilization</p>
                  <p className="text-3xl font-bold text-orange-600">{metrics.performance.networkUtilization}%</p>
                </div>
                <Network className="h-12 w-12 text-orange-600" />
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Compliance Mode */}
      {viewMode === 'compliance' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Compliance Score</p>
                  <p className="text-3xl font-bold text-green-600">{metrics.compliance.score}%</p>
                </div>
                <CheckCircle className="h-12 w-12 text-green-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Frameworks</p>
                  <p className="text-3xl font-bold text-blue-600">{metrics.compliance.frameworks}</p>
                </div>
                <Layers className="h-12 w-12 text-blue-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Compliant Assets</p>
                  <p className="text-3xl font-bold text-green-600">{metrics.compliance.compliant}</p>
                </div>
                <CheckCircle className="h-12 w-12 text-green-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Non-Compliant</p>
                  <p className="text-3xl font-bold text-red-600">{metrics.compliance.nonCompliant}</p>
                </div>
                <AlertTriangle className="h-12 w-12 text-red-600" />
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default InnovativeDashboard;
