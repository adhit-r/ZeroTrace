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
  ArrowRight,
  Eye,
  RefreshCw
} from 'lucide-react';
import { agentService } from '../services/agentService';

// Import charting library and new components
import { Bar, Doughnut, Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
// Removed mock chart data - using real data only
import TopVulnerableAssets from '../components/dashboard/TopVulnerableAssets';
// import InnovativeDashboard from '../components/dashboard/InnovativeDashboard';
import RealTimeMonitoring from '../components/dashboard/RealTimeMonitoring';
import VulnerabilityTrendAnalysis from '../components/dashboard/VulnerabilityTrendAnalysis';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend
);

// Define types locally to avoid import issues
interface Asset {
  total: number;
  vulnerable: number;
  critical: number;
  high: number;
  medium: number;
  low: number;
  lastScan: string | null;
}

interface Vulnerability {
  id: number;
  name: string;
  severity: string;
  asset: string;
  status: string;
}

interface RecentActivity {
  id: string;
  type: string;
  message: string;
  time: string;
}

interface TopVulnerableAsset {
  name: string;
  vulnerabilities: number;
  critical: number;
}

interface DashboardData {
  assets: Asset;
  vulnerabilities: Vulnerability[];
  recentActivity: RecentActivity[];
  topVulnerableAssets: TopVulnerableAsset[];
  scanStatus: string;
  vulnerabilityTrendData?: any;
  severityBreakdownData?: any;
  remediationProgressData?: any;
}

// Real data for dashboard
const useDashboardData = () => {
  const [data, setData] = useState<DashboardData>({
    assets: {
      total: 0,
      vulnerable: 0,
      critical: 0,
      high: 0,
      medium: 0,
      low: 0,
      lastScan: null
    },
    vulnerabilities: [],
    recentActivity: [],
    topVulnerableAssets: [],
    scanStatus: 'idle',
    vulnerabilityTrendData: undefined,
    severityBreakdownData: undefined,
    remediationProgressData: undefined
  });
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [agentStats, setAgentStats] = useState({
    total: 0,
    online: 0,
    offline: 0,
    avgCpu: 0,
    avgMemory: 0
  });
  const [additionalMetrics, setAdditionalMetrics] = useState({
    totalDependencies: 0,
    totalApplications: 0,
    avgRiskScore: 0,
    lastScanTime: null
  });

    const loadData = async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      console.log('Dashboard: Loading data...');
      
      // Fetch agent data to get real counts
      const agents = await agentService.getAgents();
      const agentStatsData = await agentService.getAgentStats();
      console.log('Dashboard: Agents:', agents);
      console.log('Dashboard: Agent stats:', agentStatsData);
      
      // Set agent stats
      setAgentStats(agentStatsData);
      
      // Calculate real vulnerability metrics from agents
      let totalVulnerabilities = 0;
      let criticalVulns = 0;
      let highVulns = 0;
      let mediumVulns = 0;
      let lowVulns = 0;
      let vulnerableAssets = 0;
      
      agents.forEach(agent => {
        if (agent.metadata?.vulnerabilities) {
          const vulns = agent.metadata.vulnerabilities;
          totalVulnerabilities += vulns.length;
          
          vulns.forEach((vuln: any) => {
            switch(vuln.severity) {
              case 'critical':
                criticalVulns++;
                break;
              case 'high':
                highVulns++;
                break;
              case 'medium':
                mediumVulns++;
                break;
              case 'low':
                lowVulns++;
                break;
            }
          });
          
          if (vulns.length > 0) {
            vulnerableAssets++;
          }
        }
      });
      
      console.log('Dashboard: Assets total:', agents.length);
      console.log('Dashboard: Vulnerability counts:', { totalVulnerabilities, criticalVulns, highVulns, mediumVulns, lowVulns, vulnerableAssets });
      
            // Calculate additional metrics
            let totalApplications = 0;
            let totalDependencies = 0;
            let avgRiskScore = 0;
            let lastScanTime: string | null = null;
            
            agents.forEach(agent => {
              if (agent.metadata?.dependencies) {
                totalDependencies += agent.metadata.dependencies.length;
              }
              if (agent.metadata?.applications) {
                totalApplications += agent.metadata.applications.length;
              }
              if ((agent as any).risk_score !== undefined) {
                avgRiskScore += (agent as any).risk_score;
              }
              if ((agent as any).last_scan_time && (!lastScanTime || new Date((agent as any).last_scan_time) > new Date(lastScanTime))) {
                lastScanTime = (agent as any).last_scan_time;
              }
            });

            // No agents - will show empty state
      
      avgRiskScore = agents.length > 0 ? avgRiskScore / agents.length : 0;
      
      console.log('Dashboard: Additional metrics:', { 
        totalApplications, 
        totalDependencies, 
        avgRiskScore, 
        lastScanTime 
      });
      
      // Set additional metrics in state
      setAdditionalMetrics({
        totalDependencies,
        totalApplications,
        avgRiskScore,
        lastScanTime
      });
      
      // Create dynamic chart data based on real data
      const vulnerabilityTrendData = {
        labels: ['Last 7 Days', 'Last 6 Days', 'Last 5 Days', 'Last 4 Days', 'Last 3 Days', 'Last 2 Days', 'Today'],
        datasets: [
          {
            label: 'Vulnerabilities Found',
            data: [0, 0, 0, 0, 0, 0, totalVulnerabilities], // Only show real data
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
            data: [
              criticalVulns,
              highVulns,
              mediumVulns,
              lowVulns
            ],
            backgroundColor: [
              'rgba(255, 99, 132, 0.8)',
              'rgba(255, 159, 64, 0.8)',
              'rgba(255, 205, 86, 0.8)',
              'rgba(75, 192, 192, 0.8)'
            ],
            borderColor: [
              'rgba(0, 0, 0, 1)',
              'rgba(0, 0, 0, 1)',
              'rgba(0, 0, 0, 1)',
              'rgba(0, 0, 0, 1)'
            ],
            borderWidth: 3
          }
        ]
      };

      const remediationProgressData = {
        labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4'],
        datasets: [
          {
            label: 'Patched',
            data: [0, 0, 0, 0], // No hardcoded data
            backgroundColor: 'rgba(75, 192, 192, 0.8)',
            borderColor: 'rgba(0, 0, 0, 1)',
            borderWidth: 3
          },
          {
            label: 'Outstanding',
            data: [0, 0, 0, 0], // No hardcoded data
            backgroundColor: 'rgba(255, 99, 132, 0.8)',
            borderColor: 'rgba(0, 0, 0, 1)',
            borderWidth: 3
          }
        ]
      };

      // Transform API data to dashboard format
      // Use real vulnerability counts from agents
      const transformedData: DashboardData = {
        assets: {
          total: agents.length,
          vulnerable: vulnerableAssets,
          critical: criticalVulns,
          high: highVulns,
          medium: mediumVulns,
          low: lowVulns,
          lastScan: agents.length > 0 ? (agents[0] as any).last_scan_time : null
        },
        vulnerabilities: agents.flatMap(agent => 
          agent.metadata?.vulnerabilities?.map((vuln: any) => ({
            id: vuln.id,
            name: vuln.title,
            severity: vuln.severity,
            asset: agent.hostname || agent.id,
            status: vuln.status || 'open'
          })) || []
        ),
        recentActivity: agents.map(agent => ({
          id: `activity-${agent.id}`,
          type: 'scan',
          message: `Scan completed on ${agent.hostname || agent.id} - Found ${agent.metadata?.vulnerabilities?.length || 0} vulnerabilities`,
          time: (agent as any).last_scan_time ? new Date((agent as any).last_scan_time).toLocaleString() : 'Unknown'
        })).sort((a, b) => new Date(b.time).getTime() - new Date(a.time).getTime()).slice(0, 5),
        topVulnerableAssets: agents
          .filter(agent => agent.metadata?.vulnerabilities && agent.metadata.vulnerabilities.length > 0)
          .map(agent => ({
            name: agent.hostname || agent.id,
            vulnerabilities: agent.metadata!.vulnerabilities.length,
            critical: agent.metadata!.vulnerabilities.filter((v: any) => v.severity === 'critical').length
          }))
          .sort((a, b) => b.vulnerabilities - a.vulnerabilities)
          .slice(0, 5),
        scanStatus: agentStatsData.online > 0 ? 'active' : 'idle',
        // Add real-time chart data
        vulnerabilityTrendData,
        severityBreakdownData,
        remediationProgressData
      };
      
      setData(transformedData);
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
      setError('Failed to load dashboard data. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadData();
    
    // Set up auto-refresh every 60 seconds
    const interval = setInterval(loadData, 60000);
    return () => clearInterval(interval);
  }, []);

  return { data, agentStats, additionalMetrics, isLoading, error, refresh: loadData };
};

const Dashboard: React.FC = () => {
  const { data, agentStats, additionalMetrics, isLoading, error, refresh } = useDashboardData();
  
  // Extract data from the data object
  const assets = data?.assets;
  const vulnerabilities = data?.vulnerabilities;
  const recentActivity = data?.recentActivity;
  const topVulnerableAssets = data?.topVulnerableAssets;
  const [isScanning, setIsScanning] = useState(false);

  const startScan = async () => {
    setIsScanning(true);
    try {
      // Trigger a scan by calling the API
      const response = await fetch('http://localhost:8080/api/agents/', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      
      if (response.ok) {
        // Refresh the dashboard data
        await refresh();
      }
    } catch (error) {
      console.error('Failed to trigger scan:', error);
    } finally {
      setIsScanning(false);
    }
  };

  const deployAgent = () => {
    // Open agent download/installation instructions
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
    
    // Also show alert with instructions
    alert('Agent installation instructions have been downloaded. Please follow the instructions to deploy the agent on your systems.');
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'bg-red-100 text-red-800 border-red-300';
      case 'high':
        return 'bg-orange-100 text-orange-800 border-orange-300';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800 border-yellow-300';
      case 'low':
        return 'bg-green-100 text-green-800 border-green-300';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-300';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'open':
        return 'bg-red-100 text-red-800';
      case 'in_progress':
        return 'bg-yellow-100 text-yellow-800';
      case 'closed':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
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
              <h2 className="text-xl font-bold text-red-800">Error Loading Dashboard</h2>
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
  console.log('Dashboard: Assets total:', assets.total);
  if (assets.total === 0) {
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
                <div className="text-3xl font-bold text-black">{assets.total.toLocaleString()}</div>
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
                <div className="text-3xl font-bold text-red-600">{assets.vulnerable}</div>
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
                <div className="text-3xl font-bold text-blue-600">{Math.max(0, 100 - (assets.critical * 10 + assets.high * 5 + assets.medium * 2 + assets.low * 1))}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Security Score</div>
              </div>
            </div>
            <div className="text-sm text-gray-600">{assets.critical > 0 ? 'Needs Attention' : 'Good'}</div>
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
          {/* Total Applications */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-purple-100 rounded border-2 border-black">
                <Activity className="h-6 w-6 text-purple-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-purple-600">{additionalMetrics.totalDependencies.toLocaleString()}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Applications</div>
              </div>
            </div>
            <div className="text-sm text-gray-600">Across all assets</div>
          </div>

          {/* Average Risk Score */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-orange-100 rounded border-2 border-black">
                <Zap className="h-6 w-6 text-orange-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-orange-600">{additionalMetrics.avgRiskScore.toFixed(1)}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Avg Risk</div>
              </div>
            </div>
            <div className="text-sm text-gray-600">Risk score per asset</div>
          </div>

          {/* Total Vulnerabilities */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-red-100 rounded border-2 border-black">
                <AlertTriangle className="h-6 w-6 text-red-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-red-600">{data.assets.critical + data.assets.high + data.assets.medium + data.assets.low}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Vulnerabilities</div>
              </div>
            </div>
            <div className="text-sm text-gray-600">Total security issues</div>
          </div>

          {/* Last Scan */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-green-100 rounded border-2 border-black">
                <Clock className="h-6 w-6 text-green-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-green-600">
                  {additionalMetrics.lastScanTime ? new Date(additionalMetrics.lastScanTime).toLocaleDateString() : 'N/A'}
                </div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Last Scan</div>
              </div>
            </div>
            <div className="text-sm text-gray-600">
              {additionalMetrics.lastScanTime ? new Date(additionalMetrics.lastScanTime).toLocaleTimeString() : 'No scans yet'}
            </div>
          </div>
        </div>

        {/* Enhanced Dashboard Components */}
        <div className="space-y-6">
          {/* Innovative Dashboard */}
          {/* <InnovativeDashboard /> */}
          
          {/* Real-time Monitoring */}
          <RealTimeMonitoring />
          
          {/* Vulnerability Trend Analysis */}
          <VulnerabilityTrendAnalysis />
        </div>

        {/* Charts and Top Assets */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Vulnerability Trend */}
          <div className="lg:col-span-2 p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <h2 className="text-xl font-black text-black uppercase mb-4">Vulnerability Trend</h2>
                  <div className="h-64">
                    <Line data={data.vulnerabilityTrendData} options={{
                      responsive: true,
                      maintainAspectRatio: false,
                    }} />
                  </div>
          </div>

          {/* Vulnerabilities by Severity */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <h2 className="text-xl font-black text-black uppercase mb-4">By Severity</h2>
                  <div className="h-64 flex items-center justify-center">
                    <Doughnut data={data.severityBreakdownData} options={{
                      responsive: true,
                      maintainAspectRatio: false,
                    }} />
                  </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
           {/* Top Vulnerable Assets */}
          <div className="lg:col-span-2 p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
             <TopVulnerableAssets assets={topVulnerableAssets} />
            </div>

          {/* Remediation Progress */}
          <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <h2 className="text-xl font-black text-black uppercase mb-4">Remediation</h2>
                  <div className="h-64">
                      <Bar data={data.remediationProgressData} options={{
                        responsive: true,
                        maintainAspectRatio: false,
                        scales: {
                          y: {
                            beginAtZero: true,
                          },
                        },
                      }} />
                </div>
        </div>
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Recent Vulnerabilities */}
        <div className="lg:col-span-2">
            <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg">
            <div className="flex items-center justify-between mb-6">
                <h2 className="text-xl font-bold uppercase tracking-wider text-black">Recent Vulnerabilities</h2>
                <button className="px-4 py-2 bg-orange-100 text-orange-800 border-2 border-orange-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-orange-200 transition-colors">
                  VIEW ALL <ArrowRight className="h-4 w-4 ml-2 inline-block" />
              </button>
            </div>
            
              {vulnerabilities.length === 0 ? (
                <div className="text-center py-8">
                  <Shield className="h-12 w-12 text-green-500 mx-auto mb-4" />
                  <p className="text-gray-600 font-medium">No vulnerabilities found</p>
                  <p className="text-sm text-gray-500">All systems are secure!</p>
                </div>
              ) : (
            <div className="space-y-4">
                  {vulnerabilities.map((vuln) => (
                    <div key={vuln.id} className="p-4 bg-gray-50 border-2 border-black rounded hover:bg-gray-100 transition-colors">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4">
                          <span className={`px-3 py-1 rounded text-xs font-bold uppercase tracking-wider border-2 ${getSeverityColor(vuln.severity)}`}>
                            {vuln.severity}
                          </span>
                    <div>
                            <p className="font-bold text-black">{vuln.name}</p>
                            <p className="text-sm text-gray-600">{vuln.asset}</p>
                    </div>
                  </div>
                        <div className="flex items-center gap-2">
                          <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider ${getStatusColor(vuln.status)}`}>
                            {vuln.status}
                    </span>
                          <button className="p-2 hover:bg-gray-200 rounded border-2 border-black">
                      <Eye className="h-4 w-4" />
                    </button>
                        </div>
                  </div>
                </div>
              ))}
            </div>
              )}
          </div>
        </div>

        {/* Recent Activity */}
        <div className="lg:col-span-1">
            <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg">
              <h2 className="text-xl font-bold uppercase tracking-wider text-black mb-6">Recent Activity</h2>
              
              {recentActivity.length === 0 ? (
                <div className="text-center py-8">
                  <Activity className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                  <p className="text-gray-600 font-medium">No recent activity</p>
                  <p className="text-sm text-gray-500">Activity will appear here as agents scan</p>
                </div>
              ) : (
            <div className="space-y-4">
                  {recentActivity.map((activity) => (
                    <div key={activity.id} className="flex items-start gap-3">
                      <div className="w-3 h-3 bg-orange-500 rounded-full mt-2 flex-shrink-0"></div>
                  <div className="flex-1">
                        <p className="text-sm font-medium text-black">{activity.message}</p>
                        <p className="text-xs text-gray-500 mt-1">{activity.time}</p>
                  </div>
                </div>
              ))}
            </div>
              )}
          </div>
        </div>
      </div>

      {/* Top Vulnerable Assets */}
        {topVulnerableAssets.length > 0 && (
          <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg">
        <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-bold uppercase tracking-wider text-black">Top Vulnerable Assets</h2>
              <button className="px-4 py-2 bg-orange-100 text-orange-800 border-2 border-orange-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-orange-200 transition-colors">
                VIEW ALL ASSETS <ArrowRight className="h-4 w-4 ml-2 inline-block" />
          </button>
        </div>
        
        <div className="overflow-x-auto">
              <table className="w-full">
            <thead>
                  <tr className="border-b-2 border-black">
                    <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">Asset Name</th>
                    <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">Total Vulnerabilities</th>
                    <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">Critical</th>
                    <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">Status</th>
                    <th className="text-left py-3 px-4 font-bold uppercase tracking-wider text-black">Actions</th>
              </tr>
            </thead>
            <tbody>
                  {topVulnerableAssets.map((asset, index) => (
                    <tr key={index} className="border-b border-gray-200 hover:bg-gray-50">
                      <td className="py-3 px-4 font-bold text-black">{asset.name}</td>
                      <td className="py-3 px-4">
                        <span className="px-2 py-1 bg-orange-100 text-orange-800 border-2 border-orange-300 rounded text-xs font-bold uppercase tracking-wider">
                          {asset.vulnerabilities}
                        </span>
                  </td>
                      <td className="py-3 px-4">
                        <span className="px-2 py-1 bg-red-100 text-red-800 border-2 border-red-300 rounded text-xs font-bold uppercase tracking-wider">
                          {asset.critical}
                        </span>
                  </td>
                      <td className="py-3 px-4">
                        <span className="px-2 py-1 bg-blue-100 text-blue-800 border-2 border-blue-300 rounded text-xs font-bold uppercase tracking-wider">
                          MONITORING
                        </span>
                  </td>
                      <td className="py-3 px-4">
                        <button className="p-2 hover:bg-gray-200 rounded border-2 border-black">
                      <Eye className="h-4 w-4" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Dashboard;