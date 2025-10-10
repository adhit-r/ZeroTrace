import React, { useState, useEffect } from 'react';
import KPIRibbon from './KPIRibbon';
import BranchSelector from './BranchSelector';
import RiskHeatmap from './RiskHeatmap';
import AssetInventory from './AssetInventory';
import { 
  AlertTriangle, 
  Clock, 
  Shield, 
  Target, 
  TrendingUp,
  Building2,
  Users,
  Activity
} from 'lucide-react';

interface DashboardData {
  kpis: Array<{
    id: string;
    label: string;
    value: number | string;
    delta?: {
      value: number;
      trend: 'up' | 'down' | 'neutral';
      period: string;
    };
    sparkline?: number[];
    color?: 'default' | 'critical' | 'warning' | 'success';
  }>;
  branches: Array<{
    id: string;
    name: string;
    location: string;
    type: 'headquarters' | 'branch' | 'datacenter' | 'cloud';
    status: 'active' | 'inactive' | 'maintenance';
    metrics: {
      totalAssets: number;
      criticalVulns: number;
      complianceScore: number;
      lastScan: string;
    };
    coordinates: { lat: number; lng: number };
  }>;
  assets: Array<{
    id: string;
    hostname: string;
    ip: string;
    branch: string;
    location: string;
    owner: string;
    businessCriticality: 'critical' | 'high' | 'medium' | 'low';
    tags: string[];
    lastSeen: string;
    agentStatus: 'online' | 'offline' | 'maintenance';
    vulnerabilities: {
      critical: number;
      high: number;
      medium: number;
      low: number;
    };
    complianceScore: number;
    riskScore: number;
    suggestedFixes: number;
  }>;
  recentAlerts: Array<{
    id: string;
    type: 'critical' | 'high' | 'medium' | 'low';
    title: string;
    description: string;
    timestamp: string;
    branch: string;
    asset: string;
  }>;
  topRiskyAssets: Array<{
    id: string;
    hostname: string;
    riskScore: number;
    criticalVulns: number;
    branch: string;
  }>;
}

interface GlobalDashboardProps {
  userRole: 'global_ciso' | 'branch_ciso' | 'branch_it_manager' | 'security_analyst' | 'patch_engineer';
  className?: string;
}

const GlobalDashboard: React.FC<GlobalDashboardProps> = ({ userRole, className = '' }) => {
  const [selectedBranchId, setSelectedBranchId] = useState<string>();
  const [dashboardData, setDashboardData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'overview' | 'assets' | 'branches'>('overview');

  // Mock data - replace with actual API calls
  useEffect(() => {
    const fetchDashboardData = async () => {
      setLoading(true);
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const mockData: DashboardData = {
        kpis: [
          {
            id: 'active_critical_cves',
            label: 'Active Critical CVEs',
            value: 47,
            delta: { value: 12, trend: 'up', period: 'last 7 days' },
            sparkline: [45, 42, 38, 41, 44, 47, 43],
            color: 'critical'
          },
          {
            id: 'mean_time_to_remediate',
            label: 'MTTR (days)',
            value: 12.5,
            delta: { value: 8, trend: 'down', period: 'last 30 days' },
            sparkline: [18, 16, 14, 13, 12.5, 12.8, 12.5],
            color: 'success'
          },
          {
            id: 'compliance_percent',
            label: 'Compliance %',
            value: 87,
            delta: { value: 3, trend: 'up', period: 'last 30 days' },
            sparkline: [82, 84, 85, 86, 87, 86.5, 87],
            color: 'success'
          },
          {
            id: 'scan_coverage',
            label: 'Scan Coverage %',
            value: 94,
            delta: { value: 2, trend: 'up', period: 'last 7 days' },
            sparkline: [89, 91, 92, 93, 94, 93.5, 94],
            color: 'success'
          },
          {
            id: 'exploit_presence',
            label: 'Exploitable CVEs',
            value: 23,
            delta: { value: 5, trend: 'up', period: 'last 7 days' },
            sparkline: [18, 19, 21, 20, 23, 22, 23],
            color: 'warning'
          },
          {
            id: 'total_assets',
            label: 'Total Assets',
            value: 1247,
            delta: { value: 12, trend: 'up', period: 'last 30 days' },
            sparkline: [1200, 1210, 1220, 1230, 1240, 1245, 1247],
            color: 'default'
          }
        ],
        branches: [
          {
            id: 'hq-nyc',
            name: 'Headquarters NYC',
            location: 'New York, NY',
            type: 'headquarters',
            status: 'active',
            metrics: { totalAssets: 450, criticalVulns: 12, complianceScore: 92, lastScan: '2025-01-09T10:30:00Z' },
            coordinates: { lat: 40.7128, lng: -74.0060 }
          },
          {
            id: 'branch-london',
            name: 'London Branch',
            location: 'London, UK',
            type: 'branch',
            status: 'active',
            metrics: { totalAssets: 280, criticalVulns: 8, complianceScore: 89, lastScan: '2025-01-09T08:15:00Z' },
            coordinates: { lat: 51.5074, lng: -0.1278 }
          },
          {
            id: 'datacenter-tx',
            name: 'Texas Datacenter',
            location: 'Austin, TX',
            type: 'datacenter',
            status: 'active',
            metrics: { totalAssets: 320, criticalVulns: 15, complianceScore: 85, lastScan: '2025-01-09T09:45:00Z' },
            coordinates: { lat: 30.2672, lng: -97.7431 }
          },
          {
            id: 'cloud-aws',
            name: 'AWS Cloud',
            location: 'Global',
            type: 'cloud',
            status: 'active',
            metrics: { totalAssets: 197, criticalVulns: 6, complianceScore: 94, lastScan: '2025-01-09T11:20:00Z' },
            coordinates: { lat: 39.8283, lng: -98.5795 }
          }
        ],
        assets: [
          {
            id: 'asset-001',
            hostname: 'web-server-01',
            ip: '192.168.1.100',
            branch: 'Headquarters NYC',
            location: 'New York, NY',
            owner: 'IT Team',
            businessCriticality: 'critical',
            tags: ['web', 'production'],
            lastSeen: '2025-01-09T10:30:00Z',
            agentStatus: 'online',
            vulnerabilities: { critical: 3, high: 5, medium: 8, low: 12 },
            complianceScore: 78,
            riskScore: 85,
            suggestedFixes: 7
          },
          {
            id: 'asset-002',
            hostname: 'db-server-02',
            ip: '192.168.1.101',
            branch: 'Headquarters NYC',
            location: 'New York, NY',
            owner: 'Database Team',
            businessCriticality: 'critical',
            tags: ['database', 'production'],
            lastSeen: '2025-01-09T10:25:00Z',
            agentStatus: 'online',
            vulnerabilities: { critical: 2, high: 3, medium: 6, low: 9 },
            complianceScore: 82,
            riskScore: 78,
            suggestedFixes: 5
          }
        ],
        recentAlerts: [
          {
            id: 'alert-001',
            type: 'critical',
            title: 'Critical CVE-2024-1234 Detected',
            description: 'Remote code execution vulnerability found in Apache HTTP Server',
            timestamp: '2025-01-09T10:30:00Z',
            branch: 'Headquarters NYC',
            asset: 'web-server-01'
          },
          {
            id: 'alert-002',
            type: 'high',
            title: 'High Risk CVE-2024-5678',
            description: 'Privilege escalation vulnerability in Linux kernel',
            timestamp: '2025-01-09T09:15:00Z',
            branch: 'London Branch',
            asset: 'app-server-03'
          }
        ],
        topRiskyAssets: [
          { id: 'asset-001', hostname: 'web-server-01', riskScore: 85, criticalVulns: 3, branch: 'Headquarters NYC' },
          { id: 'asset-003', hostname: 'app-server-02', riskScore: 82, criticalVulns: 2, branch: 'London Branch' },
          { id: 'asset-004', hostname: 'db-server-03', riskScore: 79, criticalVulns: 1, branch: 'Texas Datacenter' }
        ]
      };
      
      setDashboardData(mockData);
      setLoading(false);
    };

    fetchDashboardData();
  }, []);

  const handleBulkAction = (action: string, assetIds: string[]) => {
    console.log(`Bulk action: ${action} on assets:`, assetIds);
    // Implement bulk actions
  };

  const handleAssetSelect = (assetId: string) => {
    console.log('Selected asset:', assetId);
    // Navigate to asset detail
  };

  if (loading) {
    return (
      <div className={`p-6 ${className}`}>
        <div className="animate-pulse space-y-6">
          <div className="h-32 bg-gray-200 rounded border-3 border-black"></div>
          <div className="h-64 bg-gray-200 rounded border-3 border-black"></div>
          <div className="h-96 bg-gray-200 rounded border-3 border-black"></div>
        </div>
      </div>
    );
  }

  if (!dashboardData) {
    return (
      <div className={`p-6 ${className}`}>
        <div className="text-center">
          <h2 className="text-2xl font-bold text-red-600 mb-4">Failed to load dashboard data</h2>
          <p className="text-gray-600">Please check your connection and try again.</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`p-6 space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold uppercase tracking-wider">ZeroTrace Dashboard</h1>
          <p className="text-gray-600">Global security overview and asset management</p>
        </div>
        <div className="flex items-center gap-3">
          <BranchSelector
            branches={dashboardData.branches}
            selectedBranchId={selectedBranchId}
            onBranchSelect={setSelectedBranchId}
            userRole={userRole}
            className="w-80"
          />
        </div>
      </div>

      {/* KPI Ribbon */}
      <KPIRibbon metrics={dashboardData.kpis} />

      {/* Tab Navigation */}
      <div className="flex border-2 border-black rounded">
        <button
          onClick={() => setActiveTab('overview')}
          className={`px-6 py-3 text-sm font-bold uppercase tracking-wider transition-colors ${
            activeTab === 'overview' 
              ? 'bg-orange-100 text-orange-800 border-r-2 border-black' 
              : 'bg-white text-gray-600 hover:bg-gray-50'
          }`}
        >
          Overview
        </button>
        <button
          onClick={() => setActiveTab('assets')}
          className={`px-6 py-3 text-sm font-bold uppercase tracking-wider transition-colors ${
            activeTab === 'assets' 
              ? 'bg-orange-100 text-orange-800 border-r-2 border-black' 
              : 'bg-white text-gray-600 hover:bg-gray-50'
          }`}
        >
          Assets
        </button>
        <button
          onClick={() => setActiveTab('branches')}
          className={`px-6 py-3 text-sm font-bold uppercase tracking-wider transition-colors ${
            activeTab === 'branches' 
              ? 'bg-orange-100 text-orange-800' 
              : 'bg-white text-gray-600 hover:bg-gray-50'
          }`}
        >
          Branches
        </button>
      </div>

      {/* Tab Content */}
      {activeTab === 'overview' && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Risk Heatmap */}
          <div className="lg:col-span-2">
            <RiskHeatmap
              data={dashboardData.branches}
              selectedBranchId={selectedBranchId}
              onBranchSelect={setSelectedBranchId}
            />
          </div>

          {/* Recent Alerts */}
          <div className="p-6 border-3 border-black bg-white rounded shadow-lg">
            <h3 className="text-xl font-bold uppercase tracking-wider mb-4 flex items-center gap-2">
              <AlertTriangle className="w-5 h-5" />
              Recent Critical Alerts
            </h3>
            <div className="space-y-3">
              {dashboardData.recentAlerts.map((alert) => (
                <div
                  key={alert.id}
                  className={`p-3 rounded border-2 ${
                    alert.type === 'critical' 
                      ? 'border-red-500 bg-red-50' 
                      : 'border-orange-500 bg-orange-50'
                  }`}
                >
                  <div className="flex items-center justify-between mb-2">
                    <span className="font-bold text-sm">{alert.title}</span>
                    <span className="text-xs text-gray-500">
                      {new Date(alert.timestamp).toLocaleString()}
                    </span>
                  </div>
                  <p className="text-sm text-gray-600 mb-2">{alert.description}</p>
                  <div className="flex items-center gap-2 text-xs text-gray-500">
                    <Building2 className="w-3 h-3" />
                    <span>{alert.branch}</span>
                    <span>•</span>
                    <span>{alert.asset}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Top Risky Assets */}
          <div className="p-6 border-3 border-black bg-white rounded shadow-lg">
            <h3 className="text-xl font-bold uppercase tracking-wider mb-4 flex items-center gap-2">
              <Target className="w-5 h-5" />
              Top Risky Assets
            </h3>
            <div className="space-y-3">
              {dashboardData.topRiskyAssets.map((asset, index) => (
                <div
                  key={asset.id}
                  className="flex items-center justify-between p-3 border-2 border-black bg-white rounded hover:bg-gray-50 cursor-pointer"
                >
                  <div className="flex items-center gap-3">
                    <span className="w-6 h-6 bg-orange-100 text-orange-800 rounded-full flex items-center justify-center text-xs font-bold">
                      {index + 1}
                    </span>
                    <div>
                      <div className="font-bold">{asset.hostname}</div>
                      <div className="text-sm text-gray-600">{asset.branch}</div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="font-bold text-red-600">{asset.riskScore}</div>
                    <div className="text-xs text-gray-500">{asset.criticalVulns} critical</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {activeTab === 'assets' && (
        <AssetInventory
          assets={dashboardData.assets}
          userRole={userRole}
          onBulkAction={handleBulkAction}
          onAssetSelect={handleAssetSelect}
        />
      )}

      {activeTab === 'branches' && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {dashboardData.branches.map((branch) => (
            <div
              key={branch.id}
              className="p-6 border-3 border-black bg-white rounded shadow-lg hover:shadow-xl hover:translate-x-1 hover:translate-y-1 transition-all duration-150"
            >
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-2">
                  <Building2 className="w-5 h-5" />
                  <h3 className="font-bold">{branch.name}</h3>
                </div>
                <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider ${
                  branch.status === 'active' 
                    ? 'bg-green-100 text-green-800' 
                    : 'bg-red-100 text-red-800'
                }`}>
                  {branch.status}
                </span>
              </div>
              
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-gray-600">Assets:</span>
                  <span className="font-bold">{branch.metrics.totalAssets}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Critical Vulns:</span>
                  <span className="font-bold text-red-600">{branch.metrics.criticalVulns}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Compliance:</span>
                  <span className="font-bold text-green-600">{branch.metrics.complianceScore}%</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Last Scan:</span>
                  <span className="text-sm">{new Date(branch.metrics.lastScan).toLocaleDateString()}</span>
                </div>
              </div>
              
              <div className="mt-4 pt-4 border-t-2 border-gray-200">
                <button className="w-full p-2 text-sm font-bold uppercase tracking-wider bg-orange-100 text-orange-800 border-2 border-orange-300 rounded hover:bg-orange-200 transition-colors">
                  View Details
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default GlobalDashboard;
