import React, { useState, useEffect } from 'react';
import { 
  Shield, 
  AlertTriangle, 
  Search, 
  Download, 
  RefreshCw,
  Eye,
  Clock,
  Server,
  Zap
} from 'lucide-react';

interface Vulnerability {
  id: string;
  cve_id: string;
  title: string;
  description: string;
  severity: 'critical' | 'high' | 'medium' | 'low';
  cvss_score: number;
  published_date: string;
  last_modified: string;
  affected_software: string[];
  status: 'open' | 'in_progress' | 'resolved' | 'false_positive';
  asset_id: string;
  asset_name: string;
  remediation: string;
  references: string[];
}

interface VulnerabilityStats {
  total: number;
  critical: number;
  high: number;
  medium: number;
  low: number;
  resolved: number;
  open: number;
  in_progress: number;
}

const Vulnerabilities: React.FC = () => {
  const [vulnerabilities, setVulnerabilities] = useState<Vulnerability[]>([]);
  const [stats, setStats] = useState<VulnerabilityStats>({
    total: 0,
    critical: 0,
    high: 0,
    medium: 0,
    low: 0,
    resolved: 0,
    open: 0,
    in_progress: 0
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<'all' | 'critical' | 'high' | 'medium' | 'low'>('all');
  const [search, setSearch] = useState('');

  useEffect(() => {
    fetchVulnerabilities();
  }, []);

  const fetchVulnerabilities = async () => {
    try {
      setLoading(true);
      setError(null);

      // Fetch vulnerabilities from API
      const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/vulnerabilities/`);
      const data = await response.json();

      if (data.success && data.data) {
        setVulnerabilities(data.data);
        
        // Calculate stats
        const newStats: VulnerabilityStats = {
          total: data.data.length,
          critical: data.data.filter((v: Vulnerability) => v.severity === 'critical').length,
          high: data.data.filter((v: Vulnerability) => v.severity === 'high').length,
          medium: data.data.filter((v: Vulnerability) => v.severity === 'medium').length,
          low: data.data.filter((v: Vulnerability) => v.severity === 'low').length,
          resolved: data.data.filter((v: Vulnerability) => v.status === 'resolved').length,
          open: data.data.filter((v: Vulnerability) => v.status === 'open').length,
          in_progress: data.data.filter((v: Vulnerability) => v.status === 'in_progress').length
        };
        setStats(newStats);
      } else {
        // No vulnerabilities found - this is normal for a clean system
        setVulnerabilities([]);
        setStats({
          total: 0,
          critical: 0,
          high: 0,
          medium: 0,
          low: 0,
          resolved: 0,
          open: 0,
          in_progress: 0
        });
      }
    } catch (err) {
      console.error('Failed to fetch vulnerabilities:', err);
      setError('Failed to load vulnerabilities');
    } finally {
      setLoading(false);
    }
  };

  const filteredVulnerabilities = vulnerabilities.filter(vuln => {
    const matchesFilter = filter === 'all' || vuln.severity === filter;
    const matchesSearch = search === '' || 
      vuln.title.toLowerCase().includes(search.toLowerCase()) ||
      vuln.cve_id.toLowerCase().includes(search.toLowerCase()) ||
      vuln.description.toLowerCase().includes(search.toLowerCase());
    return matchesFilter && matchesSearch;
  });

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return 'bg-red-100 text-red-800 border-red-300';
      case 'high': return 'bg-orange-100 text-orange-800 border-orange-300';
      case 'medium': return 'bg-yellow-100 text-yellow-800 border-yellow-300';
      case 'low': return 'bg-green-100 text-green-800 border-green-300';
      default: return 'bg-gray-100 text-gray-800 border-gray-300';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'open': return 'bg-red-100 text-red-800';
      case 'in_progress': return 'bg-yellow-100 text-yellow-800';
      case 'resolved': return 'bg-green-100 text-green-800';
      case 'false_positive': return 'bg-gray-100 text-gray-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <RefreshCw className="h-8 w-8 animate-spin mx-auto mb-4 text-orange-500" />
          <p className="text-gray-600 font-bold">Loading vulnerabilities...</p>
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
              <h1 className="text-3xl font-bold uppercase tracking-wider text-black">Vulnerability Management</h1>
              <p className="text-gray-600 mt-1">Security vulnerabilities and threat intelligence</p>
            </div>
            <div className="flex items-center gap-4">
              <button
                onClick={fetchVulnerabilities}
                className="px-4 py-2 bg-orange-500 text-white border-3 border-black rounded shadow-neubrutalist hover:shadow-neubrutalist-hover hover:translate-x-0.5 hover:translate-y-0.5 transition-all duration-150 ease-in-out font-bold uppercase tracking-wider"
              >
                <RefreshCw className="h-4 w-4 mr-2 inline-block" />
                REFRESH
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto p-6 space-y-6">
        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {/* Total Vulnerabilities */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg hover:shadow-neubrutalist-xl transition-all duration-150">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-blue-100 rounded border-2 border-black">
                <Shield className="h-6 w-6 text-blue-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-black">{stats.total}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Total</div>
              </div>
            </div>
          </div>

          {/* Critical Vulnerabilities */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg hover:shadow-neubrutalist-xl transition-all duration-150">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-red-100 rounded border-2 border-black">
                <AlertTriangle className="h-6 w-6 text-red-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-red-600">{stats.critical}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Critical</div>
              </div>
            </div>
          </div>

          {/* High Vulnerabilities */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg hover:shadow-neubrutalist-xl transition-all duration-150">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-orange-100 rounded border-2 border-black">
                <Zap className="h-6 w-6 text-orange-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-orange-600">{stats.high}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">High</div>
              </div>
            </div>
          </div>

          {/* Resolved Vulnerabilities */}
          <div className="p-6 bg-white border-3 border-black rounded shadow-neubrutalist-lg hover:shadow-neubrutalist-xl transition-all duration-150">
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-green-100 rounded border-2 border-black">
                <Shield className="h-6 w-6 text-green-600" />
              </div>
              <div className="text-right">
                <div className="text-3xl font-bold text-green-600">{stats.resolved}</div>
                <div className="text-sm text-gray-600 uppercase tracking-wider">Resolved</div>
              </div>
            </div>
          </div>
        </div>

        {/* Filters and Search */}
        <div className="bg-white border-3 border-black rounded shadow-neubrutalist-lg p-6">
          <div className="flex flex-col md:flex-row gap-4">
            {/* Search */}
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <input
                  type="text"
                  placeholder="Search vulnerabilities..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="w-full pl-10 pr-4 py-2 border-3 border-black rounded shadow-neubrutalist focus:shadow-neubrutalist-hover focus:outline-none"
                />
              </div>
            </div>

            {/* Filters */}
            <div className="flex gap-2">
              {(['all', 'critical', 'high', 'medium', 'low'] as const).map((filterType) => (
                <button
                  key={filterType}
                  onClick={() => setFilter(filterType)}
                  className={`px-4 py-2 border-3 border-black rounded text-sm font-bold uppercase tracking-wider transition-all duration-150 ${
                    filter === filterType 
                      ? 'bg-orange-100 text-orange-800 border-orange-300 shadow-neubrutalist-md' 
                      : 'bg-white text-black hover:bg-gray-50 hover:shadow-neubrutalist-sm'
                  }`}
                >
                  {filterType} ({filterType === 'all' ? stats.total : stats[filterType]})
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Vulnerabilities List */}
        <div className="bg-white border-3 border-black rounded shadow-neubrutalist-lg">
          <div className="p-6 border-b-3 border-black">
            <h2 className="text-xl font-bold uppercase tracking-wider text-black">
              Vulnerabilities ({filteredVulnerabilities.length})
            </h2>
          </div>

          {error ? (
            <div className="p-6 text-center">
              <AlertTriangle className="h-12 w-12 text-red-500 mx-auto mb-4" />
              <p className="text-red-600 font-bold">{error}</p>
              <button
                onClick={fetchVulnerabilities}
                className="mt-4 px-4 py-2 bg-red-100 text-red-800 border-2 border-red-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-red-200 transition-colors"
              >
                <RefreshCw className="h-4 w-4 mr-2 inline-block" />
                Retry
              </button>
            </div>
          ) : filteredVulnerabilities.length === 0 ? (
            <div className="p-6 text-center">
              <Shield className="h-12 w-12 text-green-500 mx-auto mb-4" />
              <p className="text-gray-600 font-medium">No vulnerabilities found</p>
              <p className="text-sm text-gray-500 mt-2">
                {search ? 'Try adjusting your search criteria' : 'Your system appears to be secure!'}
              </p>
            </div>
          ) : (
            <div className="divide-y divide-gray-200">
              {filteredVulnerabilities.map((vuln) => (
                <div key={vuln.id} className="p-6 hover:bg-gray-50 transition-colors">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-2">
                        <span className={`px-3 py-1 rounded text-xs font-bold uppercase tracking-wider border-2 ${getSeverityColor(vuln.severity)}`}>
                          {vuln.severity}
                        </span>
                        <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider ${getStatusColor(vuln.status)}`}>
                          {vuln.status}
                        </span>
                        <span className="text-sm text-gray-500">{vuln.cve_id}</span>
                      </div>
                      
                      <h3 className="text-lg font-bold text-black mb-2">{vuln.title}</h3>
                      <p className="text-gray-600 mb-3">{vuln.description}</p>
                      
                      <div className="flex items-center gap-4 text-sm text-gray-500">
                        <div className="flex items-center gap-1">
                          <Server className="h-4 w-4" />
                          <span>{vuln.asset_name}</span>
                        </div>
                        <div className="flex items-center gap-1">
                          <Clock className="h-4 w-4" />
                          <span>Published: {new Date(vuln.published_date).toLocaleDateString()}</span>
                        </div>
                        <div className="flex items-center gap-1">
                          <Zap className="h-4 w-4" />
                          <span>CVSS: {vuln.cvss_score}</span>
                        </div>
                      </div>
                    </div>
                    
                    <div className="flex items-center gap-2 ml-4">
                      <button className="p-2 hover:bg-gray-200 rounded border-2 border-black">
                        <Eye className="h-4 w-4" />
                      </button>
                      <button className="p-2 hover:bg-gray-200 rounded border-2 border-black">
                        <Download className="h-4 w-4" />
                      </button>
                    </div>
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

export default Vulnerabilities;