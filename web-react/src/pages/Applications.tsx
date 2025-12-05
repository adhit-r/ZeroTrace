import React, { useState, useEffect } from 'react';
import { api } from '@/services/api';
import { Search, Package, AlertTriangle, CheckCircle, XCircle } from 'lucide-react';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';

interface Application {
  name: string;
  version: string;
  vendor: string;
  type: string;
  path?: string;
  agentId: string;
  agentName: string;
  vulnerabilities?: number;
  status?: 'vulnerable' | 'safe' | 'unknown';
}

const Applications: React.FC = () => {
  const [applications, setApplications] = useState<Application[]>([]);
  const [filteredApps, setFilteredApps] = useState<Application[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterType, setFilterType] = useState<'all' | 'vulnerable' | 'safe'>('all');

  useEffect(() => {
    fetchApplications();
  }, []);

  useEffect(() => {
    filterApplications();
  }, [searchTerm, filterType, applications]);

  const fetchApplications = async () => {
    try {
      setIsLoading(true);
      const response = await api.get('/api/agents/');
      const agents = response.data?.data || [];

      const allApps: Application[] = [];
      
      agents.forEach((agent: any) => {
        const dependencies = agent.metadata?.dependencies || [];
        const vulnerabilities = agent.metadata?.vulnerabilities || [];
        
        dependencies.forEach((dep: any) => {
          const appVulns = vulnerabilities.filter((v: any) => 
            v.package_name === dep.name || v.package_name === dep.package_name
          );
          
          allApps.push({
            name: dep.name || dep.package_name,
            version: dep.version || 'unknown',
            vendor: dep.vendor || 'Unknown',
            type: dep.type || 'application',
            path: dep.path,
            agentId: agent.id,
            agentName: agent.name || agent.hostname,
            vulnerabilities: appVulns.length,
            status: appVulns.length > 0 ? 'vulnerable' : 'safe'
          });
        });
      });

      // Deduplicate by name+version+agent
      const uniqueApps = allApps.filter((app, index, self) =>
        index === self.findIndex((a) => 
          a.name === app.name && 
          a.version === app.version && 
          a.agentId === app.agentId
        )
      );

      setApplications(uniqueApps);
      setFilteredApps(uniqueApps);
    } catch (error) {
      console.error('Failed to fetch applications:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const filterApplications = () => {
    let filtered = [...applications];

    // Search filter
    if (searchTerm) {
      filtered = filtered.filter(app =>
        app.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        app.vendor.toLowerCase().includes(searchTerm.toLowerCase()) ||
        app.version.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Status filter
    if (filterType !== 'all') {
      filtered = filtered.filter(app => app.status === filterType);
    }

    setFilteredApps(filtered);
  };

  const getStatusIcon = (status?: string) => {
    switch (status) {
      case 'vulnerable':
        return <AlertTriangle className="h-4 w-4 text-red-600" />;
      case 'safe':
        return <CheckCircle className="h-4 w-4 text-green-600" />;
      default:
        return <XCircle className="h-4 w-4 text-gray-400" />;
    }
  };

  const getStatusBadge = (status?: string) => {
    switch (status) {
      case 'vulnerable':
        return <Badge className="bg-red-100 text-red-800 border-red-300">Vulnerable</Badge>;
      case 'safe':
        return <Badge className="bg-green-100 text-green-800 border-green-300">Safe</Badge>;
      default:
        return <Badge className="bg-gray-100 text-gray-800 border-gray-300">Unknown</Badge>;
    }
  };

  const stats = {
    total: applications.length,
    vulnerable: applications.filter(a => a.status === 'vulnerable').length,
    safe: applications.filter(a => a.status === 'safe').length,
  };

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading applications...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-black text-black uppercase mb-2">Applications</h1>
        <p className="text-gray-600">Manage and monitor all installed applications across your agents</p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 uppercase font-bold mb-1">Total Applications</p>
              <p className="text-3xl font-black text-black">{stats.total}</p>
            </div>
            <Package className="h-8 w-8 text-blue-600" />
          </div>
        </Card>

        <Card className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 uppercase font-bold mb-1">Vulnerable</p>
              <p className="text-3xl font-black text-red-600">{stats.vulnerable}</p>
            </div>
            <AlertTriangle className="h-8 w-8 text-red-600" />
          </div>
        </Card>

        <Card className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 uppercase font-bold mb-1">Safe</p>
              <p className="text-3xl font-black text-green-600">{stats.safe}</p>
            </div>
            <CheckCircle className="h-8 w-8 text-green-600" />
          </div>
        </Card>
      </div>

      {/* Filters */}
      <Card className="p-4 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
        <div className="flex flex-col md:flex-row gap-4">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
            <Input
              type="text"
              placeholder="Search applications..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10 border-3 border-black"
            />
          </div>
          <div className="flex gap-2">
            <button
              onClick={() => setFilterType('all')}
              className={`px-4 py-2 font-bold uppercase border-3 border-black rounded ${
                filterType === 'all' 
                  ? 'bg-orange-500 text-white' 
                  : 'bg-white text-black hover:bg-gray-100'
              }`}
            >
              All
            </button>
            <button
              onClick={() => setFilterType('vulnerable')}
              className={`px-4 py-2 font-bold uppercase border-3 border-black rounded ${
                filterType === 'vulnerable' 
                  ? 'bg-red-500 text-white' 
                  : 'bg-white text-black hover:bg-gray-100'
              }`}
            >
              Vulnerable
            </button>
            <button
              onClick={() => setFilterType('safe')}
              className={`px-4 py-2 font-bold uppercase border-3 border-black rounded ${
                filterType === 'safe' 
                  ? 'bg-green-500 text-white' 
                  : 'bg-white text-black hover:bg-gray-100'
              }`}
            >
              Safe
            </button>
          </div>
        </div>
      </Card>

      {/* Applications List */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {filteredApps.length === 0 ? (
          <Card className="col-span-full p-12 text-center bg-white border-3 border-black rounded-lg shadow-neo-brutal">
            <Package className="h-16 w-16 text-gray-400 mx-auto mb-4" />
            <h3 className="text-xl font-bold text-gray-900 mb-2">No Applications Found</h3>
            <p className="text-gray-600">
              {searchTerm || filterType !== 'all' 
                ? 'Try adjusting your search or filters'
                : 'No applications detected yet. Agents will report installed applications after scanning.'}
            </p>
          </Card>
        ) : (
          filteredApps.map((app, index) => (
            <Card
              key={`${app.agentId}-${app.name}-${index}`}
              className="p-4 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all"
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-2">
                  {getStatusIcon(app.status)}
                  <div>
                    <h3 className="font-black text-black uppercase text-lg">{app.name}</h3>
                    <p className="text-sm text-gray-600">{app.vendor}</p>
                  </div>
                </div>
                {getStatusBadge(app.status)}
              </div>

              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="font-bold text-gray-700">Version:</span>
                  <span className="text-black">{app.version}</span>
                </div>
                <div className="flex justify-between">
                  <span className="font-bold text-gray-700">Type:</span>
                  <span className="text-black uppercase">{app.type}</span>
                </div>
                <div className="flex justify-between">
                  <span className="font-bold text-gray-700">Agent:</span>
                  <span className="text-black">{app.agentName}</span>
                </div>
                {app.vulnerabilities !== undefined && app.vulnerabilities > 0 && (
                  <div className="flex justify-between">
                    <span className="font-bold text-red-700">Vulnerabilities:</span>
                    <span className="text-red-600 font-bold">{app.vulnerabilities}</span>
                  </div>
                )}
              </div>
            </Card>
          ))
        )}
      </div>
    </div>
  );
};

export default Applications;

