import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import {
  Server,
  Cpu,
  HardDrive,
  MemoryStick,
  Monitor,
  MapPin,
  Shield,
  AlertTriangle,
  CheckCircle,
  Clock,
  Activity,
  Wifi,
  Globe,
  Tag,
  RefreshCw,
  Download,
  Eye,
  Settings,
  BarChart3,
  Zap
} from 'lucide-react';

interface AssetDetailData {
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
}

const AssetDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [asset, setAsset] = useState<AssetDetailData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchAssetDetail();
  }, [id]);

  const fetchAssetDetail = async () => {
    try {
      setLoading(true);
      
      // Fetch real agent data from API
      const response = await fetch(`http://localhost:8080/api/agents/`);
      const agentsData = await response.json();
      
      if (!agentsData.success || !agentsData.data) {
        throw new Error('Failed to fetch agents data');
      }
      
      // Find the specific agent by ID
      const agent = agentsData.data.find((a: any) => a.id === id);
      
      if (!agent) {
        throw new Error('Agent not found');
      }
      
      
      // Transform API data to AssetDetailData format
      // Check both main agent fields and metadata for system information
      const assetData: AssetDetailData = {
        id: agent.id,
        hostname: agent.hostname || 'Unknown',
        ipAddress: agent.ip_address || 'Unknown',
        macAddress: agent.mac_address || agent.metadata?.mac_address || 'Unknown',
        osName: agent.os || agent.metadata?.os_name || 'Unknown',
        osVersion: agent.os_version || agent.metadata?.os_version || 'Unknown',
        osBuild: agent.metadata?.os_build || 'Unknown',
        kernelVersion: agent.metadata?.kernel_version || 'Unknown',
        cpuModel: agent.metadata?.cpu_model || 'Unknown',
        cpuCores: agent.metadata?.cpu_cores || 0,
        memoryTotalGB: agent.metadata?.memory_total_gb || 0,
        storageTotalGB: agent.metadata?.storage_total_gb || 0,
        gpuModel: agent.metadata?.gpu_model || 'Unknown',
        serialNumber: agent.serial_number || agent.metadata?.serial_number || 'Unknown',
        platform: agent.metadata?.platform || 'Unknown',
        city: agent.city || agent.metadata?.city || 'Unknown',
        region: agent.region || agent.metadata?.region || 'Unknown',
        country: agent.country || agent.metadata?.country || 'Unknown',
        timezone: agent.timezone || agent.metadata?.timezone || 'Unknown',
        riskScore: agent.risk_score || agent.metadata?.risk_score || 0,
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
          total: agent.metadata?.total_config_checks || 0,
          issues: agent.metadata?.config_issues || 0
        },
        installedApps: agent.metadata?.dependencies || []
      };
      
      setAsset(assetData);
    } catch (err) {
      setError('Failed to fetch asset details');
      console.error('Error fetching asset:', err);
    } finally {
      setLoading(false);
    }
  };

  const getRiskColor = (score: number) => {
    if (score >= 8) return 'text-red-600 bg-red-100 border-red-300';
    if (score >= 6) return 'text-orange-600 bg-orange-100 border-orange-300';
    if (score >= 4) return 'text-yellow-600 bg-yellow-100 border-yellow-300';
    return 'text-green-600 bg-green-100 border-green-300';
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online': return 'text-green-600 bg-green-100 border-green-300';
      case 'offline': return 'text-red-600 bg-red-100 border-red-300';
      default: return 'text-gray-600 bg-gray-100 border-gray-300';
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <RefreshCw className="h-8 w-8 animate-spin mx-auto mb-4 text-orange-500" />
          <p className="text-gray-600 font-bold">Loading asset details...</p>
        </div>
      </div>
    );
  }

  if (error || !asset) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <AlertTriangle className="h-8 w-8 mx-auto mb-4 text-red-500" />
          <p className="text-red-600 font-bold">{error || 'Asset not found'}</p>
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
              <h1 className="text-3xl font-bold uppercase tracking-wider text-black">
                {asset.hostname}
              </h1>
              <p className="text-gray-600 mt-1">Asset Details & Security Analysis</p>
            </div>
            <div className="flex items-center gap-4">
              <div className={`px-4 py-2 rounded-lg border-3 font-bold uppercase tracking-wider ${getStatusColor(asset.status)}`}>
                {asset.status}
              </div>
              <div className={`px-4 py-2 rounded-lg border-3 font-bold uppercase tracking-wider ${getRiskColor(asset.riskScore)}`}>
                Risk: {asset.riskScore}/10
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto p-6 space-y-6">
        {/* Overview Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {/* Vulnerabilities */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-red-100 rounded border-2 border-black">
                <AlertTriangle className="h-6 w-6 text-red-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-red-600">{asset.vulnerabilities.total}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Vulnerabilities</div>
              </div>
            </div>
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-red-600 font-bold">Critical: {asset.vulnerabilities.critical}</span>
                <span className="text-orange-600 font-bold">High: {asset.vulnerabilities.high}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-yellow-600 font-bold">Medium: {asset.vulnerabilities.medium}</span>
                <span className="text-green-600 font-bold">Low: {asset.vulnerabilities.low}</span>
              </div>
            </div>
          </div>

          {/* Applications */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-blue-100 rounded border-2 border-black">
                <Server className="h-6 w-6 text-blue-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-blue-600">{asset.applications.total}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Applications</div>
              </div>
            </div>
            <div className="text-sm text-orange-600 font-bold">
              {asset.applications.vulnerable} vulnerable
            </div>
          </div>

          {/* Configuration */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-purple-100 rounded border-2 border-black">
                <Settings className="h-6 w-6 text-purple-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-purple-600">{asset.configuration.total}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Config Checks</div>
              </div>
            </div>
            <div className="text-sm text-red-600 font-bold">
              {asset.configuration.issues} issues found
            </div>
          </div>

          {/* Last Seen */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-green-100 rounded border-2 border-black">
                <Clock className="h-6 w-6 text-green-600" />
              </div>
              <div className="text-right">
                <div className="text-sm font-bold text-green-600">Last Seen</div>
                <div className="text-xs text-gray-600">
                  {new Date(asset.lastSeen).toLocaleString()}
                </div>
              </div>
            </div>
            <div className="text-sm text-gray-600">
              Status: <span className={`font-bold ${getStatusColor(asset.status).split(' ')[0]}`}>
                {asset.status.toUpperCase()}
              </span>
            </div>
          </div>
        </div>

        {/* System Information */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Operating System */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <h2 className="text-xl font-black text-black uppercase mb-4 flex items-center gap-2">
              <Monitor className="h-6 w-6" />
              Operating System
            </h2>
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">OS:</span>
                <span className="text-black font-bold">{asset.osName} {asset.osVersion}</span>
              </div>
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">Build:</span>
                <span className="text-black font-bold">{asset.osBuild}</span>
              </div>
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">Kernel:</span>
                <span className="text-black font-bold">{asset.kernelVersion}</span>
              </div>
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">Platform:</span>
                <span className="text-black font-bold">{asset.platform}</span>
              </div>
            </div>
          </div>

          {/* Hardware */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <h2 className="text-xl font-black text-black uppercase mb-4 flex items-center gap-2">
              <Cpu className="h-6 w-6" />
              Hardware
            </h2>
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">CPU:</span>
                <span className="text-black font-bold">{asset.cpuModel}</span>
              </div>
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">Cores:</span>
                <span className="text-black font-bold">{asset.cpuCores}</span>
              </div>
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">Memory:</span>
                <span className="text-black font-bold">{asset.memoryTotalGB} GB</span>
              </div>
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">Storage:</span>
                <span className="text-black font-bold">{asset.storageTotalGB} GB</span>
              </div>
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">GPU:</span>
                <span className="text-black font-bold">{asset.gpuModel}</span>
              </div>
            </div>
          </div>

          {/* Installed Applications */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <h2 className="text-xl font-black text-black uppercase mb-4 flex items-center gap-2">
              <Settings className="h-6 w-6" />
              Installed Applications
            </h2>
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">Total Apps:</span>
                <span className="text-black font-bold">{asset.installedApps.length}</span>
              </div>
              {asset.installedApps.length > 0 && (
                <div className="mt-4">
                  <h3 className="font-bold text-gray-700 mb-2">Application List:</h3>
                  <div className="max-h-64 overflow-y-auto space-y-2">
                    {asset.installedApps.slice(0, 20).map((app, index) => (
                      <div key={index} className="flex justify-between items-center p-2 bg-gray-50 border border-gray-300 rounded">
                        <div>
                          <span className="font-bold text-black">{app.name}</span>
                          {app.version && app.version !== 'unknown' && (
                            <span className="text-gray-600 ml-2">v{app.version}</span>
                          )}
                        </div>
                        {app.vendor && app.vendor !== 'Unknown' && (
                          <span className="text-sm text-gray-500">{app.vendor}</span>
                        )}
                      </div>
                    ))}
                    {asset.installedApps.length > 20 && (
                      <div className="text-center text-gray-500 text-sm">
                        ... and {asset.installedApps.length - 20} more applications
                      </div>
                    )}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Network & Location */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Network Information */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <h2 className="text-xl font-black text-black uppercase mb-4 flex items-center gap-2">
              <Wifi className="h-6 w-6" />
              Network
            </h2>
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">IP Address:</span>
                <span className="text-black font-bold">{asset.ipAddress}</span>
              </div>
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">MAC Address:</span>
                <span className="text-black font-bold font-mono text-sm">{asset.macAddress}</span>
              </div>
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">Serial Number:</span>
                <span className="text-black font-bold font-mono text-sm">{asset.serialNumber}</span>
              </div>
            </div>
          </div>

          {/* Location Information */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <h2 className="text-xl font-black text-black uppercase mb-4 flex items-center gap-2">
              <MapPin className="h-6 w-6" />
              Location
            </h2>
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">City:</span>
                <span className="text-black font-bold">{asset.city}</span>
              </div>
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">Region:</span>
                <span className="text-black font-bold">{asset.region}</span>
              </div>
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">Country:</span>
                <span className="text-black font-bold">{asset.country}</span>
              </div>
              <div className="flex justify-between">
                <span className="font-bold text-gray-700">Timezone:</span>
                <span className="text-black font-bold">{asset.timezone}</span>
              </div>
            </div>
          </div>
        </div>

        {/* Tags */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <h2 className="text-xl font-black text-black uppercase mb-4 flex items-center gap-2">
            <Tag className="h-6 w-6" />
            Tags
          </h2>
          <div className="flex flex-wrap gap-2">
            {asset.tags.map((tag, index) => (
              <span
                key={index}
                className="px-3 py-1 bg-orange-100 text-orange-800 border-2 border-orange-300 rounded-lg font-bold uppercase tracking-wider text-sm"
              >
                {tag}
              </span>
            ))}
          </div>
        </div>

        {/* Actions */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <h2 className="text-xl font-black text-black uppercase mb-4 flex items-center gap-2">
            <Zap className="h-6 w-6" />
            Actions
          </h2>
          <div className="flex flex-wrap gap-4">
            <button className="px-6 py-3 bg-blue-500 text-white border-3 border-black rounded-lg shadow-neo-brutal-small hover:shadow-neo-brutal-small-hover transition-all font-bold uppercase tracking-wider">
              <Eye className="h-4 w-4 mr-2 inline-block" />
              View Vulnerabilities
            </button>
            <button className="px-6 py-3 bg-green-500 text-white border-3 border-black rounded-lg shadow-neo-brutal-small hover:shadow-neo-brutal-small-hover transition-all font-bold uppercase tracking-wider">
              <RefreshCw className="h-4 w-4 mr-2 inline-block" />
              Rescan Asset
            </button>
            <button className="px-6 py-3 bg-purple-500 text-white border-3 border-black rounded-lg shadow-neo-brutal-small hover:shadow-neo-brutal-small-hover transition-all font-bold uppercase tracking-wider">
              <BarChart3 className="h-4 w-4 mr-2 inline-block" />
              View Reports
            </button>
            <button className="px-6 py-3 bg-orange-500 text-white border-3 border-black rounded-lg shadow-neo-brutal-small hover:shadow-neo-brutal-small-hover transition-all font-bold uppercase tracking-wider">
              <Download className="h-4 w-4 mr-2 inline-block" />
              Export Data
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AssetDetail;
