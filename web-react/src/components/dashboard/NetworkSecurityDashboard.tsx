import React, { useMemo } from 'react';
import {
  Network,
  Shield,
  Activity,
  AlertTriangle,
  CheckCircle,
  Server,
  Wifi,
  Globe,
  Eye
} from 'lucide-react';
import type { Agent } from '../../services/agentService';

interface NetworkSecurityDashboardProps {
  agents: Agent[];
}

interface NetworkFinding {
  id: string;
  host: string;
  port: number;
  protocol: string;
  finding_type: string;
  severity: string;
  description: string;
  service_name?: string;
  service_version?: string;
  status: string;
  discovered_at: string;
}

const NetworkSecurityDashboard: React.FC<NetworkSecurityDashboardProps> = ({ agents }) => {

  // Aggregate network data from all agents
  const { networkFindings, uniqueHosts, openPorts, criticalFindings } = useMemo(() => {
    let allFindings: NetworkFinding[] = [];
    const hostSet = new Set<string>();
    let critical = 0;

    agents.forEach(agent => {
      // Check for network scan results in metadata
      const scanResult = agent.metadata?.network_scan_result;
      if (scanResult?.network_findings && Array.isArray(scanResult.network_findings)) {
        scanResult.network_findings.forEach((finding: NetworkFinding) => {
          allFindings.push({ ...finding, agent_name: agent.name } as any);
          hostSet.add(finding.host);
          if (finding.severity === 'critical' || finding.severity === 'high') {
            critical++;
          }
        });
      }
    });

    return {
      networkFindings: allFindings,
      uniqueHosts: Array.from(hostSet),
      openPorts: allFindings.filter(f => f.finding_type === 'port').length,
      criticalFindings: critical
    };
  }, [agents]);

  // Group findings by host
  const hostGroups = useMemo(() => {
    const groups: Record<string, NetworkFinding[]> = {};
    networkFindings.forEach(finding => {
      if (!groups[finding.host]) {
        groups[finding.host] = [];
      }
      groups[finding.host].push(finding);
    });
    return groups;
  }, [networkFindings]);

  const getSeverityColor = (severity: string) => {
    switch (severity?.toLowerCase()) {
      case 'critical': return 'text-red-600 bg-red-100 border-red-300';
      case 'high': return 'text-orange-600 bg-orange-100 border-orange-300';
      case 'medium': return 'text-yellow-600 bg-yellow-100 border-yellow-300';
      case 'low': return 'text-green-600 bg-green-100 border-green-300';
      default: return 'text-blue-600 bg-blue-100 border-blue-300';
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="p-6 bg-gradient-to-r from-slate-900 to-slate-800 border-3 border-black rounded-lg shadow-neo-brutal text-white">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-black uppercase mb-1">Network Security</h2>
            <p className="text-slate-300">Real-time network monitoring and threat detection</p>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-green-500 rounded-full animate-pulse"></div>
            <span className="text-sm text-slate-300">Live Scanning</span>
          </div>
        </div>
      </div>

      {/* Network Metrics Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Discovered Hosts */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
          <div className="flex items-center justify-between mb-4">
            <div className="p-3 bg-blue-100 rounded-lg border-2 border-black">
              <Globe className="h-6 w-6 text-blue-600" />
            </div>
            <div className="text-right">
              <div className="text-3xl font-black text-blue-600">{uniqueHosts.length}</div>
              <div className="text-xs text-gray-500 uppercase tracking-wider">Found</div>
            </div>
          </div>
          <h3 className="text-lg font-black text-black uppercase">Discovered Hosts</h3>
          <p className="text-sm text-gray-500 mt-1">Unique IPs on network</p>
        </div>

        {/* Open Ports */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
          <div className="flex items-center justify-between mb-4">
            <div className="p-3 bg-orange-100 rounded-lg border-2 border-black">
              <Activity className="h-6 w-6 text-orange-600" />
            </div>
            <div className="text-right">
              <div className="text-3xl font-black text-orange-600">{openPorts}</div>
              <div className="text-xs text-gray-500 uppercase tracking-wider">Detected</div>
            </div>
          </div>
          <h3 className="text-lg font-black text-black uppercase">Open Ports</h3>
          <p className="text-sm text-gray-500 mt-1">Active services found</p>
        </div>

        {/* Scanning Agents */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
          <div className="flex items-center justify-between mb-4">
            <div className="p-3 bg-green-100 rounded-lg border-2 border-black">
              <Eye className="h-6 w-6 text-green-600" />
            </div>
            <div className="text-right">
              <div className="text-3xl font-black text-green-600">
                {agents.filter(a => a.status === 'online').length}
              </div>
              <div className="text-xs text-gray-500 uppercase tracking-wider">Active</div>
            </div>
          </div>
          <h3 className="text-lg font-black text-black uppercase">Scanning Agents</h3>
          <p className="text-sm text-gray-500 mt-1">Running network scans</p>
        </div>

        {/* Security Issues */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
          <div className="flex items-center justify-between mb-4">
            <div className={`p-3 rounded-lg border-2 border-black ${criticalFindings > 0 ? 'bg-red-100' : 'bg-green-100'}`}>
              <AlertTriangle className={`h-6 w-6 ${criticalFindings > 0 ? 'text-red-600' : 'text-green-600'}`} />
            </div>
            <div className="text-right">
              <div className={`text-3xl font-black ${criticalFindings > 0 ? 'text-red-600' : 'text-green-600'}`}>
                {criticalFindings}
              </div>
              <div className="text-xs text-gray-500 uppercase tracking-wider">Issues</div>
            </div>
          </div>
          <h3 className="text-lg font-black text-black uppercase">Security Issues</h3>
          <p className="text-sm text-gray-500 mt-1">Critical/High findings</p>
        </div>
      </div>

      {/* Network Overview */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Discovered Hosts List */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <h3 className="text-xl font-black text-black uppercase mb-4 flex items-center gap-2">
            <Network className="h-5 w-5" />
            Discovered Hosts
          </h3>

          {Object.keys(hostGroups).length === 0 ? (
            <div className="h-64 flex flex-col items-center justify-center bg-gray-50 border-2 border-gray-300 rounded-lg">
              <Wifi className="h-12 w-12 text-gray-400 mx-auto mb-3" />
              <p className="text-gray-600 font-bold">No Hosts Discovered</p>
              <p className="text-sm text-gray-500 mt-1">Network scan in progress or not started</p>
            </div>
          ) : (
            <div className="overflow-y-auto max-h-80 space-y-2">
              {Object.entries(hostGroups).map(([host, findings]) => (
                <div key={host} className="p-4 bg-gray-50 border-2 border-gray-200 rounded-lg hover:bg-gray-100 transition-colors">
                  <div className="flex justify-between items-start">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-blue-100 rounded border border-blue-300">
                        <Server className="w-4 h-4 text-blue-600" />
                      </div>
                      <div>
                        <div className="font-bold text-gray-800">{host}</div>
                        <div className="text-xs text-gray-500">
                          {findings.length} port{findings.length !== 1 ? 's' : ''} open
                        </div>
                      </div>
                    </div>
                    <div className="flex flex-wrap gap-1 justify-end max-w-[50%]">
                      {findings.slice(0, 4).map((f, idx) => (
                        <span key={idx} className="text-xs font-mono px-2 py-0.5 bg-gray-200 rounded text-gray-700">
                          {f.port}/{f.protocol}
                        </span>
                      ))}
                      {findings.length > 4 && (
                        <span className="text-xs px-2 py-0.5 bg-gray-300 rounded text-gray-600">
                          +{findings.length - 4} more
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Port Distribution */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <h3 className="text-xl font-black text-black uppercase mb-4 flex items-center gap-2">
            <Activity className="h-5 w-5" />
            Port Analysis
          </h3>

          {networkFindings.length === 0 ? (
            <div className="h-64 flex flex-col items-center justify-center bg-gray-50 border-2 border-gray-300 rounded-lg">
              <Activity className="h-12 w-12 text-gray-400 mx-auto mb-3" />
              <p className="text-gray-600 font-bold">No Ports Detected</p>
              <p className="text-sm text-gray-500 mt-1">Waiting for scan results</p>
            </div>
          ) : (
            <div className="space-y-3">
              {/* Common ports */}
              {Array.from(new Set(networkFindings.map(f => f.port)))
                .sort((a, b) => a - b)
                .slice(0, 8)
                .map(port => {
                  const count = networkFindings.filter(f => f.port === port).length;
                  const percentage = Math.min(100, (count / uniqueHosts.length) * 100);
                  const finding = networkFindings.find(f => f.port === port);
                  return (
                    <div key={port} className="space-y-1">
                      <div className="flex justify-between text-sm">
                        <span className="font-mono font-bold">{port}</span>
                        <span className="text-gray-500">{finding?.service_name || 'Unknown'} - {count} host{count !== 1 ? 's' : ''}</span>
                      </div>
                      <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
                        <div
                          className="h-full bg-gradient-to-r from-blue-500 to-blue-600 rounded-full transition-all"
                          style={{ width: `${percentage}%` }}
                        />
                      </div>
                    </div>
                  );
                })}
            </div>
          )}
        </div>
      </div>

      {/* Recent Findings Table */}
      <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
        <h3 className="text-xl font-black text-black uppercase mb-4 flex items-center gap-2">
          <Shield className="h-5 w-5" />
          Recent Network Findings
        </h3>

        {networkFindings.length === 0 ? (
          <div className="text-center py-8">
            <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-3" />
            <p className="text-gray-600 font-bold">No Network Findings</p>
            <p className="text-sm text-gray-500">Your network scan has not discovered any issues yet</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b-2 border-black">
                  <th className="text-left py-3 px-4 font-black uppercase text-sm">Host</th>
                  <th className="text-left py-3 px-4 font-black uppercase text-sm">Port</th>
                  <th className="text-left py-3 px-4 font-black uppercase text-sm">Type</th>
                  <th className="text-left py-3 px-4 font-black uppercase text-sm">Severity</th>
                  <th className="text-left py-3 px-4 font-black uppercase text-sm">Description</th>
                </tr>
              </thead>
              <tbody>
                {networkFindings.slice(0, 10).map((finding, idx) => (
                  <tr key={finding.id || idx} className="border-b border-gray-200 hover:bg-gray-50">
                    <td className="py-3 px-4 font-mono text-sm">{finding.host}</td>
                    <td className="py-3 px-4">
                      <span className="font-mono bg-gray-100 px-2 py-1 rounded text-sm">
                        {finding.port}/{finding.protocol}
                      </span>
                    </td>
                    <td className="py-3 px-4">
                      <span className="text-sm uppercase">{finding.finding_type}</span>
                    </td>
                    <td className="py-3 px-4">
                      <span className={`px-2 py-1 rounded text-xs font-bold uppercase border ${getSeverityColor(finding.severity)}`}>
                        {finding.severity}
                      </span>
                    </td>
                    <td className="py-3 px-4 text-sm text-gray-600 max-w-xs truncate">{finding.description}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};

export default NetworkSecurityDashboard;
