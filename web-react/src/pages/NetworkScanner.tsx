import React, { useState, useEffect } from 'react';
import toast from 'react-hot-toast';
import { Network, Play, RefreshCw, AlertTriangle, CheckCircle, Search } from 'lucide-react';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { networkScanService } from '@/services/networkScanService';
import { agentService } from '@/services/agentService';

interface NetworkScan {
  id: string;
  agent_id: string;
  status: string;
  start_time: string;
  end_time?: string;
  network_findings: any[];
  metadata?: any;
}

const NetworkScanner: React.FC = () => {
  const [scans, setScans] = useState<NetworkScan[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isScanning, setIsScanning] = useState(false);
  const [selectedScan, setSelectedScan] = useState<NetworkScan | null>(null);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    fetchScans();
    // Auto-refresh every 30 seconds
    const interval = setInterval(fetchScans, 30000);
    return () => clearInterval(interval);
  }, []);

  const fetchScans = async () => {
    try {
      const results = await networkScanService.getNetworkScanResults();
      setScans(results);
    } catch (error) {
      console.error('Failed to fetch network scans:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const initiateScan = async () => {
    try {
      setIsScanning(true);
      // Get first available agent
      const agents = await agentService.getAgents();
      
      if (agents.length === 0) {
        toast.error('No agents available. Please register an agent first.');
        return;
      }

      const agentId = agents[0].id;
      
      // Get scan configuration from environment or use defaults
      const scanTargets = import.meta.env.VITE_DEFAULT_SCAN_TARGETS?.split(',') || ['192.168.1.0/24'];
      const scanType = import.meta.env.VITE_DEFAULT_SCAN_TYPE || 'tcp';
      const timeout = parseInt(import.meta.env.VITE_SCAN_TIMEOUT || '30', 10);
      const concurrency = parseInt(import.meta.env.VITE_SCAN_CONCURRENCY || '10', 10);

      // Initiate network scan
      const scanId = await networkScanService.initiateScan({
        agent_id: agentId,
        targets: scanTargets,
        scan_type: scanType,
        timeout,
        concurrency
      });

      if (scanId) {
        toast.success(`Network scan initiated! Scan ID: ${scanId}`);
        // Refresh scans after a delay
        setTimeout(fetchScans, 2000);
      }
    } catch (error: any) {
      console.error('Failed to initiate scan:', error);
      const errorMessage = error.response?.data?.error || error.message || 'Unknown error';
      toast.error(`Failed to initiate scan: ${errorMessage}`);
    } finally {
      setIsScanning(false);
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status.toLowerCase()) {
      case 'completed':
      case 'success':
        return <Badge className="bg-green-100 text-green-800 border-green-300">Completed</Badge>;
      case 'running':
      case 'in_progress':
        return <Badge className="bg-blue-100 text-blue-800 border-blue-300">Running</Badge>;
      case 'failed':
      case 'error':
        return <Badge className="bg-red-100 text-red-800 border-red-300">Failed</Badge>;
      default:
        return <Badge className="bg-gray-100 text-gray-800 border-gray-300">{status}</Badge>;
    }
  };

  const getSeverityBadge = (severity: string) => {
    const severityLower = severity.toLowerCase();
    if (severityLower.includes('critical')) {
      return <Badge className="bg-red-100 text-red-800 border-red-300">Critical</Badge>;
    } else if (severityLower.includes('high')) {
      return <Badge className="bg-orange-100 text-orange-800 border-orange-300">High</Badge>;
    } else if (severityLower.includes('medium')) {
      return <Badge className="bg-yellow-100 text-yellow-800 border-yellow-300">Medium</Badge>;
    } else {
      return <Badge className="bg-blue-100 text-blue-800 border-blue-300">Low</Badge>;
    }
  };

  const filteredScans = scans.filter(scan =>
    scan.id.toLowerCase().includes(searchTerm.toLowerCase()) ||
    scan.agent_id.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const totalFindings = scans.reduce((sum, scan) => sum + (scan.network_findings?.length || 0), 0);
  const criticalFindings = scans.reduce((sum, scan) => {
    const critical = scan.network_findings?.filter((f: any) => 
      f.severity?.toLowerCase().includes('critical')
    ).length || 0;
    return sum + critical;
  }, 0);

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading network scans...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-black text-black uppercase mb-2">Network Scanner</h1>
          <p className="text-gray-600">Agentless network scanning using Nmap and Nuclei</p>
        </div>
        <Button
          onClick={initiateScan}
          disabled={isScanning}
          className="bg-orange-500 text-white font-bold uppercase border-3 border-black hover:bg-orange-600"
        >
          {isScanning ? (
            <>
              <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
              Scanning...
            </>
          ) : (
            <>
              <Play className="h-4 w-4 mr-2" />
              Start Scan
            </>
          )}
        </Button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 uppercase font-bold mb-1">Total Scans</p>
              <p className="text-3xl font-black text-black">{scans.length}</p>
            </div>
            <Network className="h-8 w-8 text-blue-600" />
          </div>
        </Card>

        <Card className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 uppercase font-bold mb-1">Total Findings</p>
              <p className="text-3xl font-black text-black">{totalFindings}</p>
            </div>
            <AlertTriangle className="h-8 w-8 text-orange-600" />
          </div>
        </Card>

        <Card className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 uppercase font-bold mb-1">Critical</p>
              <p className="text-3xl font-black text-red-600">{criticalFindings}</p>
            </div>
            <AlertTriangle className="h-8 w-8 text-red-600" />
          </div>
        </Card>

        <Card className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 uppercase font-bold mb-1">Completed</p>
              <p className="text-3xl font-black text-green-600">
                {scans.filter(s => s.status === 'completed').length}
              </p>
            </div>
            <CheckCircle className="h-8 w-8 text-green-600" />
          </div>
        </Card>
      </div>

      {/* Search */}
      <Card className="p-4 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
          <Input
            type="text"
            placeholder="Search scans..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-10 border-3 border-black"
          />
        </div>
      </Card>

      {/* Scans List */}
      <div className="space-y-4">
        {filteredScans.length === 0 ? (
          <Card className="p-12 text-center bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <Network className="h-16 w-16 text-gray-400 mx-auto mb-4" />
            <h3 className="text-xl font-bold text-gray-900 mb-2">No Network Scans</h3>
            <p className="text-gray-600 mb-4">
              Start a network scan to discover devices and vulnerabilities on your network.
            </p>
            <Button
              onClick={initiateScan}
              className="bg-orange-500 text-white font-bold uppercase border-3 border-black"
            >
              <Play className="h-4 w-4 mr-2" />
              Start First Scan
            </Button>
          </Card>
        ) : (
          filteredScans.map((scan) => (
            <Card
              key={scan.id}
              className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all cursor-pointer"
              onClick={() => setSelectedScan(selectedScan?.id === scan.id ? null : scan)}
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <Network className="h-6 w-6 text-blue-600" />
                    <h3 className="text-xl font-black text-black uppercase">Scan {scan.id.slice(0, 8)}</h3>
                    {getStatusBadge(scan.status)}
                  </div>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div>
                      <p className="text-gray-600 font-bold">Agent:</p>
                      <p className="text-black">{scan.agent_id.slice(0, 8)}...</p>
                    </div>
                    <div>
                      <p className="text-gray-600 font-bold">Started:</p>
                      <p className="text-black">
                        {new Date(scan.start_time).toLocaleString()}
                      </p>
                    </div>
                    <div>
                      <p className="text-gray-600 font-bold">Findings:</p>
                      <p className="text-black font-bold">{scan.network_findings?.length || 0}</p>
                    </div>
                    <div>
                      <p className="text-gray-600 font-bold">Duration:</p>
                      <p className="text-black">
                        {scan.end_time
                          ? `${Math.round((new Date(scan.end_time).getTime() - new Date(scan.start_time).getTime()) / 1000)}s`
                          : 'Running...'}
                      </p>
                    </div>
                  </div>
                </div>
                <Button
                  variant="ghost"
                  onClick={(e) => {
                    e.stopPropagation();
                    setSelectedScan(selectedScan?.id === scan.id ? null : scan);
                  }}
                >
                  {selectedScan?.id === scan.id ? 'Hide' : 'Show'} Details
                </Button>
              </div>

              {/* Findings Preview */}
              {scan.network_findings && scan.network_findings.length > 0 && (
                <div className="mt-4 pt-4 border-t-2 border-gray-300">
                  <p className="text-sm font-bold text-gray-700 mb-2">Top Findings:</p>
                  <div className="space-y-2">
                    {scan.network_findings.slice(0, 3).map((finding: any, idx: number) => (
                      <div key={idx} className="flex items-center justify-between p-2 bg-gray-50 rounded border-2 border-gray-200">
                        <div className="flex items-center gap-2">
                          <AlertTriangle className="h-4 w-4 text-orange-600" />
                          <span className="font-bold text-black">{finding.title || finding.type}</span>
                        </div>
                        {finding.severity && getSeverityBadge(finding.severity)}
                      </div>
                    ))}
                    {scan.network_findings.length > 3 && (
                      <p className="text-sm text-gray-600 text-center">
                        +{scan.network_findings.length - 3} more findings
                      </p>
                    )}
                  </div>
                </div>
              )}

              {/* Expanded Details */}
              {selectedScan?.id === scan.id && scan.network_findings && scan.network_findings.length > 0 && (
                <div className="mt-4 pt-4 border-t-2 border-black">
                  <h4 className="font-black text-black uppercase mb-3">All Findings</h4>
                  <div className="space-y-3 max-h-96 overflow-y-auto">
                    {scan.network_findings.map((finding: any, idx: number) => (
                      <Card key={idx} className="p-4 bg-gray-50 border-2 border-gray-300">
                        <div className="flex items-start justify-between mb-2">
                          <div>
                            <h5 className="font-bold text-black">{finding.title || finding.type || 'Finding'}</h5>
                            {finding.description && (
                              <p className="text-sm text-gray-600 mt-1">{finding.description}</p>
                            )}
                          </div>
                          {finding.severity && getSeverityBadge(finding.severity)}
                        </div>
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-2 text-xs mt-2">
                          {finding.host && (
                            <div>
                              <span className="font-bold text-gray-700">Host:</span>
                              <span className="text-black ml-1">{finding.host}</span>
                            </div>
                          )}
                          {finding.port && (
                            <div>
                              <span className="font-bold text-gray-700">Port:</span>
                              <span className="text-black ml-1">{finding.port}</span>
                            </div>
                          )}
                          {finding.protocol && (
                            <div>
                              <span className="font-bold text-gray-700">Protocol:</span>
                              <span className="text-black ml-1 uppercase">{finding.protocol}</span>
                            </div>
                          )}
                          {finding.service && (
                            <div>
                              <span className="font-bold text-gray-700">Service:</span>
                              <span className="text-black ml-1">{finding.service}</span>
                            </div>
                          )}
                        </div>
                      </Card>
                    ))}
                  </div>
                </div>
              )}
            </Card>
          ))
        )}
      </div>
    </div>
  );
};

export default NetworkScanner;

