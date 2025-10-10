import React, { useState, useEffect } from 'react';
import { 
  AlertTriangle, 
  Shield, 
  Clock, 
  CheckCircle, 
  X, 
  Download, 
  ExternalLink,
  TrendingUp,
  Activity,
  HardDrive,
  Network,
  User,
  Tag,
  Calendar,
  Zap
} from 'lucide-react';

interface Vulnerability {
  id: string;
  cve: string;
  title: string;
  description: string;
  severity: 'critical' | 'high' | 'medium' | 'low';
  cvss: number;
  exploitability: 'exploitable' | 'poc' | 'theoretical' | 'unknown';
  status: 'open' | 'patched' | 'ignored' | 'in_progress';
  published: string;
  suggestedFixes: string[];
  references: string[];
}

interface AssetDetailData {
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
  vulnerabilities: Vulnerability[];
  complianceScore: number;
  riskScore: number;
  suggestedFixes: number;
  metadata: {
    os: string;
    architecture: string;
    kernel: string;
    uptime: string;
    memory: string;
    cpu: string;
    disk: string;
  };
  scanHistory: Array<{
    timestamp: string;
    status: 'completed' | 'failed' | 'in_progress';
    vulnerabilitiesFound: number;
    duration: string;
  }>;
  networkInfo: {
    interfaces: Array<{
      name: string;
      ip: string;
      mac: string;
      status: 'up' | 'down';
    }>;
    openPorts: Array<{
      port: number;
      protocol: string;
      service: string;
      status: 'open' | 'closed' | 'filtered';
    }>;
  };
}

interface AssetDetailProps {
  assetId: string;
  onClose: () => void;
  className?: string;
}

const AssetDetail: React.FC<AssetDetailProps> = ({ assetId, onClose, className = '' }) => {
  const [assetData, setAssetData] = useState<AssetDetailData | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'overview' | 'vulnerabilities' | 'network' | 'history'>('overview');
  const [selectedVulns, setSelectedVulns] = useState<string[]>([]);

  useEffect(() => {
    const fetchAssetData = async () => {
      setLoading(true);
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const mockData: AssetDetailData = {
        id: assetId,
        hostname: 'web-server-01',
        ip: '192.168.1.100',
        branch: 'Headquarters NYC',
        location: 'New York, NY',
        owner: 'IT Team',
        businessCriticality: 'critical',
        tags: ['web', 'production', 'nginx'],
        lastSeen: '2025-01-09T10:30:00Z',
        agentStatus: 'online',
        vulnerabilities: [
          {
            id: 'vuln-001',
            cve: 'CVE-2024-1234',
            title: 'Apache HTTP Server Remote Code Execution',
            description: 'A critical vulnerability in Apache HTTP Server allows remote attackers to execute arbitrary code.',
            severity: 'critical',
            cvss: 9.8,
            exploitability: 'exploitable',
            status: 'open',
            published: '2025-01-05T00:00:00Z',
            suggestedFixes: ['Update to Apache 2.4.58', 'Apply security patch APACHE-2024-001'],
            references: ['https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2024-1234']
          },
          {
            id: 'vuln-002',
            cve: 'CVE-2024-5678',
            title: 'Linux Kernel Privilege Escalation',
            description: 'A vulnerability in the Linux kernel allows local users to escalate privileges.',
            severity: 'high',
            cvss: 7.8,
            exploitability: 'poc',
            status: 'open',
            published: '2025-01-03T00:00:00Z',
            suggestedFixes: ['Update kernel to version 5.15.0-91', 'Apply kernel patch KERNEL-2024-002'],
            references: ['https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2024-5678']
          }
        ],
        complianceScore: 78,
        riskScore: 85,
        suggestedFixes: 7,
        metadata: {
          os: 'Ubuntu 22.04 LTS',
          architecture: 'x86_64',
          kernel: '5.15.0-89-generic',
          uptime: '15 days, 3 hours',
          memory: '16 GB',
          cpu: 'Intel Xeon E5-2680 v4',
          disk: '500 GB SSD'
        },
        scanHistory: [
          { timestamp: '2025-01-09T10:30:00Z', status: 'completed', vulnerabilitiesFound: 2, duration: '5m 23s' },
          { timestamp: '2025-01-08T10:30:00Z', status: 'completed', vulnerabilitiesFound: 1, duration: '4m 45s' },
          { timestamp: '2025-01-07T10:30:00Z', status: 'completed', vulnerabilitiesFound: 3, duration: '6m 12s' }
        ],
        networkInfo: {
          interfaces: [
            { name: 'eth0', ip: '192.168.1.100', mac: '00:1a:2b:3c:4d:5e', status: 'up' },
            { name: 'lo', ip: '127.0.0.1', mac: '00:00:00:00:00:00', status: 'up' }
          ],
          openPorts: [
            { port: 22, protocol: 'tcp', service: 'ssh', status: 'open' },
            { port: 80, protocol: 'tcp', service: 'http', status: 'open' },
            { port: 443, protocol: 'tcp', service: 'https', status: 'open' },
            { port: 3306, protocol: 'tcp', service: 'mysql', status: 'open' }
          ]
        }
      };
      
      setAssetData(mockData);
      setLoading(false);
    };

    fetchAssetData();
  }, [assetId]);

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
      case 'online':
        return 'bg-green-100 text-green-800';
      case 'offline':
        return 'bg-red-100 text-red-800';
      case 'maintenance':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getCriticalityColor = (criticality: string) => {
    switch (criticality) {
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

  const handleBulkVulnAction = (action: string) => {
    console.log(`Bulk vulnerability action: ${action} on vulns:`, selectedVulns);
    // Implement bulk vulnerability actions
  };

  if (loading) {
    return (
      <div className={`fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 ${className}`}>
        <div className="bg-white p-8 rounded border-3 border-black shadow-xl">
          <div className="animate-pulse space-y-4">
            <div className="h-8 bg-gray-200 rounded border-3 border-black"></div>
            <div className="h-64 bg-gray-200 rounded border-3 border-black"></div>
          </div>
        </div>
      </div>
    );
  }

  if (!assetData) {
    return (
      <div className={`fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 ${className}`}>
        <div className="bg-white p-8 rounded border-3 border-black shadow-xl">
          <div className="text-center">
            <h2 className="text-2xl font-bold text-red-600 mb-4">Asset not found</h2>
            <p className="text-gray-600 mb-4">The requested asset could not be found.</p>
            <button
              onClick={onClose}
              className="px-4 py-2 bg-orange-100 text-orange-800 border-2 border-orange-300 rounded font-bold uppercase tracking-wider hover:bg-orange-200 transition-colors"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={`fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4 ${className}`}>
      <div className="bg-white rounded border-3 border-black shadow-xl max-w-6xl w-full max-h-[90vh] overflow-hidden">
        {/* Header */}
        <div className="p-6 border-b-3 border-black bg-gray-50">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div>
                <h2 className="text-2xl font-bold uppercase tracking-wider">{assetData.hostname}</h2>
                <p className="text-gray-600">{assetData.ip} • {assetData.branch}</p>
              </div>
              <div className="flex items-center gap-2">
                <span className={`px-3 py-1 rounded text-sm font-bold uppercase tracking-wider border-2 ${getCriticalityColor(assetData.businessCriticality)}`}>
                  {assetData.businessCriticality}
                </span>
                <span className={`px-3 py-1 rounded text-sm font-bold uppercase tracking-wider ${getStatusColor(assetData.agentStatus)}`}>
                  {assetData.agentStatus}
                </span>
              </div>
            </div>
            <button
              onClick={onClose}
              className="p-2 hover:bg-gray-200 rounded border-2 border-black"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Tab Navigation */}
        <div className="flex border-b-2 border-black">
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
            onClick={() => setActiveTab('vulnerabilities')}
            className={`px-6 py-3 text-sm font-bold uppercase tracking-wider transition-colors ${
              activeTab === 'vulnerabilities' 
                ? 'bg-orange-100 text-orange-800 border-r-2 border-black' 
                : 'bg-white text-gray-600 hover:bg-gray-50'
            }`}
          >
            Vulnerabilities ({assetData.vulnerabilities.length})
          </button>
          <button
            onClick={() => setActiveTab('network')}
            className={`px-6 py-3 text-sm font-bold uppercase tracking-wider transition-colors ${
              activeTab === 'network' 
                ? 'bg-orange-100 text-orange-800 border-r-2 border-black' 
                : 'bg-white text-gray-600 hover:bg-gray-50'
            }`}
          >
            Network
          </button>
          <button
            onClick={() => setActiveTab('history')}
            className={`px-6 py-3 text-sm font-bold uppercase tracking-wider transition-colors ${
              activeTab === 'history' 
                ? 'bg-orange-100 text-orange-800' 
                : 'bg-white text-gray-600 hover:bg-gray-50'
            }`}
          >
            Scan History
          </button>
        </div>

        {/* Content */}
        <div className="p-6 overflow-y-auto max-h-[60vh]">
          {activeTab === 'overview' && (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Asset Info */}
              <div className="space-y-4">
                <div className="p-4 border-3 border-black bg-white rounded shadow-lg">
                  <h3 className="text-lg font-bold uppercase tracking-wider mb-4 flex items-center gap-2">
                    <HardDrive className="w-5 h-5" />
                    Asset Information
                  </h3>
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span className="text-gray-600">Hostname:</span>
                      <span className="font-bold">{assetData.hostname}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">IP Address:</span>
                      <span className="font-bold">{assetData.ip}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Owner:</span>
                      <span className="font-bold">{assetData.owner}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Last Seen:</span>
                      <span className="font-bold">{new Date(assetData.lastSeen).toLocaleString()}</span>
                    </div>
                  </div>
                </div>

                <div className="p-4 border-3 border-black bg-white rounded shadow-lg">
                  <h3 className="text-lg font-bold uppercase tracking-wider mb-4 flex items-center gap-2">
                    <Activity className="w-5 h-5" />
                    System Information
                  </h3>
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span className="text-gray-600">OS:</span>
                      <span className="font-bold">{assetData.metadata.os}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Architecture:</span>
                      <span className="font-bold">{assetData.metadata.architecture}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Kernel:</span>
                      <span className="font-bold">{assetData.metadata.kernel}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Uptime:</span>
                      <span className="font-bold">{assetData.metadata.uptime}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Memory:</span>
                      <span className="font-bold">{assetData.metadata.memory}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">CPU:</span>
                      <span className="font-bold">{assetData.metadata.cpu}</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Risk Metrics */}
              <div className="space-y-4">
                <div className="p-4 border-3 border-black bg-white rounded shadow-lg">
                  <h3 className="text-lg font-bold uppercase tracking-wider mb-4 flex items-center gap-2">
                    <AlertTriangle className="w-5 h-5" />
                    Risk Assessment
                  </h3>
                  <div className="space-y-4">
                    <div className="flex justify-between items-center">
                      <span className="text-gray-600">Risk Score:</span>
                      <span className="text-2xl font-bold text-red-600">{assetData.riskScore}</span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-gray-600">Compliance Score:</span>
                      <span className="text-2xl font-bold text-green-600">{assetData.complianceScore}%</span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-gray-600">Suggested Fixes:</span>
                      <span className="text-2xl font-bold text-orange-600">{assetData.suggestedFixes}</span>
                    </div>
                  </div>
                </div>

                <div className="p-4 border-3 border-black bg-white rounded shadow-lg">
                  <h3 className="text-lg font-bold uppercase tracking-wider mb-4 flex items-center gap-2">
                    <Tag className="w-5 h-5" />
                    Tags
                  </h3>
                  <div className="flex flex-wrap gap-2">
                    {assetData.tags.map((tag, index) => (
                      <span
                        key={index}
                        className="px-2 py-1 bg-gray-100 text-gray-800 border-2 border-gray-300 rounded text-sm font-bold uppercase tracking-wider"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'vulnerabilities' && (
            <div className="space-y-4">
              {/* Bulk Actions */}
              {selectedVulns.length > 0 && (
                <div className="p-4 border-3 border-black bg-orange-50 rounded">
                  <div className="flex items-center justify-between">
                    <span className="font-bold">{selectedVulns.length} vulnerabilities selected</span>
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleBulkVulnAction('patch')}
                        className="px-3 py-1 bg-green-100 text-green-800 border-2 border-green-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-green-200 transition-colors"
                      >
                        Apply Patches
                      </button>
                      <button
                        onClick={() => handleBulkVulnAction('ignore')}
                        className="px-3 py-1 bg-yellow-100 text-yellow-800 border-2 border-yellow-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-yellow-200 transition-colors"
                      >
                        Ignore
                      </button>
                      <button
                        onClick={() => handleBulkVulnAction('export')}
                        className="px-3 py-1 bg-blue-100 text-blue-800 border-2 border-blue-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-blue-200 transition-colors"
                      >
                        Export
                      </button>
                    </div>
                  </div>
                </div>
              )}

              {/* Vulnerabilities List */}
              <div className="space-y-4">
                {assetData.vulnerabilities.map((vuln) => (
                  <div
                    key={vuln.id}
                    className="p-4 border-3 border-black bg-white rounded shadow-lg"
                  >
                    <div className="flex items-start justify-between mb-4">
                      <div className="flex items-center gap-3">
                        <input
                          type="checkbox"
                          checked={selectedVulns.includes(vuln.id)}
                          onChange={(e) => {
                            if (e.target.checked) {
                              setSelectedVulns(prev => [...prev, vuln.id]);
                            } else {
                              setSelectedVulns(prev => prev.filter(id => id !== vuln.id));
                            }
                          }}
                          className="rounded border-2 border-black"
                        />
                        <div>
                          <h4 className="font-bold text-lg">{vuln.title}</h4>
                          <p className="text-gray-600">{vuln.cve}</p>
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        <span className={`px-3 py-1 rounded text-sm font-bold uppercase tracking-wider border-2 ${getSeverityColor(vuln.severity)}`}>
                          {vuln.severity}
                        </span>
                        <span className="font-bold text-lg">CVSS {vuln.cvss}</span>
                      </div>
                    </div>
                    
                    <p className="text-gray-700 mb-4">{vuln.description}</p>
                    
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <h5 className="font-bold mb-2">Suggested Fixes:</h5>
                        <ul className="space-y-1">
                          {vuln.suggestedFixes.map((fix, index) => (
                            <li key={index} className="text-sm text-gray-600">• {fix}</li>
                          ))}
                        </ul>
                      </div>
                      <div>
                        <h5 className="font-bold mb-2">References:</h5>
                        <div className="space-y-1">
                          {vuln.references.map((ref, index) => (
                            <a
                              key={index}
                              href={ref}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="flex items-center gap-1 text-sm text-blue-600 hover:text-blue-800"
                            >
                              <ExternalLink className="w-3 h-3" />
                              {ref}
                            </a>
                          ))}
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {activeTab === 'network' && (
            <div className="space-y-6">
              <div className="p-4 border-3 border-black bg-white rounded shadow-lg">
                <h3 className="text-lg font-bold uppercase tracking-wider mb-4 flex items-center gap-2">
                  <Network className="w-5 h-5" />
                  Network Interfaces
                </h3>
                <div className="space-y-2">
                  {assetData.networkInfo.interfaces.map((iface, index) => (
                    <div key={index} className="flex items-center justify-between p-3 border-2 border-black bg-gray-50 rounded">
                      <div>
                        <span className="font-bold">{iface.name}</span>
                        <span className="text-gray-600 ml-2">{iface.ip}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <span className="text-sm text-gray-600">{iface.mac}</span>
                        <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider ${
                          iface.status === 'up' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                        }`}>
                          {iface.status}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <div className="p-4 border-3 border-black bg-white rounded shadow-lg">
                <h3 className="text-lg font-bold uppercase tracking-wider mb-4 flex items-center gap-2">
                  <Zap className="w-5 h-5" />
                  Open Ports
                </h3>
                <div className="space-y-2">
                  {assetData.networkInfo.openPorts.map((port, index) => (
                    <div key={index} className="flex items-center justify-between p-3 border-2 border-black bg-gray-50 rounded">
                      <div>
                        <span className="font-bold">{port.port}/{port.protocol}</span>
                        <span className="text-gray-600 ml-2">{port.service}</span>
                      </div>
                      <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider ${
                        port.status === 'open' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                      }`}>
                        {port.status}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}

          {activeTab === 'history' && (
            <div className="space-y-4">
              {assetData.scanHistory.map((scan, index) => (
                <div key={index} className="p-4 border-3 border-black bg-white rounded shadow-lg">
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      <Calendar className="w-4 h-4" />
                      <span className="font-bold">{new Date(scan.timestamp).toLocaleString()}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider ${
                        scan.status === 'completed' ? 'bg-green-100 text-green-800' : 
                        scan.status === 'failed' ? 'bg-red-100 text-red-800' : 
                        'bg-yellow-100 text-yellow-800'
                      }`}>
                        {scan.status}
                      </span>
                      <span className="text-sm text-gray-600">{scan.duration}</span>
                    </div>
                  </div>
                  <div className="flex items-center gap-4 text-sm text-gray-600">
                    <span>Vulnerabilities Found: <span className="font-bold">{scan.vulnerabilitiesFound}</span></span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default AssetDetail;
