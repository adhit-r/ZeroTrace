import React, { useState, useEffect } from 'react';
import {
  Server,
  Cpu,
  HardDrive,
  MemoryStick,
  Shield,
  AlertTriangle,
  CheckCircle,
  Activity,
  Wifi,
  Globe,
  RefreshCw,
  Download,
  Eye,
  Settings,
  BarChart3,
  Network,
  Layers,
  X
} from 'lucide-react';
import { agentService } from '../../services/agentService';

interface SystemMetrics {
  cpu: {
    usage: number;
    cores: number;
    model: string;
    temperature: number;
    frequency: number;
  };
  memory: {
    total: number;
    used: number;
    available: number;
    usage: number;
  };
  storage: {
    total: number;
    used: number;
    available: number;
    usage: number;
    devices: Array<{
      name: string;
      size: number;
      used: number;
      type: string;
      health: number;
    }>;
  };
  network: {
    interfaces: Array<{
      name: string;
      ip: string;
      mac: string;
      status: 'up' | 'down';
      speed: number;
      bytesIn: number;
      bytesOut: number;
    }>;
    connections: number;
    bandwidth: number;
  };
  processes: {
    total: number;
    running: number;
    sleeping: number;
    zombie: number;
    topProcesses: Array<{
      name: string;
      pid: number;
      cpu: number;
      memory: number;
      user: string;
    }>;
  };
}

interface SecurityMetrics {
  vulnerabilities: {
    total: number;
    critical: number;
    high: number;
    medium: number;
    low: number;
    trend: 'increasing' | 'decreasing' | 'stable';
  };
  compliance: {
    score: number;
    frameworks: Array<{
      name: string;
      score: number;
      status: 'compliant' | 'non_compliant' | 'partial';
    }>;
  };
  risk: {
    score: number;
    factors: Array<{
      name: string;
      weight: number;
      score: number;
      description: string;
    }>;
    trend: 'increasing' | 'decreasing' | 'stable';
  };
  patches: {
    available: number;
    critical: number;
    installed: number;
    pending: number;
  };
}

interface AssetData {
  id: string;
  hostname: string;
  ipAddress: string;
  macAddress: string;
  osName: string;
  osVersion: string;
  osBuild: string;
  kernelVersion: string;
  cpuModel: string;
  cpuCores: number;
  memoryTotalGB: number;
  storageTotalGB: number;
  gpuModel: string;
  serialNumber: string;
  platform: string;
  city: string;
  region: string;
  country: string;
  timezone: string;
  riskScore: number;
  tags: string[];
  lastSeen: string;
  status: string;
  vulnerabilities: {
    total: number;
    critical: number;
    high: number;
    medium: number;
    low: number;
  };
  applications: {
    total: number;
    vulnerable: number;
  };
  configuration: {
    total: number;
    issues: number;
  };
  installedApps: Array<{
    name: string;
    version: string;
    vendor: string;
    path?: string;
  }>;
  systemMetrics: SystemMetrics;
  securityMetrics: SecurityMetrics;
}

interface EnhancedAssetDetailsProps {
  assetId: string;
  onClose?: () => void;
  className?: string;
}

const EnhancedAssetDetails: React.FC<EnhancedAssetDetailsProps> = ({
  assetId,
  onClose,
  className = ''
}) => {
  const [assetData, setAssetData] = useState<AssetData | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'overview' | 'security' | 'performance' | 'applications' | 'network' | 'compliance'>('overview');

  useEffect(() => {
    const fetchAssetData = async () => {
      setLoading(true);
      try {
        // Fetch real agent data from API
        const agents = await agentService.getAgents();

        // Find the specific agent by ID
        const agent = agents.find((a: any) => a.id === assetId);

        if (!agent) {
          throw new Error('Agent not found');
        }

        // Transform API data to enhanced AssetData format
        const assetData: AssetData = {
          id: agent.id,
          hostname: agent.hostname || 'Unknown',
          ipAddress: agent.ip_address || 'Unknown',
          macAddress: agent.mac_address || 'Unknown',
          osName: agent.os || 'Unknown',
          osVersion: agent.os_version || 'Unknown',
          osBuild: agent.metadata?.os_build || 'Unknown',
          kernelVersion: agent.metadata?.kernel_version || 'Unknown',
          cpuModel: agent.metadata?.cpu_model || 'Unknown',
          cpuCores: agent.metadata?.cpu_cores || 0,
          memoryTotalGB: agent.metadata?.memory_total_gb || 0,
          storageTotalGB: agent.metadata?.storage_total_gb || 0,
          gpuModel: agent.metadata?.gpu_model || 'Unknown',
          serialNumber: agent.serial_number || 'Unknown',
          platform: agent.metadata?.platform || 'Unknown',
          city: agent.city || 'Unknown',
          region: agent.region || 'Unknown',
          country: agent.country || 'Unknown',
          timezone: agent.timezone || 'Unknown',
          riskScore: agent.risk_score || 0,
          tags: agent.metadata?.tags || [],
          lastSeen: agent.last_seen || new Date().toISOString(),
          status: agent.status || 'unknown',
          vulnerabilities: {
            total: agent.metadata?.total_vulnerabilities || 0,
            critical: agent.metadata?.critical_vulnerabilities || 0,
            high: agent.metadata?.high_vulnerabilities || 0,
            medium: agent.metadata?.medium_vulnerabilities || 0,
            low: agent.metadata?.low_vulnerabilities || 0
          },
          applications: {
            total: agent.metadata?.total_applications || 0,
            vulnerable: agent.metadata?.vulnerable_applications || 0
          },
          configuration: {
            total: agent.metadata?.total_configurations || 0,
            issues: agent.metadata?.configuration_issues || 0
          },
          installedApps: agent.metadata?.installed_apps || [],
          systemMetrics: {
            cpu: {
              usage: agent.cpu_usage || 0,
              cores: agent.metadata?.cpu_cores || 0,
              model: agent.metadata?.cpu_model || 'Unknown',
              temperature: agent.metadata?.cpu_temperature || 0,
              frequency: agent.metadata?.cpu_frequency || 0
            },
            memory: {
              total: agent.metadata?.memory_total_gb || 0,
              used: agent.metadata?.memory_used_gb || 0,
              available: agent.metadata?.memory_available_gb || 0,
              usage: agent.memory_usage || 0
            },
            storage: {
              total: agent.metadata?.storage_total_gb || 0,
              used: agent.metadata?.storage_used_gb || 0,
              available: agent.metadata?.storage_available_gb || 0,
              usage: agent.metadata?.storage_usage || 0,
              devices: agent.metadata?.storage_devices || []
            },
            network: {
              interfaces: agent.metadata?.network_interfaces || [],
              connections: agent.metadata?.network_connections || 0,
              bandwidth: agent.metadata?.network_bandwidth || 0
            },
            processes: {
              total: agent.metadata?.total_processes || 0,
              running: agent.metadata?.running_processes || 0,
              sleeping: agent.metadata?.sleeping_processes || 0,
              zombie: agent.metadata?.zombie_processes || 0,
              topProcesses: agent.metadata?.top_processes || []
            }
          },
          securityMetrics: {
            vulnerabilities: {
              total: agent.metadata?.total_vulnerabilities || 0,
              critical: agent.metadata?.critical_vulnerabilities || 0,
              high: agent.metadata?.high_vulnerabilities || 0,
              medium: agent.metadata?.medium_vulnerabilities || 0,
              low: agent.metadata?.low_vulnerabilities || 0,
              trend: agent.metadata?.vulnerability_trend || 'stable'
            },
            compliance: {
              score: agent.metadata?.compliance_score || 0,
              frameworks: agent.metadata?.compliance_frameworks || []
            },
            risk: {
              score: agent.risk_score || 0,
              factors: agent.metadata?.risk_factors || [],
              trend: agent.metadata?.risk_trend || 'stable'
            },
            patches: {
              available: agent.metadata?.patches_available || 0,
              critical: agent.metadata?.critical_patches || 0,
              installed: agent.metadata?.patches_installed || 0,
              pending: agent.metadata?.patches_pending || 0
            }
          }
        };

        setAssetData(assetData);
      } catch (error) {
        console.error('Failed to fetch asset data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchAssetData();
  }, [assetId]);

  const getRiskColor = (score: number) => {
    if (score >= 80) return 'text-red-600 bg-red-100';
    if (score >= 60) return 'text-orange-600 bg-orange-100';
    if (score >= 40) return 'text-yellow-600 bg-yellow-100';
    return 'text-green-600 bg-green-100';
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online': return 'text-green-600 bg-green-100';
      case 'offline': return 'text-red-600 bg-red-100';
      case 'maintenance': return 'text-yellow-600 bg-yellow-100';
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

  if (loading) {
    return (
      <div className={`p-6 ${className}`}>
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  if (!assetData) {
    return (
      <div className={`p-6 ${className}`}>
        <div className="text-center text-gray-500">
          <Server className="h-12 w-12 mx-auto mb-4" />
          <p>Asset data not available</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">{assetData.hostname}</h2>
          <p className="text-gray-600">
            {assetData.ipAddress} • {assetData.osName} {assetData.osVersion}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <span className={`px-3 py-1 rounded text-sm font-bold ${getStatusColor(assetData.status)}`}>
            {assetData.status.toUpperCase()}
          </span>
          <span className={`px-3 py-1 rounded text-sm font-bold ${getRiskColor(assetData.riskScore)}`}>
            Risk: {assetData.riskScore}%
          </span>
          {onClose && (
            <button
              onClick={onClose}
              className="p-2 hover:bg-gray-100 rounded border-2 border-black"
            >
              <X className="h-4 w-4" />
            </button>
          )}
        </div>
      </div>

      {/* Navigation Tabs */}
      <div className="flex space-x-1 bg-gray-100 p-1 rounded-lg">
        {[
          { id: 'overview', label: 'Overview', icon: BarChart3 },
          { id: 'security', label: 'Security', icon: Shield },
          { id: 'performance', label: 'Performance', icon: Activity },
          { id: 'applications', label: 'Applications', icon: Layers },
          { id: 'network', label: 'Network', icon: Network },
          { id: 'compliance', label: 'Compliance', icon: CheckCircle }
        ].map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id as any)}
            className={`flex items-center gap-2 px-4 py-2 rounded-md transition-colors ${activeTab === tab.id
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
              }`}
          >
            <tab.icon className="h-4 w-4" />
            {tab.label}
          </button>
        ))}
      </div>

      {/* Content based on active tab */}
      {activeTab === 'overview' && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* System Information */}
          <div className="lg:col-span-2 space-y-6">
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
                <Server className="h-5 w-5" />
                System Information
              </h3>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-gray-600">Hostname</p>
                  <p className="font-bold">{assetData.hostname}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">IP Address</p>
                  <p className="font-bold">{assetData.ipAddress}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">Operating System</p>
                  <p className="font-bold">{assetData.osName} {assetData.osVersion}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">Platform</p>
                  <p className="font-bold">{assetData.platform}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">CPU</p>
                  <p className="font-bold">{assetData.cpuModel} ({assetData.cpuCores} cores)</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">Memory</p>
                  <p className="font-bold">{assetData.memoryTotalGB}GB</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">Storage</p>
                  <p className="font-bold">{assetData.storageTotalGB}GB</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">Last Seen</p>
                  <p className="font-bold">{new Date(assetData.lastSeen).toLocaleString()}</p>
                </div>
              </div>
            </div>

            {/* Performance Metrics */}
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
                <Activity className="h-5 w-5" />
                Performance Metrics
              </h3>
              <div className="grid grid-cols-2 gap-6">
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm text-gray-600">CPU Usage</span>
                    <span className="font-bold">{assetData.systemMetrics.cpu.usage}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-blue-600 h-2 rounded-full"
                      style={{ width: `${assetData.systemMetrics.cpu.usage}%` }}
                    ></div>
                  </div>
                </div>
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm text-gray-600">Memory Usage</span>
                    <span className="font-bold">{assetData.systemMetrics.memory.usage}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-green-600 h-2 rounded-full"
                      style={{ width: `${assetData.systemMetrics.memory.usage}%` }}
                    ></div>
                  </div>
                </div>
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm text-gray-600">Storage Usage</span>
                    <span className="font-bold">{assetData.systemMetrics.storage.usage}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-purple-600 h-2 rounded-full"
                      style={{ width: `${assetData.systemMetrics.storage.usage}%` }}
                    ></div>
                  </div>
                </div>
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm text-gray-600">Processes</span>
                    <span className="font-bold">{assetData.systemMetrics.processes.total}</span>
                  </div>
                  <div className="text-sm text-gray-600">
                    Running: {assetData.systemMetrics.processes.running}
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Security Overview */}
          <div className="space-y-6">
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
                <Shield className="h-5 w-5" />
                Security Overview
              </h3>
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-gray-600">Risk Score</span>
                  <span className={`px-2 py-1 rounded text-sm font-bold ${getRiskColor(assetData.riskScore)}`}>
                    {assetData.riskScore}%
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-gray-600">Vulnerabilities</span>
                  <span className="font-bold text-red-600">{assetData.vulnerabilities.total}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-gray-600">Critical</span>
                  <span className="font-bold text-red-600">{assetData.vulnerabilities.critical}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-gray-600">High</span>
                  <span className="font-bold text-orange-600">{assetData.vulnerabilities.high}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-gray-600">Compliance Score</span>
                  <span className="font-bold text-green-600">{assetData.securityMetrics.compliance.score}%</span>
                </div>
              </div>
            </div>

            {/* Quick Actions */}
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <h3 className="text-lg font-bold mb-4">Quick Actions</h3>
              <div className="space-y-2">
                <button className="w-full px-4 py-2 bg-blue-600 text-white rounded border-2 border-blue-700 hover:bg-blue-700 transition-colors">
                  <RefreshCw className="h-4 w-4 inline mr-2" />
                  Scan Now
                </button>
                <button className="w-full px-4 py-2 bg-green-600 text-white rounded border-2 border-green-700 hover:bg-green-700 transition-colors">
                  <Download className="h-4 w-4 inline mr-2" />
                  Download Report
                </button>
                <button className="w-full px-4 py-2 bg-orange-600 text-white rounded border-2 border-orange-700 hover:bg-orange-700 transition-colors">
                  <Settings className="h-4 w-4 inline mr-2" />
                  Configure
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'security' && (
        <div className="space-y-6">
          {/* Security Metrics */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="p-4 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Total Vulnerabilities</p>
                  <p className="text-2xl font-bold text-red-600">{assetData.vulnerabilities.total}</p>
                </div>
                <AlertTriangle className="h-8 w-8 text-red-600" />
              </div>
            </div>
            <div className="p-4 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Critical</p>
                  <p className="text-2xl font-bold text-red-600">{assetData.vulnerabilities.critical}</p>
                </div>
                <AlertTriangle className="h-8 w-8 text-red-600" />
              </div>
            </div>
            <div className="p-4 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Compliance Score</p>
                  <p className="text-2xl font-bold text-green-600">{assetData.securityMetrics.compliance.score}%</p>
                </div>
                <CheckCircle className="h-8 w-8 text-green-600" />
              </div>
            </div>
            <div className="p-4 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Patches Available</p>
                  <p className="text-2xl font-bold text-orange-600">{assetData.securityMetrics.patches.available}</p>
                </div>
                <Download className="h-8 w-8 text-orange-600" />
              </div>
            </div>
          </div>

          {/* Vulnerability Breakdown */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
            <h3 className="text-lg font-bold mb-4">Vulnerability Breakdown</h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {Object.entries(assetData.vulnerabilities).map(([severity, count]) => (
                <div key={severity} className="text-center">
                  <div className={`inline-flex items-center px-3 py-1 rounded text-sm font-bold ${getSeverityColor(severity)}`}>
                    {severity.toUpperCase()}
                  </div>
                  <p className="text-2xl font-bold mt-2">{count}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Risk Factors */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
            <h3 className="text-lg font-bold mb-4">Risk Factors</h3>
            <div className="space-y-3">
              {assetData.securityMetrics.risk.factors.map((factor, index) => (
                <div key={index} className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">{factor.name}</p>
                    <p className="text-sm text-gray-600">{factor.description}</p>
                  </div>
                  <div className="text-right">
                    <p className="font-bold">{factor.score}</p>
                    <p className="text-sm text-gray-600">Weight: {factor.weight}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {activeTab === 'performance' && (
        <div className="space-y-6">
          {/* Performance Metrics */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
                <Cpu className="h-5 w-5" />
                CPU Performance
              </h3>
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-gray-600">Usage</span>
                  <span className="font-bold">{assetData.systemMetrics.cpu.usage}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-blue-600 h-2 rounded-full"
                    style={{ width: `${assetData.systemMetrics.cpu.usage}%` }}
                  ></div>
                </div>
                <div className="text-sm text-gray-600">
                  Temperature: {assetData.systemMetrics.cpu.temperature}°C
                </div>
                <div className="text-sm text-gray-600">
                  Frequency: {assetData.systemMetrics.cpu.frequency}MHz
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
                  <span className="text-sm text-gray-600">Usage</span>
                  <span className="font-bold">{assetData.systemMetrics.memory.usage}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-green-600 h-2 rounded-full"
                    style={{ width: `${assetData.systemMetrics.memory.usage}%` }}
                  ></div>
                </div>
                <div className="text-sm text-gray-600">
                  Used: {assetData.systemMetrics.memory.used}GB
                </div>
                <div className="text-sm text-gray-600">
                  Available: {assetData.systemMetrics.memory.available}GB
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
                  <span className="text-sm text-gray-600">Usage</span>
                  <span className="font-bold">{assetData.systemMetrics.storage.usage}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-purple-600 h-2 rounded-full"
                    style={{ width: `${assetData.systemMetrics.storage.usage}%` }}
                  ></div>
                </div>
                <div className="text-sm text-gray-600">
                  Used: {assetData.systemMetrics.storage.used}GB
                </div>
                <div className="text-sm text-gray-600">
                  Available: {assetData.systemMetrics.storage.available}GB
                </div>
              </div>
            </div>
          </div>

          {/* Top Processes */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
            <h3 className="text-lg font-bold mb-4">Top Processes</h3>
            <div className="space-y-2">
              {assetData.systemMetrics.processes.topProcesses.map((process, index) => (
                <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                  <div>
                    <p className="font-medium">{process.name}</p>
                    <p className="text-sm text-gray-600">PID: {process.pid} • User: {process.user}</p>
                  </div>
                  <div className="text-right">
                    <p className="font-bold">CPU: {process.cpu}%</p>
                    <p className="text-sm text-gray-600">Memory: {process.memory}%</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {activeTab === 'applications' && (
        <div className="space-y-6">
          {/* Application Summary */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="p-4 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Total Applications</p>
                  <p className="text-2xl font-bold">{assetData.applications.total}</p>
                </div>
                <Layers className="h-8 w-8 text-blue-600" />
              </div>
            </div>
            <div className="p-4 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Vulnerable Apps</p>
                  <p className="text-2xl font-bold text-red-600">{assetData.applications.vulnerable}</p>
                </div>
                <AlertTriangle className="h-8 w-8 text-red-600" />
              </div>
            </div>
            <div className="p-4 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Configuration Issues</p>
                  <p className="text-2xl font-bold text-orange-600">{assetData.configuration.issues}</p>
                </div>
                <Settings className="h-8 w-8 text-orange-600" />
              </div>
            </div>
          </div>

          {/* Installed Applications */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
            <h3 className="text-lg font-bold mb-4">Installed Applications</h3>
            <div className="space-y-2">
              {assetData.installedApps.map((app, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded border-2 border-gray-200">
                  <div>
                    <p className="font-medium">{app.name}</p>
                    <p className="text-sm text-gray-600">{app.vendor} • v{app.version}</p>
                    {app.path && <p className="text-xs text-gray-500">{app.path}</p>}
                  </div>
                  <div className="flex items-center gap-2">
                    <button className="p-1 text-blue-600 hover:bg-blue-100 rounded">
                      <Eye className="h-4 w-4" />
                    </button>
                    <button className="p-1 text-gray-600 hover:bg-gray-100 rounded">
                      <Settings className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {activeTab === 'network' && (
        <div className="space-y-6">
          {/* Network Interfaces */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
            <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
              <Network className="h-5 w-5" />
              Network Interfaces
            </h3>
            <div className="space-y-3">
              {assetData.systemMetrics.network.interfaces.map((iface, index) => (
                <div key={index} className="p-3 bg-gray-50 rounded border-2 border-gray-200">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium">{iface.name}</p>
                      <p className="text-sm text-gray-600">IP: {iface.ip} • MAC: {iface.mac}</p>
                    </div>
                    <div className="text-right">
                      <span className={`px-2 py-1 rounded text-xs font-bold ${iface.status === 'up' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                        }`}>
                        {iface.status.toUpperCase()}
                      </span>
                      <p className="text-sm text-gray-600 mt-1">{iface.speed} Mbps</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Network Statistics */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="p-4 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Active Connections</p>
                  <p className="text-2xl font-bold">{assetData.systemMetrics.network.connections}</p>
                </div>
                <Wifi className="h-8 w-8 text-blue-600" />
              </div>
            </div>
            <div className="p-4 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Bandwidth</p>
                  <p className="text-2xl font-bold">{assetData.systemMetrics.network.bandwidth} Mbps</p>
                </div>
                <Globe className="h-8 w-8 text-green-600" />
              </div>
            </div>
            <div className="p-4 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Interfaces</p>
                  <p className="text-2xl font-bold">{assetData.systemMetrics.network.interfaces.length}</p>
                </div>
                <Network className="h-8 w-8 text-purple-600" />
              </div>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'compliance' && (
        <div className="space-y-6">
          {/* Compliance Overview */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
            <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
              <CheckCircle className="h-5 w-5" />
              Compliance Overview
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm text-gray-600">Overall Score</span>
                  <span className="text-2xl font-bold text-green-600">{assetData.securityMetrics.compliance.score}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-3">
                  <div
                    className="bg-green-600 h-3 rounded-full"
                    style={{ width: `${assetData.securityMetrics.compliance.score}%` }}
                  ></div>
                </div>
              </div>
              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm text-gray-600">Frameworks</span>
                  <span className="text-2xl font-bold">{assetData.securityMetrics.compliance.frameworks.length}</span>
                </div>
                <div className="text-sm text-gray-600">
                  Active compliance frameworks
                </div>
              </div>
            </div>
          </div>

          {/* Compliance Frameworks */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
            <h3 className="text-lg font-bold mb-4">Compliance Frameworks</h3>
            <div className="space-y-3">
              {assetData.securityMetrics.compliance.frameworks.map((framework, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded border-2 border-gray-200">
                  <div>
                    <p className="font-medium">{framework.name}</p>
                    <p className="text-sm text-gray-600">Score: {framework.score}%</p>
                  </div>
                  <span className={`px-3 py-1 rounded text-sm font-bold ${framework.status === 'compliant' ? 'bg-green-100 text-green-800' :
                      framework.status === 'non_compliant' ? 'bg-red-100 text-red-800' :
                        'bg-yellow-100 text-yellow-800'
                    }`}>
                    {framework.status.replace('_', ' ').toUpperCase()}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default EnhancedAssetDetails;
