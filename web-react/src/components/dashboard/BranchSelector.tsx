import React, { useState } from 'react';
import { ChevronDown, Building2, MapPin, Users, AlertTriangle, CheckCircle } from 'lucide-react';

interface Branch {
  id: string;
  name: string;
  location: string;
  type: 'headquarters' | 'branch' | 'datacenter' | 'cloud';
  status: 'active' | 'inactive' | 'maintenance';
  metrics: {
    totalAssets: number;
    criticalVulns: number;
    complianceScore: number;
    lastScan: string;
  };
  children?: Branch[];
}

interface BranchSelectorProps {
  branches: Branch[];
  selectedBranchId?: string;
  onBranchSelect: (branchId: string) => void;
  userRole: 'global_ciso' | 'branch_ciso' | 'branch_it_manager' | 'security_analyst' | 'patch_engineer';
  className?: string;
}

const BranchSelector: React.FC<BranchSelectorProps> = ({
  branches,
  selectedBranchId,
  onBranchSelect,
  userRole,
  className = ''
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  const getBranchIcon = (type: string) => {
    switch (type) {
      case 'headquarters':
        return <Building2 className="w-4 h-4" />;
      case 'datacenter':
        return <Building2 className="w-4 h-4" />;
      case 'cloud':
        return <Building2 className="w-4 h-4" />;
      default:
        return <MapPin className="w-4 h-4" />;
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'active':
        return <CheckCircle className="w-4 h-4 text-green-600" />;
      case 'inactive':
        return <AlertTriangle className="w-4 h-4 text-red-600" />;
      case 'maintenance':
        return <AlertTriangle className="w-4 h-4 text-yellow-600" />;
      default:
        return <div className="w-4 h-4" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'border-green-500 bg-green-50';
      case 'inactive':
        return 'border-red-500 bg-red-50';
      case 'maintenance':
        return 'border-yellow-500 bg-yellow-50';
      default:
        return 'border-gray-500 bg-gray-50';
    }
  };

  const filteredBranches = branches.filter(branch =>
    branch.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    branch.location.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const selectedBranch = branches.find(b => b.id === selectedBranchId);

  const renderBranchItem = (branch: Branch, level: number = 0) => (
    <div key={branch.id} className="relative">
      <button
        onClick={() => {
          onBranchSelect(branch.id);
          setIsOpen(false);
        }}
        className={`
          w-full text-left p-3 rounded border-2 border-black bg-white shadow-sm
          hover:shadow-md hover:translate-x-1 hover:translate-y-1
          transition-all duration-150 ease-in-out
          ${selectedBranchId === branch.id ? 'bg-orange-100 border-orange-500' : ''}
          ${level > 0 ? 'ml-6' : ''}
        `}
        style={{ paddingLeft: `${12 + level * 24}px` }}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            {getBranchIcon(branch.type)}
            <div>
              <div className="font-bold text-sm">{branch.name}</div>
              <div className="text-xs text-gray-600">{branch.location}</div>
            </div>
          </div>
          <div className="flex items-center gap-2">
            {getStatusIcon(branch.status)}
            <div className="text-xs text-gray-500">
              {branch.metrics.totalAssets} assets
            </div>
          </div>
        </div>
        
        {/* Metrics */}
        <div className="mt-2 grid grid-cols-3 gap-2 text-xs">
          <div className="text-center">
            <div className="font-bold text-red-600">{branch.metrics.criticalVulns}</div>
            <div className="text-gray-500">Critical</div>
          </div>
          <div className="text-center">
            <div className="font-bold text-green-600">{branch.metrics.complianceScore}%</div>
            <div className="text-gray-500">Compliance</div>
          </div>
          <div className="text-center">
            <div className="font-bold text-blue-600">
              {new Date(branch.metrics.lastScan).toLocaleDateString()}
            </div>
            <div className="text-gray-500">Last Scan</div>
          </div>
        </div>
      </button>
      
      {/* Children branches */}
      {branch.children && branch.children.length > 0 && (
        <div className="ml-4 mt-2 space-y-1">
          {branch.children.map(child => renderBranchItem(child, level + 1))}
        </div>
      )}
    </div>
  );

  return (
    <div className={`relative ${className}`}>
      {/* Selected Branch Display */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full p-4 rounded border-3 border-black bg-white shadow-lg hover:shadow-xl transition-all duration-150 ease-in-out hover:translate-x-1 hover:translate-y-1 hover:shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]"
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            {selectedBranch ? (
              <>
                {getBranchIcon(selectedBranch.type)}
                <div>
                  <div className="font-bold">{selectedBranch.name}</div>
                  <div className="text-sm text-gray-600">{selectedBranch.location}</div>
                </div>
              </>
            ) : (
              <div>
                <div className="font-bold">Select Branch</div>
                <div className="text-sm text-gray-600">Choose a branch to view</div>
              </div>
            )}
          </div>
          <ChevronDown className={`w-5 h-5 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
        </div>
      </button>

      {/* Dropdown */}
      {isOpen && (
        <div className="absolute top-full left-0 right-0 mt-2 bg-white border-3 border-black shadow-xl rounded z-50 max-h-96 overflow-y-auto">
          {/* Search */}
          <div className="p-3 border-b-2 border-black">
            <input
              type="text"
              placeholder="Search branches..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full p-2 border-2 border-black rounded text-sm focus:outline-none focus:border-orange-500 focus:shadow-[4px_4px_0px_0px_rgba(255,107,0,1)]"
            />
          </div>

          {/* Branch List */}
          <div className="p-2 space-y-1">
            {filteredBranches.length > 0 ? (
              filteredBranches.map(branch => renderBranchItem(branch))
            ) : (
              <div className="p-4 text-center text-gray-500">
                No branches found
              </div>
            )}
          </div>

          {/* Role-based actions */}
          {userRole === 'global_ciso' && (
            <div className="p-3 border-t-2 border-black bg-gray-50">
              <button className="w-full p-2 text-sm font-bold uppercase tracking-wider text-orange-600 hover:text-orange-700 transition-colors">
                + Add New Branch
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default BranchSelector;
