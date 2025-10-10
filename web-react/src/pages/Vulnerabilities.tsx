import React, { useState, useEffect } from 'react';
import { Search, Filter, AlertTriangle, AlertCircle, Info, ExternalLink, Shield, Eye } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';

const Vulnerabilities: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedSeverity, setSelectedSeverity] = useState('all');

  // Fetch real vulnerability data from API
  const { data: vulnData, isLoading, error } = useQuery({
    queryKey: ['vulnerabilities'],
    queryFn: async () => {
      const response = await fetch('http://localhost:8080/api/vulnerabilities/', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      if (!response.ok) {
        throw new Error('Failed to fetch vulnerabilities');
      }
      return response.json();
    },
  });

    const vulnerabilities = vulnData?.data?.vulnerabilities || [];

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'bg-red-100 text-red-800';
      case 'high':
        return 'bg-orange-100 text-orange-800';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800';
      case 'low':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getSeverityIcon = (severity: string) => {
    switch (severity) {
      case 'critical':
        return <AlertTriangle className="h-4 w-4" />;
      case 'high':
        return <AlertCircle className="h-4 w-4" />;
      case 'medium':
      case 'low':
        return <Info className="h-4 w-4" />;
      default:
        return <Info className="h-4 w-4" />;
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  const filteredVulnerabilities = vulnerabilities.filter(vuln => {
    const matchesSearch = vuln.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         vuln.cve_id.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         vuln.package_name.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesSeverity = selectedSeverity === 'all' || vuln.severity === selectedSeverity;
    return matchesSearch && matchesSeverity;
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading vulnerabilities...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <AlertTriangle className="h-12 w-12 text-red-500 mx-auto" />
          <p className="mt-4 text-red-600">Failed to load vulnerabilities</p>
          <p className="text-gray-600">Please check your connection and try again</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Vulnerabilities</h1>
          <p className="text-gray-600">
            {vulnerabilities.length} vulnerabilities found
            {vulnData?.data?.critical && ` • ${vulnData.data.critical} critical`}
            {vulnData?.data?.high && ` • ${vulnData.data.high} high`}
            {vulnData?.data?.medium && ` • ${vulnData.data.medium} medium`}
            {vulnData?.data?.low && ` • ${vulnData.data.low} low`}
          </p>
        </div>
      </div>

      {/* Filters and Search */}
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="flex-1">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
            <input
              type="text"
              placeholder="Search vulnerabilities..."
              className="input-field pl-10"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </div>
        </div>
        <select
          className="input-field max-w-xs"
          value={selectedSeverity}
          onChange={(e) => setSelectedSeverity(e.target.value)}
        >
          <option value="all">All Severities</option>
          <option value="critical">Critical</option>
          <option value="high">High</option>
          <option value="medium">Medium</option>
          <option value="low">Low</option>
        </select>
        <button className="btn-secondary flex items-center">
          <Filter className="h-4 w-4 mr-2" />
          Filter
        </button>
      </div>

      {/* Vulnerabilities List */}
      <div className="space-y-4">
        {filteredVulnerabilities.map((vuln) => (
          <div key={vuln.id} className="card hover:shadow-md transition-shadow">
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <div className="flex items-center space-x-3 mb-2">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getSeverityColor(vuln.severity)}`}>
                    {getSeverityIcon(vuln.severity)}
                    <span className="ml-1 capitalize">{vuln.severity}</span>
                  </span>
                  {vuln.cve_id && (
                    <span className="text-sm text-gray-500 font-mono">{vuln.cve_id}</span>
                  )}
                  {vuln.cvss_score && (
                    <span className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-blue-100 text-blue-800">
                      <Shield className="h-3 w-3 mr-1" />
                      CVSS {vuln.cvss_score}
                    </span>
                  )}
                  {vuln.exploit_available && (
                    <span className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-red-100 text-red-800">
                      <AlertTriangle className="h-3 w-3 mr-1" />
                      Exploit Available
                    </span>
                  )}
                </div>
                
                <h3 className="text-lg font-medium text-gray-900 mb-1">
                  {vuln.title}
                </h3>
                
                <p className="text-gray-600 mb-3">
                  {vuln.description}
                </p>
                
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                  <div>
                    <span className="font-medium text-gray-700">Package:</span>
                    <span className="ml-1 text-gray-600">{vuln.package_name}@{vuln.package_version}</span>
                  </div>
                  <div>
                    <span className="font-medium text-gray-700">Location:</span>
                    <span className="ml-1 text-gray-600">{vuln.location}</span>
                  </div>
                  <div>
                    <span className="font-medium text-gray-700">Detected:</span>
                    <span className="ml-1 text-gray-600">{formatDate(vuln.created_at)}</span>
                  </div>
                </div>
              </div>
              
              <div className="ml-4 flex flex-col space-y-2">
                <button className="btn-primary text-sm">
                  View Details
                </button>
                <button className="btn-secondary text-sm">
                  Mark as Fixed
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Empty State */}
      {filteredVulnerabilities.length === 0 && (
        <div className="text-center py-12">
          <AlertTriangle className="mx-auto h-12 w-12 text-gray-400" />
          <h3 className="mt-2 text-sm font-medium text-gray-900">No vulnerabilities found</h3>
          <p className="mt-1 text-sm text-gray-500">
            {searchTerm || selectedSeverity !== 'all' 
              ? 'Try adjusting your search or filter criteria.' 
              : 'Great! No vulnerabilities detected in your scans.'
            }
          </p>
        </div>
      )}
    </div>
  );
};

export default Vulnerabilities;
