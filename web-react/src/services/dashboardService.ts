import { api } from './api';

export interface DashboardMetrics {
  assets: {
    total: number;
    online: number;
    offline: number;
    critical: number;
    high: number;
    medium: number;
    low: number;
  };
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
    frameworks: number;
    compliant: number;
    nonCompliant: number;
  };
  performance: {
    avgCpu: number;
    avgMemory: number;
    avgStorage: number;
    networkUtilization: number;
  };
  security: {
    riskScore: number;
    threatLevel: 'low' | 'medium' | 'high' | 'critical';
    lastScan: string;
    nextScan: string;
  };
}

export interface RealTimeMetric {
  id: string;
  name: string;
  value: number;
  unit: string;
  trend: 'up' | 'down' | 'stable';
  status: 'normal' | 'warning' | 'critical';
  timestamp: string;
  history: Array<{
    timestamp: string;
    value: number;
  }>;
}

export interface RealTimeAlert {
  id: string;
  type: 'security' | 'performance' | 'system' | 'network';
  severity: 'critical' | 'high' | 'medium' | 'low';
  title: string;
  description: string;
  asset: string;
  timestamp: string;
  status: 'active' | 'acknowledged' | 'resolved';
  acknowledgedBy?: string;
  resolvedAt?: string;
}

export interface VulnerabilityTrend {
  date: string;
  total: number;
  critical: number;
  high: number;
  medium: number;
  low: number;
  resolved: number;
  new: number;
  exploited: number;
}

export interface RiskScore {
  assetId: string;
  hostname: string;
  score: number;
  factors: Array<{
    name: string;
    weight: number;
    score: number;
    impact: 'high' | 'medium' | 'low';
  }>;
  trend: 'increasing' | 'decreasing' | 'stable';
  lastUpdated: string;
}

export interface VulnerabilityPattern {
  type: string;
  frequency: number;
  severity: 'critical' | 'high' | 'medium' | 'low';
  affectedAssets: number;
  commonCauses: string[];
  mitigation: string[];
}

class DashboardService {
  // Get comprehensive dashboard metrics
  async getDashboardMetrics(): Promise<DashboardMetrics> {
    try {
      const [agentsResponse, vulnerabilitiesResponse] = await Promise.all([
        api.get('/api/agents/'),
        api.get('/api/vulnerabilities/')
      ]);

      const agents = agentsResponse.data.data || [];
      const vulnerabilities = vulnerabilitiesResponse.data.data || [];

      // Calculate asset metrics
      const totalAssets = agents.length;
      const onlineAssets = agents.filter((agent: any) => agent.status === 'online').length;
      const offlineAssets = totalAssets - onlineAssets;

      // Calculate vulnerability metrics
      const totalVulns = vulnerabilities.length;
      const criticalVulns = vulnerabilities.filter((v: any) => v.severity === 'critical').length;
      const highVulns = vulnerabilities.filter((v: any) => v.severity === 'high').length;
      const mediumVulns = vulnerabilities.filter((v: any) => v.severity === 'medium').length;
      const lowVulns = vulnerabilities.filter((v: any) => v.severity === 'low').length;

      // Calculate performance metrics from agents
      let totalCpu = 0;
      let totalMemory = 0;
      let totalStorage = 0;
      let networkUtilization = 0;

      agents.forEach((agent: any) => {
        if (agent.metadata?.cpu_usage) totalCpu += agent.metadata.cpu_usage;
        if (agent.metadata?.memory_usage) totalMemory += agent.metadata.memory_usage;
        if (agent.metadata?.storage_usage) totalStorage += agent.metadata.storage_usage;
        if (agent.metadata?.network_utilization) networkUtilization += agent.metadata.network_utilization;
      });

      const avgCpu = totalAssets > 0 ? totalCpu / totalAssets : 0;
      const avgMemory = totalAssets > 0 ? totalMemory / totalAssets : 0;
      const avgStorage = totalAssets > 0 ? totalStorage / totalAssets : 0;

      // Calculate security metrics
      let totalRiskScore = 0;
      agents.forEach((agent: any) => {
        if (agent.risk_score) totalRiskScore += agent.risk_score;
      });
      const avgRiskScore = totalAssets > 0 ? totalRiskScore / totalAssets : 0;

      const threatLevel = avgRiskScore >= 8 ? 'critical' :
        avgRiskScore >= 6 ? 'high' :
          avgRiskScore >= 4 ? 'medium' : 'low';

      // Calculate compliance metrics (approximate from risk)
      const compliantAssets = agents.filter((a: any) => (a.risk_score || 0) < 4).length; // Low risk = compliant
      const nonCompliantAssets = totalAssets - compliantAssets;
      const complianceScore = totalAssets > 0 ? Math.round((compliantAssets / totalAssets) * 100) : 100;

      return {
        assets: {
          total: totalAssets,
          online: onlineAssets,
          offline: offlineAssets,
          critical: criticalVulns,
          high: highVulns,
          medium: mediumVulns,
          low: lowVulns
        },
        vulnerabilities: {
          total: totalVulns,
          critical: criticalVulns,
          high: highVulns,
          medium: mediumVulns,
          low: lowVulns,
          trend: 'stable' // Trends require historical data snapshoting
        },
        compliance: {
          score: complianceScore,
          frameworks: 3, // Placeholder until compliance frameworks are fully implemented
          compliant: compliantAssets,
          nonCompliant: nonCompliantAssets
        },
        performance: {
          avgCpu,
          avgMemory,
          avgStorage,
          networkUtilization: totalAssets > 0 ? networkUtilization / totalAssets : 0
        },
        security: {
          riskScore: avgRiskScore,
          threatLevel,
          lastScan: agents.length > 0 ? (agents[0] as any).last_scan_time || new Date().toISOString() : new Date().toISOString(),
          nextScan: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString() // Next day
        }
      };
    } catch (error) {
      console.error('Failed to fetch dashboard metrics:', error);
      throw error;
    }
  }

  // Get real-time metrics
  async getRealTimeMetrics(): Promise<RealTimeMetric[]> {
    try {
      const response = await api.get('/api/agents/');
      const agents = response.data.data || [];

      const metrics: RealTimeMetric[] = [];

      // CPU Usage metric
      let totalCpu = 0;
      let cpuCount = 0;
      agents.forEach((agent: any) => {
        if (agent.metadata?.cpu_usage) {
          totalCpu += agent.metadata.cpu_usage;
          cpuCount++;
        }
      });

      if (cpuCount > 0) {
        metrics.push({
          id: 'cpu-usage',
          name: 'CPU Usage',
          value: totalCpu / cpuCount,
          unit: '%',
          trend: 'stable',
          status: totalCpu / cpuCount > 80 ? 'critical' : totalCpu / cpuCount > 60 ? 'warning' : 'normal',
          timestamp: new Date().toISOString(),
          history: [] // TODO: Add historical data
        });
      }

      // Memory Usage metric
      let totalMemory = 0;
      let memoryCount = 0;
      agents.forEach((agent: any) => {
        if (agent.metadata?.memory_usage) {
          totalMemory += agent.metadata.memory_usage;
          memoryCount++;
        }
      });

      if (memoryCount > 0) {
        metrics.push({
          id: 'memory-usage',
          name: 'Memory Usage',
          value: totalMemory / memoryCount,
          unit: '%',
          trend: 'stable',
          status: totalMemory / memoryCount > 85 ? 'critical' : totalMemory / memoryCount > 70 ? 'warning' : 'normal',
          timestamp: new Date().toISOString(),
          history: []
        });
      }

      // Vulnerability count metric
      const vulnResponse = await api.get('/api/vulnerabilities/');
      const vulnerabilities = vulnResponse.data.data || [];

      metrics.push({
        id: 'vulnerabilities',
        name: 'Active Vulnerabilities',
        value: vulnerabilities.length,
        unit: 'count',
        trend: 'stable',
        status: vulnerabilities.length > 50 ? 'critical' : vulnerabilities.length > 20 ? 'warning' : 'normal',
        timestamp: new Date().toISOString(),
        history: []
      });

      return metrics;
    } catch (error) {
      console.error('Failed to fetch real-time metrics:', error);
      return [];
    }
  }

  // Get real-time alerts
  async getRealTimeAlerts(): Promise<RealTimeAlert[]> {
    try {
      const [agentsResponse, vulnerabilitiesResponse] = await Promise.all([
        api.get('/api/agents/'),
        api.get('/api/vulnerabilities/')
      ]);

      const agents = agentsResponse.data.data || [];
      const vulnerabilities = vulnerabilitiesResponse.data.data || [];

      const alerts: RealTimeAlert[] = [];

      // Create alerts for critical vulnerabilities
      vulnerabilities.filter((v: any) => v.severity === 'critical').forEach((vuln: any) => {
        alerts.push({
          id: `vuln-${vuln.id}`,
          type: 'security',
          severity: 'critical',
          title: `Critical Vulnerability: ${vuln.title}`,
          description: vuln.description || 'Critical security vulnerability detected',
          asset: vuln.asset_name || 'Unknown',
          timestamp: vuln.created_at || new Date().toISOString(),
          status: 'active'
        });
      });

      // Create alerts for high CPU usage
      agents.forEach((agent: any) => {
        if (agent.metadata?.cpu_usage && agent.metadata.cpu_usage > 90) {
          alerts.push({
            id: `cpu-${agent.id}`,
            type: 'performance',
            severity: 'high',
            title: `High CPU Usage on ${agent.hostname}`,
            description: `CPU usage is at ${agent.metadata.cpu_usage}%`,
            asset: agent.hostname || agent.id,
            timestamp: new Date().toISOString(),
            status: 'active'
          });
        }
      });

      // Create alerts for high memory usage
      agents.forEach((agent: any) => {
        if (agent.metadata?.memory_usage && agent.metadata.memory_usage > 90) {
          alerts.push({
            id: `memory-${agent.id}`,
            type: 'performance',
            severity: 'high',
            title: `High Memory Usage on ${agent.hostname}`,
            description: `Memory usage is at ${agent.metadata.memory_usage}%`,
            asset: agent.hostname || agent.id,
            timestamp: new Date().toISOString(),
            status: 'active'
          });
        }
      });

      return alerts.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime());
    } catch (error) {
      console.error('Failed to fetch real-time alerts:', error);
      return [];
    }
  }

  // Get vulnerability trends
  async getVulnerabilityTrends(): Promise<VulnerabilityTrend[]> {
    try {
      const response = await api.get('/api/vulnerabilities/');
      const vulnerabilities = response.data.data || [];

      // Group vulnerabilities by date (last 7 days)
      const trends: VulnerabilityTrend[] = [];
      const today = new Date();

      for (let i = 6; i >= 0; i--) {
        const date = new Date(today);
        date.setDate(date.getDate() - i);
        const dateStr = date.toISOString().split('T')[0];

        const dayVulns = vulnerabilities.filter((v: any) => {
          const vulnDate = new Date(v.created_at || v.updated_at);
          return vulnDate.toISOString().split('T')[0] === dateStr;
        });

        // Filter for this day's activity
        const resolvedCount = vulnerabilities.filter((v: any) => {
          const upDate = new Date(v.updated_at);
          return (v.status === 'resolved' || v.status === 'fixed') &&
            upDate.toISOString().split('T')[0] === dateStr;
        }).length;

        const exploitedCount = dayVulns.filter((v: any) => v.exploit_available).length;

        trends.push({
          date: dateStr,
          total: dayVulns.length,
          critical: dayVulns.filter((v: any) => v.severity === 'critical').length,
          high: dayVulns.filter((v: any) => v.severity === 'high').length,
          medium: dayVulns.filter((v: any) => v.severity === 'medium').length,
          low: dayVulns.filter((v: any) => v.severity === 'low').length,
          resolved: resolvedCount,
          new: dayVulns.length,
          exploited: exploitedCount
        });
      }

      return trends;
    } catch (error) {
      console.error('Failed to fetch vulnerability trends:', error);
      return [];
    }
  }

  // Get risk scores for assets
  async getRiskScores(): Promise<RiskScore[]> {
    try {
      const response = await api.get('/api/agents/');
      const agents = response.data.data || [];

      return agents.map((agent: any) => ({
        assetId: agent.id,
        hostname: agent.hostname || agent.id,
        score: agent.risk_score || 0,
        factors: [
          {
            name: 'Vulnerability Count',
            weight: 0.4,
            score: Math.min((agent.metadata?.vulnerabilities?.length || 0) * 2, 10),
            impact: 'high'
          },
          {
            name: 'System Health',
            weight: 0.3,
            score: agent.metadata?.cpu_usage > 80 ? 8 : agent.metadata?.memory_usage > 80 ? 7 : 3,
            impact: 'medium'
          },
          {
            name: 'Network Exposure',
            weight: 0.3,
            score: agent.metadata?.network_utilization > 90 ? 6 : 2,
            impact: 'medium'
          }
        ],
        trend: 'stable', // TODO: Calculate actual trend
        lastUpdated: agent.updated_at || new Date().toISOString()
      }));
    } catch (error) {
      console.error('Failed to fetch risk scores:', error);
      return [];
    }
  }

  // Get vulnerability patterns
  async getVulnerabilityPatterns(): Promise<VulnerabilityPattern[]> {
    try {
      const response = await api.get('/api/vulnerabilities/');
      const vulnerabilities = response.data.data || [];

      // Group vulnerabilities by type/pattern
      const patterns: { [key: string]: VulnerabilityPattern } = {};

      vulnerabilities.forEach((vuln: any) => {
        const type = vuln.type || 'Unknown';
        if (!patterns[type]) {
          patterns[type] = {
            type,
            frequency: 0,
            severity: vuln.severity,
            affectedAssets: 0,
            commonCauses: [],
            mitigation: []
          };
        }
        patterns[type].frequency++;
      });

      return Object.values(patterns).sort((a, b) => b.frequency - a.frequency);
    } catch (error) {
      console.error('Failed to fetch vulnerability patterns:', error);
      return [];
    }
  }
}

export const dashboardService = new DashboardService();