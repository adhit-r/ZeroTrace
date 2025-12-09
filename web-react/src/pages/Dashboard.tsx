import React, { useState, useEffect } from 'react';
import {
  Shield,
  AlertTriangle,
  TrendingUp,
  TrendingDown,
  Clock,
  Activity,
  Zap,
  Server,
  RefreshCw
} from 'lucide-react';
import { agentService } from '../services/agentService';
import { useQuery } from '@tanstack/react-query';
import { AssetAppList } from '../components/dashboard/AssetAppList';
import NetworkSecurityDashboard from '../components/dashboard/NetworkSecurityDashboard';
import TopVulnerableAssets from '../components/dashboard/TopVulnerableAssets';
import VulnerabilityScanner from '../components/dashboard/VulnerabilityScanner';

let ChartComponents: any = null;

const loadChartJS = async () => {
  const { Chart, registerables } = await import('chart.js');
  Chart.register(...registerables);
};



const loadCharts = async () => {
  if (!ChartComponents) {
    await loadChartJS();
    const chartModule = await import('react-chartjs-2');
    ChartComponents = {
      Bar: chartModule.Bar,
      Doughnut: chartModule.Doughnut,
      Line: chartModule.Line,
    };
  }
  return ChartComponents;
};

const Dashboard: React.FC = () => {
  const [isScanning, setIsScanning] = useState(false);

  // Use React Query for data fetching
  const {
    data: overview,
    isLoading: isOverviewLoading,
    error: overviewError,
    refetch: refetchOverview
  } = useQuery({
    queryKey: ['dashboard-overview'],
    queryFn: agentService.getDashboardOverview,
    refetchInterval: 30000, // Poll every 30s
  });

  const {
    data: vulnerabilities,
    isLoading: isVulnsLoading,
    error: vulnsError,
    refetch: refetchVulns
  } = useQuery({
    queryKey: ['vulnerabilities'],
    queryFn: agentService.getVulnerabilities,
    refetchInterval: 30000,
  });

  const {
    data: agentStats,
    isLoading: isStatsLoading
  } = useQuery({
    queryKey: ['agent-stats'],
    queryFn: agentService.getAgentStats,
    initialData: { total: 0, online: 0, offline: 0, avgCpu: 0, avgMemory: 0 }
  });

  const {
    data: agents = [],
    isLoading: isAgentsLoading
  } = useQuery({
    queryKey: ['agents'],
    queryFn: agentService.getAgents,
    refetchInterval: 30000,
  });

  // Load Charts
  const [Charts, setCharts] = useState<{ Bar: any; Doughnut: any; Line: any } | null>(null);
  useEffect(() => {
    loadCharts().then(setCharts);
  }, []);

  const isLoading = isOverviewLoading || isVulnsLoading || isStatsLoading || isAgentsLoading || !Charts;
  const error = overviewError || vulnsError;

  const refresh = async () => {
    await Promise.all([refetchOverview(), refetchVulns()]);
  };

  const startScan = async () => {
    setIsScanning(true);
    try {
      await refresh();
    } catch (error) {
      console.error('Failed to trigger scan:', error);
    } finally {
      setIsScanning(false);
    }
  };

  const deployAgent = () => {
    const agentInstructions = `
# ZeroTrace Agent Installation

## macOS
1. Download the agent binary
2. Run: chmod +x zerotrace-agent
3. Run: ./zerotrace-agent

## Linux
1. Download the agent binary
2. Run: chmod +x zerotrace-agent
3. Run: ./zerotrace-agent

## Windows
1. Download the agent executable
2. Run as Administrator: zerotrace-agent.exe

The agent will automatically register with the ZeroTrace API.
    `;

    // Create a temporary file with instructions
    const blob = new Blob([agentInstructions], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'zerotrace-agent-instructions.txt';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);

    alert('Agent installation instructions have been downloaded.');
  };

  const getSeverityColor = (severity: string) => {
    switch (severity?.toLowerCase()) {
      case 'critical': return 'bg-red-100 text-red-800 border-red-300';
      case 'high': return 'bg-orange-100 text-orange-800 border-orange-300';
      case 'medium': return 'bg-yellow-100 text-yellow-800 border-yellow-300';
      case 'low': return 'bg-green-100 text-green-800 border-green-300';
      default: return 'bg-gray-100 text-gray-800 border-gray-300';
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

  if (error || !overview) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="p-6 bg-red-50 border-3 border-red-300 rounded shadow-neubrutalist-lg">
            <div className="flex items-center gap-3 mb-4">
              <AlertTriangle className="h-6 w-6 text-red-600" />
              <h2 className="text-xl font-bold text-red-800">Error Loading Dashboard</h2>
            </div>
            <p className="text-red-700 mb-4">{error instanceof Error ? error.message : 'Unknown error'}</p>
            <button
              onClick={() => refresh()}
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

  // Extract data
  // Safe defaults if API returns partial data
  const totalAssets = overview.total_assets || 0;
  const vulnerableAssets = overview.vulnerable_assets || 0;
  const criticalVulns = overview.critical_vulnerabilities || 0;
  const highVulns = overview.high_vulnerabilities || 0;
  const mediumVulns = overview.medium_vulnerabilities || 0;
  const lowVulns = overview.low_vulnerabilities || 0;
  const totalVulns = overview.total_vulnerabilities || 0;

  if (totalAssets === 0) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="text-center py-20">
            <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg max-w-2xl mx-auto">
              <div className="mb-6">
                <Server className="h-16 w-16 text-gray-400 mx-auto mb-4" />
                <h1 className="text-3xl font-bold text-black mb-2">No Agents Deployed</h1>
                <p className="text-gray-600 text-lg">
                  Deploy the ZeroTrace agent on your systems to start monitoring vulnerabilities
                </p>
              </div>

              <div className="space-y-4">
                <button
                  onClick={deployAgent}
                  className="px-8 py-4 bg-orange-500 text-white border-3 border-black rounded shadow-neubrutalist-lg hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-neubrutalist-md transition-all duration-150 ease-in-out font-bold uppercase tracking-wider"
                >
                  <Activity className="h-5 w-5 mr-2 inline-block" />
                  DEPLOY AGENT
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

  // Dynamic Chart Data
  const vulnerabilityTrendData = {
    labels: ['Last 7 Days', 'Last 6 Days', 'Last 5 Days', 'Last 4 Days', 'Last 3 Days', 'Last 2 Days', 'Today'],
    datasets: [
      {
        label: 'Vulnerabilities Found',
        data: [0, 0, 0, 0, 0, 0, totalVulns],
        fill: false,
        backgroundColor: 'rgb(255, 99, 132)',
        borderColor: 'rgba(255, 99, 132, 0.8)',
        borderWidth: 3,
        tension: 0.4
      }
    ]
  };

  const severityBreakdownData = {
    labels: ['Critical', 'High', 'Medium', 'Low'],
    datasets: [
      {
        label: '# of Vulnerabilities',
        data: [criticalVulns, highVulns, mediumVulns, lowVulns],
        backgroundColor: [
          'rgba(255, 99, 132, 0.8)',
          'rgba(255, 159, 64, 0.8)',
          'rgba(255, 205, 86, 0.8)',
          'rgba(75, 192, 192, 0.8)'
        ],
        borderColor: ['rgba(0, 0, 0, 1)', 'rgba(0, 0, 0, 1)', 'rgba(0, 0, 0, 1)', 'rgba(0, 0, 0, 1)'],
        borderWidth: 3
      }
    ]
  };

  const remediationProgressData = {
    labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4'],
    datasets: [
      {
        label: 'Patched',
        data: [0, 0, 0, 0],
        backgroundColor: 'rgba(75, 192, 192, 0.8)',
        borderColor: 'rgba(0, 0, 0, 1)',
        borderWidth: 3
      },
      {
        label: 'Outstanding',
        data: [0, 0, 0, 0],
        backgroundColor: 'rgba(255, 99, 132, 0.8)',
        borderColor: 'rgba(0, 0, 0, 1)',
        borderWidth: 3
      }
    ]
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b-3 border-black shadow-lg">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold uppercase tracking-wider text-black">ZeroTrace Dashboard</h1>
              <p className="text-gray-600 mt-1">Enterprise Security Management Platform</p>
            </div>
            <div className="flex items-center gap-4">
              <button
                onClick={startScan}
                disabled={isScanning}
                className="px-6 py-3 bg-orange-500 text-white border-3 border-black rounded shadow-neubrutalist-lg hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-neubrutalist-md transition-all duration-150 ease-in-out font-bold uppercase tracking-wider disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isScanning ? (
                  <>
                    <div className="animate-spin h-4 w-4 mr-2 border-2 border-current border-t-transparent rounded-full inline-block"></div>
                    SCANNING...
                  </>
                ) : (
                  <>
                    <Activity className="h-4 w-4 mr-2 inline-block" />
                    START SCAN
                  </>
                )}
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto p-6 space-y-6">
        {/* KPI Cards & New Charts */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {/* Total Assets */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-blue-100 rounded border-2 border-black">
                <Server className="h-6 w-6 text-blue-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-black">{totalAssets.toLocaleString()}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Total Assets</div>
              </div>
            </div>
            <div className="flex items-center gap-2 text-sm">
              <TrendingUp className="h-4 w-4 text-green-600" />
              <span className="text-green-600 font-bold">+12%</span>
              <span className="text-gray-500">vs last month</span>
            </div>
          </div>

          {/* Vulnerable Assets */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-red-100 rounded border-2 border-black">
                <AlertTriangle className="h-6 w-6 text-red-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-red-600">{vulnerableAssets}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Vulnerable</div>
              </div>
            </div>
            <div className="flex items-center gap-2 text-sm">
              <TrendingDown className="h-4 w-4 text-red-600" />
              <span className="text-red-600 font-bold">-8%</span>
              <span className="text-gray-500">vs last week</span>
            </div>
          </div>

          {/* Security Score */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-blue-100 rounded border-2 border-black">
                <Shield className="h-6 w-6 text-blue-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-blue-600">
                  {Math.max(0, 100 - (criticalVulns * 10 + highVulns * 5 + mediumVulns * 2 + lowVulns * 1))}
                </div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Security Score</div>
              </div>
            </div>
            <div className="text-sm text-gray-600">{criticalVulns > 0 ? 'Needs Attention' : 'Good'}</div>
          </div>

          {/* Agent Status */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-green-100 rounded border-2 border-black">
                <Shield className="h-6 w-6 text-green-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-green-600">{agentStats.online}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Online</div>
              </div>
            </div>
            <div className="text-sm text-gray-600">{agentStats.offline} offline</div>
          </div>
        </div>

        {/* Additional Metrics Row */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {/* Placeholders for metrics we don't have yet but want layout consistency */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-purple-100 rounded border-2 border-black"><Activity className="h-6 w-6 text-purple-600" /></div>
              <div className="text-right"><div className="text-3xl font-bold text-purple-600">-</div><div className="text-sm text-gray-600 uppercase">Applications</div></div>
            </div>
          </div>
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-orange-100 rounded border-2 border-black"><Zap className="h-6 w-6 text-orange-600" /></div>
              <div className="text-right"><div className="text-3xl font-bold text-orange-600">-</div><div className="text-sm text-gray-600 uppercase">Avg Risk</div></div>
            </div>
          </div>

          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-red-100 rounded border-2 border-black"><AlertTriangle className="h-6 w-6 text-red-600" /></div>
              <div className="text-right"><div className="text-3xl font-bold text-red-600">{totalVulns}</div><div className="text-sm text-gray-600 uppercase">Vulnerabilities</div></div>
            </div>
          </div>

          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-green-100 rounded border-2 border-black"><Clock className="h-6 w-6 text-green-600" /></div>
              <div className="text-right">
                <div className="text-3xl font-bold text-green-600">{overview.last_scan ? new Date(overview.last_scan).toLocaleDateString() : 'N/A'}</div>
                <div className="text-sm text-gray-600 uppercase">Last Scan</div>
              </div>
            </div>
          </div>
        </div>

        {/* Enhanced Dashboard Components */}

        {/* Vulnerability Scanner */}
        <div className="mb-6">
          <VulnerabilityScanner agents={agents || []} />
        </div>

        {/* Software Inventory */}
        <div className="grid grid-cols-1 gap-6">
          <div className="bg-white border-3 border-black rounded-lg shadow-neo-brutal h-[500px] overflow-hidden">
            <AssetAppList agents={agents} />
          </div>
        </div>

        {/* Network Security */}
        <div className="grid grid-cols-1 gap-6">
          <NetworkSecurityDashboard agents={agents} />
        </div>

        {/* Charts and Top Assets */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <h2 className="text-xl font-black text-black uppercase mb-4">Vulnerability Trend</h2>
            <div className="h-64">
              <Charts.Line data={vulnerabilityTrendData} options={{ responsive: true, maintainAspectRatio: false }} />
            </div>
          </div>

          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <h2 className="text-xl font-black text-black uppercase mb-4">By Severity</h2>
            <div className="h-64 flex items-center justify-center">
              <Charts.Doughnut data={severityBreakdownData} options={{ responsive: true, maintainAspectRatio: false }} />
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <TopVulnerableAssets assets={overview.top_vulnerable_assets || []} />
          </div>

          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <h2 className="text-xl font-black text-black uppercase mb-4">Remediation</h2>
            <div className="h-64">
              <Charts.Bar data={remediationProgressData} options={{ responsive: true, maintainAspectRatio: false }} />
            </div>
          </div>
        </div>

        {/* Recent Vulnerabilities */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2">
            <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-xl font-bold uppercase tracking-wider text-black">Recent Vulnerabilities</h2>
              </div>
              {!vulnerabilities || vulnerabilities.length === 0 ? (
                <div className="text-center py-8">
                  <Shield className="h-12 w-12 text-green-500 mx-auto mb-4" />
                  <p className="text-gray-600 font-medium">No vulnerabilities found</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {vulnerabilities.slice(0, 5).map((vuln: any) => (
                    <div key={vuln.id} className="p-4 bg-gray-50 border-2 border-black rounded hover:bg-gray-100 transition-colors">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4">
                          <span className={`px-3 py-1 rounded text-xs font-bold uppercase tracking-wider border-2 ${getSeverityColor(vuln.severity)}`}>
                            {vuln.severity}
                          </span>
                          <div>
                            <p className="font-bold text-black">{vuln.title || vuln.name}</p>
                            <p className="text-sm text-gray-600">{vuln.agent_hostname || vuln.asset || 'Unknown Asset'}</p>
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Recent Activity (Using recent_scans as proxy) */}
          <div className="lg:col-span-1">
            <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg">
              <h2 className="text-xl font-bold uppercase tracking-wider text-black mb-6">Recent Activity</h2>
              {!overview.recent_scans || overview.recent_scans.length === 0 ? (
                <div className="text-center py-8">
                  <Activity className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                  <p className="text-gray-600 font-medium">No recent activity</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {overview.recent_scans.slice(0, 5).map((scan: any) => (
                    <div key={scan.id} className="flex items-start gap-3">
                      <div className="w-3 h-3 bg-orange-500 rounded-full mt-2 flex-shrink-0"></div>
                      <div className="flex-1">
                        <p className="text-sm font-medium text-black">Scan completed on {scan.hostname}</p>
                        <p className="text-xs text-gray-500">{new Date(scan.created_at).toLocaleString()}</p>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;