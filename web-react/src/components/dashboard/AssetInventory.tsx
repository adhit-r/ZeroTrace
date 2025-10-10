import React, { useState, useMemo } from 'react';
import { 
  Search, 
  Filter, 
  Download, 
  MoreHorizontal, 
  AlertTriangle, 
  Shield, 
  Clock, 
  CheckCircle,
  X,
  ChevronDown,
  ChevronUp
} from 'lucide-react';

interface Asset {
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
  vulnerabilities: {
    critical: number;
    high: number;
    medium: number;
    low: number;
  };
  complianceScore: number;
  riskScore: number;
  suggestedFixes: number;
}

interface AssetInventoryProps {
  assets: Asset[];
  loading?: boolean;
  userRole: string;
  onBulkAction: (action: string, assetIds: string[]) => void;
  onAssetSelect: (assetId: string) => void;
  className?: string;
}

const AssetInventory: React.FC<AssetInventoryProps> = ({
  assets,
  loading = false,
  userRole,
  onBulkAction,
  onAssetSelect,
  className = ''
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedAssets, setSelectedAssets] = useState<string[]>([]);
  const [sortField, setSortField] = useState<keyof Asset>('riskScore');
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc');
  const [showFilters, setShowFilters] = useState(false);
  const [filters, setFilters] = useState({
    criticality: [] as string[],
    status: [] as string[],
    riskLevel: [] as string[]
  });

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

  const getRiskColor = (score: number) => {
    if (score >= 80) return 'text-red-600';
    if (score >= 60) return 'text-orange-600';
    if (score >= 40) return 'text-yellow-600';
    return 'text-green-600';
  };

  const filteredAndSortedAssets = useMemo(() => {
    let filtered = assets.filter(asset => {
      const matchesSearch = 
        asset.hostname.toLowerCase().includes(searchTerm.toLowerCase()) ||
        asset.ip.includes(searchTerm) ||
        asset.owner.toLowerCase().includes(searchTerm.toLowerCase());
      
      const matchesCriticality = filters.criticality.length === 0 || 
        filters.criticality.includes(asset.businessCriticality);
      
      const matchesStatus = filters.status.length === 0 || 
        filters.status.includes(asset.agentStatus);
      
      const matchesRisk = filters.riskLevel.length === 0 || 
        filters.riskLevel.some(level => {
          const score = asset.riskScore;
          switch (level) {
            case 'critical': return score >= 80;
            case 'high': return score >= 60 && score < 80;
            case 'medium': return score >= 40 && score < 60;
            case 'low': return score < 40;
            default: return true;
          }
        });

      return matchesSearch && matchesCriticality && matchesStatus && matchesRisk;
    });

    return filtered.sort((a, b) => {
      const aVal = a[sortField];
      const bVal = b[sortField];
      
      if (typeof aVal === 'string' && typeof bVal === 'string') {
        return sortDirection === 'asc' 
          ? aVal.localeCompare(bVal)
          : bVal.localeCompare(aVal);
      }
      
      if (typeof aVal === 'number' && typeof bVal === 'number') {
        return sortDirection === 'asc' ? aVal - bVal : bVal - aVal;
      }
      
      return 0;
    });
  }, [assets, searchTerm, filters, sortField, sortDirection]);

  const handleSort = (field: keyof Asset) => {
    if (sortField === field) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortDirection('desc');
    }
  };

  const handleSelectAll = () => {
    if (selectedAssets.length === filteredAndSortedAssets.length) {
      setSelectedAssets([]);
    } else {
      setSelectedAssets(filteredAndSortedAssets.map(asset => asset.id));
    }
  };

  const handleSelectAsset = (assetId: string) => {
    setSelectedAssets(prev => 
      prev.includes(assetId) 
        ? prev.filter(id => id !== assetId)
        : [...prev, assetId]
    );
  };

  const bulkActions = [
    { id: 'scan', label: 'Start Scan', icon: <Shield className="w-4 h-4" /> },
    { id: 'patch', label: 'Apply Patches', icon: <CheckCircle className="w-4 h-4" /> },
    { id: 'ignore', label: 'Ignore', icon: <X className="w-4 h-4" /> },
    { id: 'export', label: 'Export', icon: <Download className="w-4 h-4" /> }
  ];

  if (loading) {
    return (
      <div className={`p-6 ${className}`}>
        <div className="animate-pulse space-y-4">
          <div className="h-12 bg-gray-200 rounded border-3 border-black"></div>
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="h-16 bg-gray-200 rounded border-3 border-black"></div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className={`p-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold uppercase tracking-wider">Asset Inventory</h2>
          <p className="text-gray-600">{filteredAndSortedAssets.length} assets</p>
        </div>
        <div className="flex items-center gap-3">
          {selectedAssets.length > 0 && (
            <div className="flex items-center gap-2">
              <span className="text-sm font-bold">{selectedAssets.length} selected</span>
              <div className="flex gap-1">
                {bulkActions.map(action => (
                  <button
                    key={action.id}
                    onClick={() => onBulkAction(action.id, selectedAssets)}
                    className="flex items-center gap-1 px-3 py-1 text-xs font-bold uppercase tracking-wider bg-orange-100 text-orange-800 border-2 border-orange-300 rounded hover:bg-orange-200 transition-colors"
                  >
                    {action.icon}
                    {action.label}
                  </button>
                ))}
              </div>
            </div>
          )}
          <button
            onClick={() => setShowFilters(!showFilters)}
            className="flex items-center gap-2 px-4 py-2 border-2 border-black bg-white rounded hover:bg-gray-50 transition-colors"
          >
            <Filter className="w-4 h-4" />
            Filters
          </button>
        </div>
      </div>

      {/* Search and Filters */}
      <div className="mb-6 space-y-4">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            type="text"
            placeholder="Search assets by hostname, IP, or owner..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-3 border-3 border-black rounded focus:outline-none focus:border-orange-500 focus:shadow-[4px_4px_0px_0px_rgba(255,107,0,1)]"
          />
        </div>

        {showFilters && (
          <div className="p-4 border-3 border-black bg-gray-50 rounded">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-bold mb-2">Criticality</label>
                <div className="space-y-1">
                  {['critical', 'high', 'medium', 'low'].map(level => (
                    <label key={level} className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        checked={filters.criticality.includes(level)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setFilters(prev => ({ ...prev, criticality: [...prev.criticality, level] }));
                          } else {
                            setFilters(prev => ({ ...prev, criticality: prev.criticality.filter(c => c !== level) }));
                          }
                        }}
                        className="rounded border-2 border-black"
                      />
                      <span className="text-sm capitalize">{level}</span>
                    </label>
                  ))}
                </div>
              </div>
              <div>
                <label className="block text-sm font-bold mb-2">Status</label>
                <div className="space-y-1">
                  {['online', 'offline', 'maintenance'].map(status => (
                    <label key={status} className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        checked={filters.status.includes(status)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setFilters(prev => ({ ...prev, status: [...prev.status, status] }));
                          } else {
                            setFilters(prev => ({ ...prev, status: prev.status.filter(s => s !== status) }));
                          }
                        }}
                        className="rounded border-2 border-black"
                      />
                      <span className="text-sm capitalize">{status}</span>
                    </label>
                  ))}
                </div>
              </div>
              <div>
                <label className="block text-sm font-bold mb-2">Risk Level</label>
                <div className="space-y-1">
                  {['critical', 'high', 'medium', 'low'].map(level => (
                    <label key={level} className="flex items-center gap-2">
                      <input
                        type="checkbox"
                        checked={filters.riskLevel.includes(level)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setFilters(prev => ({ ...prev, riskLevel: [...prev.riskLevel, level] }));
                          } else {
                            setFilters(prev => ({ ...prev, riskLevel: prev.riskLevel.filter(r => r !== level) }));
                          }
                        }}
                        className="rounded border-2 border-black"
                      />
                      <span className="text-sm capitalize">{level}</span>
                    </label>
                  ))}
                </div>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Table */}
      <div className="border-3 border-black rounded overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-100 border-b-3 border-black">
              <tr>
                <th className="p-3 text-left">
                  <input
                    type="checkbox"
                    checked={selectedAssets.length === filteredAndSortedAssets.length && filteredAndSortedAssets.length > 0}
                    onChange={handleSelectAll}
                    className="rounded border-2 border-black"
                  />
                </th>
                <th 
                  className="p-3 text-left cursor-pointer hover:bg-gray-200"
                  onClick={() => handleSort('hostname')}
                >
                  <div className="flex items-center gap-2">
                    Hostname
                    {sortField === 'hostname' && (
                      sortDirection === 'asc' ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />
                    )}
                  </div>
                </th>
                <th className="p-3 text-left">IP Address</th>
                <th className="p-3 text-left">Location</th>
                <th className="p-3 text-left">Owner</th>
                <th 
                  className="p-3 text-left cursor-pointer hover:bg-gray-200"
                  onClick={() => handleSort('businessCriticality')}
                >
                  <div className="flex items-center gap-2">
                    Criticality
                    {sortField === 'businessCriticality' && (
                      sortDirection === 'asc' ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />
                    )}
                  </div>
                </th>
                <th className="p-3 text-left">Status</th>
                <th 
                  className="p-3 text-left cursor-pointer hover:bg-gray-200"
                  onClick={() => handleSort('riskScore')}
                >
                  <div className="flex items-center gap-2">
                    Risk Score
                    {sortField === 'riskScore' && (
                      sortDirection === 'asc' ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />
                    )}
                  </div>
                </th>
                <th className="p-3 text-left">Vulnerabilities</th>
                <th className="p-3 text-left">Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredAndSortedAssets.map((asset) => (
                <tr 
                  key={asset.id} 
                  className="border-b border-gray-200 hover:bg-gray-50 cursor-pointer"
                  onClick={() => onAssetSelect(asset.id)}
                >
                  <td className="p-3">
                    <input
                      type="checkbox"
                      checked={selectedAssets.includes(asset.id)}
                      onChange={(e) => {
                        e.stopPropagation();
                        handleSelectAsset(asset.id);
                      }}
                      className="rounded border-2 border-black"
                    />
                  </td>
                  <td className="p-3 font-bold">{asset.hostname}</td>
                  <td className="p-3 text-gray-600">{asset.ip}</td>
                  <td className="p-3 text-gray-600">{asset.location}</td>
                  <td className="p-3 text-gray-600">{asset.owner}</td>
                  <td className="p-3">
                    <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider border-2 ${getCriticalityColor(asset.businessCriticality)}`}>
                      {asset.businessCriticality}
                    </span>
                  </td>
                  <td className="p-3">
                    <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider ${getStatusColor(asset.agentStatus)}`}>
                      {asset.agentStatus}
                    </span>
                  </td>
                  <td className="p-3">
                    <span className={`font-bold ${getRiskColor(asset.riskScore)}`}>
                      {asset.riskScore}
                    </span>
                  </td>
                  <td className="p-3">
                    <div className="flex gap-1">
                      {asset.vulnerabilities.critical > 0 && (
                        <span className="px-2 py-1 bg-red-100 text-red-800 text-xs font-bold rounded border-2 border-red-300">
                          {asset.vulnerabilities.critical}C
                        </span>
                      )}
                      {asset.vulnerabilities.high > 0 && (
                        <span className="px-2 py-1 bg-orange-100 text-orange-800 text-xs font-bold rounded border-2 border-orange-300">
                          {asset.vulnerabilities.high}H
                        </span>
                      )}
                      {asset.vulnerabilities.medium > 0 && (
                        <span className="px-2 py-1 bg-yellow-100 text-yellow-800 text-xs font-bold rounded border-2 border-yellow-300">
                          {asset.vulnerabilities.medium}M
                        </span>
                      )}
                      {asset.vulnerabilities.low > 0 && (
                        <span className="px-2 py-1 bg-green-100 text-green-800 text-xs font-bold rounded border-2 border-green-300">
                          {asset.vulnerabilities.low}L
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="p-3">
                    <button
                      onClick={(e) => e.stopPropagation()}
                      className="p-1 hover:bg-gray-200 rounded"
                    >
                      <MoreHorizontal className="w-4 h-4" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Pagination */}
      <div className="mt-6 flex items-center justify-between">
        <div className="text-sm text-gray-600">
          Showing {filteredAndSortedAssets.length} of {assets.length} assets
        </div>
        <div className="flex items-center gap-2">
          <button className="px-3 py-1 border-2 border-black bg-white rounded hover:bg-gray-50 transition-colors">
            Previous
          </button>
          <span className="px-3 py-1 bg-orange-100 text-orange-800 border-2 border-orange-300 rounded font-bold">
            1
          </span>
          <button className="px-3 py-1 border-2 border-black bg-white rounded hover:bg-gray-50 transition-colors">
            Next
          </button>
        </div>
      </div>
    </div>
  );
};

export default AssetInventory;
