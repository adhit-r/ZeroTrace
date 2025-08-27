import React, { useState, useEffect } from 'react';
import { 
  Shield, 
  AlertTriangle, 
  CheckCircle, 
  TrendingUp, 
  TrendingDown,
  Clock,
  Activity,
  Zap,
  Target,
  Users,
  Server,
  Package,
  FileText,
  ArrowRight,
  Eye,
  Filter,
  Search,
  Bell,
  Settings,
  Download,
  Upload,
  Terminal,
  Cpu,
  HardDrive
} from 'lucide-react';

// Realistic demo data for software vulnerability management
const useDemoData = () => {
  const [data, setData] = useState({
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
    scanStatus: 'idle'
  });

  useEffect(() => {
    // Simulate loading real data
    const loadData = async () => {
      // Simulate API call delay
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setData({
        assets: {
          total: 1247,
          vulnerable: 89,
          critical: 12,
          high: 23,
          medium: 34,
          low: 20,
          lastScan: new Date(Date.now() - 2 * 60 * 60 * 1000) // 2 hours ago
        },
        vulnerabilities: [
          { id: 1, name: 'CVE-2023-1234', severity: 'critical', asset: 'Web Server 01', status: 'open' },
          { id: 2, name: 'CVE-2023-5678', severity: 'high', asset: 'Database Server', status: 'open' },
          { id: 3, name: 'CVE-2023-9012', severity: 'medium', asset: 'File Server', status: 'mitigated' },
        ],
        recentActivity: [
          { id: 1, type: 'scan', message: 'Full system scan completed', time: '2 hours ago' },
          { id: 2, type: 'vulnerability', message: 'New critical vulnerability detected', time: '4 hours ago' },
          { id: 3, type: 'update', message: 'Security patches applied', time: '6 hours ago' },
        ],
        topVulnerableAssets: [
          { name: 'Web Server 01', vulnerabilities: 8, critical: 2 },
          { name: 'Database Server', vulnerabilities: 6, critical: 1 },
          { name: 'File Server', vulnerabilities: 4, critical: 0 },
        ],
        scanStatus: 'idle'
      });
    };

    loadData();
  }, []);

  return data;
};

const Dashboard: React.FC = () => {
  const data = useDemoData();
  const [isScanning, setIsScanning] = useState(false);

  const startScan = async () => {
    setIsScanning(true);
    // Simulate scan process
    setTimeout(() => {
      setIsScanning(false);
    }, 3000);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gold text-glow">SECURITY OVERVIEW</h1>
          <p className="text-text-secondary mt-1">MONITOR SOFTWARE VULNERABILITIES AND ASSET SECURITY</p>
        </div>
        <div className="flex items-center space-x-3">
          <button className="btn btn-secondary btn-sm">
            <Filter className="h-4 w-4 mr-2" />
            FILTER
          </button>
          <button 
            onClick={startScan}
            disabled={isScanning}
            className="btn btn-primary btn-sm"
          >
            {isScanning ? (
              <>
                <div className="animate-spin h-4 w-4 mr-2 border-2 border-current border-t-transparent rounded-full"></div>
                SCANNING...
              </>
            ) : (
              <>
                <Activity className="h-4 w-4 mr-2" />
                START SCAN
              </>
            )}
          </button>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Total Assets */}
        <div className="card card-terminal glow-border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-text-secondary text-sm uppercase tracking-wide">TOTAL ASSETS</p>
              <p className="text-3xl font-bold text-text-primary">{data.assets.total.toLocaleString()}</p>
            </div>
            <div className="h-12 w-12 bg-medium-gray rounded-lg flex items-center justify-center border border-gold">
              <Server className="h-6 w-6 text-gold" />
            </div>
          </div>
        </div>

        {/* Vulnerable Assets */}
        <div className="card card-terminal glow-border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-text-secondary text-sm uppercase tracking-wide">VULNERABLE</p>
              <p className="text-3xl font-bold text-warning">{data.assets.vulnerable}</p>
            </div>
            <div className="h-12 w-12 bg-medium-gray rounded-lg flex items-center justify-center border border-warning">
              <AlertTriangle className="h-6 w-6 text-warning" />
            </div>
          </div>
        </div>

        {/* Critical Vulnerabilities */}
        <div className="card card-critical glow-border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-text-secondary text-sm uppercase tracking-wide">CRITICAL</p>
              <p className="text-3xl font-bold text-critical">{data.assets.critical}</p>
            </div>
            <div className="h-12 w-12 bg-medium-gray rounded-lg flex items-center justify-center border border-critical">
              <Zap className="h-6 w-6 text-critical" />
            </div>
          </div>
        </div>

        {/* Last Scan */}
        <div className="card card-terminal glow-border">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-text-secondary text-sm uppercase tracking-wide">LAST SCAN</p>
              <p className="text-lg font-bold text-text-primary">
                {data.assets.lastScan ? data.assets.lastScan.toLocaleTimeString() : 'NEVER'}
              </p>
            </div>
            <div className="h-12 w-12 bg-medium-gray rounded-lg flex items-center justify-center border border-gold">
              <Clock className="h-6 w-6 text-gold" />
            </div>
          </div>
        </div>
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Recent Vulnerabilities */}
        <div className="lg:col-span-2">
          <div className="card card-terminal">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-bold text-text-primary">RECENT VULNERABILITIES</h2>
              <button className="btn btn-ghost btn-sm">
                VIEW ALL <ArrowRight className="h-4 w-4 ml-2" />
              </button>
            </div>
            
            <div className="space-y-4">
              {data.vulnerabilities.map((vuln) => (
                <div key={vuln.id} className="flex items-center justify-between p-4 bg-medium-gray rounded border border-light-gray">
                  <div className="flex items-center space-x-4">
                    <div className={`badge badge-${vuln.severity}`}>
                      {vuln.severity.toUpperCase()}
                    </div>
                    <div>
                      <p className="text-text-primary font-medium">{vuln.name}</p>
                      <p className="text-text-secondary text-sm">{vuln.asset}</p>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    <span className={`badge ${vuln.status === 'open' ? 'badge-critical' : 'badge-low'}`}>
                      {vuln.status.toUpperCase()}
                    </span>
                    <button className="btn btn-ghost btn-sm">
                      <Eye className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Recent Activity */}
        <div className="lg:col-span-1">
          <div className="card card-terminal">
            <h2 className="text-xl font-bold text-text-primary mb-6">RECENT ACTIVITY</h2>
            
            <div className="space-y-4">
              {data.recentActivity.map((activity) => (
                <div key={activity.id} className="flex items-start space-x-3">
                  <div className="h-2 w-2 bg-gold rounded-full mt-2 flex-shrink-0"></div>
                  <div className="flex-1">
                    <p className="text-text-primary text-sm">{activity.message}</p>
                    <p className="text-text-muted text-xs">{activity.time}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Top Vulnerable Assets */}
      <div className="card card-terminal">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-text-primary">TOP VULNERABLE ASSETS</h2>
          <button className="btn btn-ghost btn-sm">
            VIEW ALL ASSETS <ArrowRight className="h-4 w-4 ml-2" />
          </button>
        </div>
        
        <div className="overflow-x-auto">
          <table className="table">
            <thead>
              <tr>
                <th>ASSET NAME</th>
                <th>TOTAL VULNERABILITIES</th>
                <th>CRITICAL</th>
                <th>STATUS</th>
                <th>ACTIONS</th>
              </tr>
            </thead>
            <tbody>
              {data.topVulnerableAssets.map((asset, index) => (
                <tr key={index}>
                  <td className="text-text-primary font-medium">{asset.name}</td>
                  <td>
                    <span className="badge badge-warning">{asset.vulnerabilities}</span>
                  </td>
                  <td>
                    <span className="badge badge-critical">{asset.critical}</span>
                  </td>
                  <td>
                    <span className="badge badge-info">MONITORING</span>
                  </td>
                  <td>
                    <button className="btn btn-ghost btn-sm">
                      <Eye className="h-4 w-4" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
