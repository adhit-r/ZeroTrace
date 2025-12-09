import { api } from './api';

export interface InitiateScanRequest {
  agent_id: string;
  targets: string[];
  scan_type: string;
  timeout: number;
  concurrency: number;
}

export interface NetworkScanResult {
  id: string;
  agent_id: string;
  company_id: string;
  start_time: string;
  end_time: string;
  status: string;
  network_findings: NetworkFinding[];
  metadata: {
    total_hosts?: number;
    total_findings?: number;
    port_findings?: number;
    vuln_findings?: number;
    scan_method?: string;
  };
}

export interface NetworkFinding {
  id: string;
  finding_type: string;
  severity: string;
  host: string;
  port: number;
  protocol: string;
  service_name: string;
  service_version: string;
  banner: string;
  description: string;
  remediation: string;
  discovered_at: string;
  status: string;
  device_type?: string;
  os?: string;
  os_version?: string;
  metadata?: any;
}

export interface NetworkTopologyData {
  nodes: NetworkNode[];
  edges: NetworkEdge[];
}

export interface NetworkNode {
  id: string;
  type: 'device' | 'vulnerability' | 'service';
  label: string;
  deviceType?: string;
  ipAddress?: string;
  os?: string;
  riskScore?: number;
  vulnerabilities?: number;
  openPorts?: NetworkFinding[];
  status?: string;
  severity?: string;
  description?: string;
  cve?: string;
  port?: number;
  protocol?: string;
}

export interface NetworkEdge {
  source: string;
  target: string;
  type: string;
  strength?: number;
}

export const networkScanService = {
  async initiateScan(request: InitiateScanRequest): Promise<string | null> {
    try {
      const response = await api.post<{ scan_id: string }>('/api/v2/scans/network', request);
      return response.data.scan_id || null;
    } catch (error) {
      console.error('Failed to initiate network scan:', error);
      throw error;
    }
  },

  async getNetworkScanResults(): Promise<NetworkScanResult[]> {
    try {
      // Fetch agents and extract network scan results from metadata
      const response = await api.get('/api/agents/');
      const agents = response.data?.data || [];

      const scanResults: NetworkScanResult[] = [];

      agents.forEach((agent: any) => {
        if (agent.metadata?.network_scan_result) {
          const scanResult = agent.metadata.network_scan_result;
          scanResults.push({
            id: scanResult.id || scanResult.scan_id,
            agent_id: scanResult.agent_id || agent.id,
            company_id: scanResult.company_id || '',
            start_time: scanResult.start_time || scanResult.startTime,
            end_time: scanResult.end_time || scanResult.endTime,
            status: scanResult.status || 'completed',
            network_findings: scanResult.network_findings || scanResult.NetworkFindings || [],
            metadata: scanResult.metadata || {},
          });
        }
      });

      return scanResults;
    } catch (error) {
      console.error('Failed to fetch network scan results:', error);
      return [];
    }
  },

  async getNetworkTopology(): Promise<NetworkTopologyData> {
    try {
      const scanResults = await this.getNetworkScanResults();
      
      const deviceMap = new Map<string, NetworkNode>();
      const vulnNodes: NetworkNode[] = [];
      const serviceNodes: NetworkNode[] = [];
      const edges: NetworkEdge[] = [];

      // Process all findings
      scanResults.forEach((scanResult) => {
        scanResult.network_findings.forEach((finding) => {
          // Create or update device node
          if (!deviceMap.has(finding.host)) {
            deviceMap.set(finding.host, {
              id: `device-${finding.host}`,
              type: 'device',
              label: finding.host,
              deviceType: finding.device_type || 'unknown',
              ipAddress: finding.host,
              os: finding.os || 'Unknown',
              riskScore: 0,
              vulnerabilities: 0,
              openPorts: [],
              status: 'online',
            });
          }

          const device = deviceMap.get(finding.host)!;

          // Add finding to appropriate category
          if (finding.finding_type === 'vuln') {
            device.vulnerabilities = (device.vulnerabilities || 0) + 1;
            vulnNodes.push({
              id: `vuln-${finding.id}`,
              type: 'vulnerability',
              label: finding.description.substring(0, 40) + '...',
              severity: finding.severity,
              description: finding.description,
              cve: finding.metadata?.cve || finding.metadata?.template_id,
              ipAddress: finding.host,
              port: finding.port,
            });

            // Create edge from device to vulnerability
            edges.push({
              source: device.id,
              target: `vuln-${finding.id}`,
              type: 'vulnerability',
            });
          } else if (finding.finding_type === 'port') {
            if (!device.openPorts) {
              device.openPorts = [];
            }
            device.openPorts.push(finding);

            serviceNodes.push({
              id: `service-${finding.host}-${finding.port}`,
              type: 'service',
              label: finding.service_name || `Port ${finding.port}`,
              port: finding.port,
              protocol: finding.protocol,
              ipAddress: finding.host,
            });

            // Create edge from device to service
            edges.push({
              source: device.id,
              target: `service-${finding.host}-${finding.port}`,
              type: 'service',
            });
          }

          // Update risk score based on findings
          if (finding.finding_type === 'vuln' || finding.finding_type === 'config') {
            let riskIncrease = 0;
            switch (finding.severity) {
              case 'critical':
                riskIncrease = 20;
                break;
              case 'high':
                riskIncrease = 15;
                break;
              case 'medium':
                riskIncrease = 10;
                break;
              case 'low':
                riskIncrease = 5;
                break;
            }
            device.riskScore = Math.min((device.riskScore || 10) + riskIncrease, 100);
          }
        });
      });

      return {
        nodes: [...Array.from(deviceMap.values()), ...vulnNodes, ...serviceNodes],
        edges,
      };
    } catch (error) {
      console.error('Failed to build network topology:', error);
      return { nodes: [], edges: [] };
    }
  },
};

