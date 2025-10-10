import React, { useState, useEffect } from 'react';
import GlobalDashboard from './GlobalDashboard';
import { 
  Shield, 
  Building2, 
  Users, 
  AlertTriangle, 
  Target, 
  Clock,
  TrendingUp,
  CheckCircle,
  X,
  Download
} from 'lucide-react';

interface RoleBasedDashboardProps {
  userRole: 'global_ciso' | 'branch_ciso' | 'branch_it_manager' | 'security_analyst' | 'patch_engineer';
  className?: string;
}

const RoleBasedDashboard: React.FC<RoleBasedDashboardProps> = ({ userRole, className = '' }) => {
  const [selectedBranch, setSelectedBranch] = useState<string>();
  const [loading, setLoading] = useState(true);

  const getRoleConfig = (role: string) => {
    switch (role) {
      case 'global_ciso':
        return {
          title: 'Global CISO Dashboard',
          description: 'Enterprise-wide security oversight and risk management',
          quickActions: [
            { id: 'scan-all', label: 'Scan All Branches', icon: <Target className="w-4 h-4" />, color: 'orange' },
            { id: 'compliance-report', label: 'Generate Compliance Report', icon: <Shield className="w-4 h-4" />, color: 'green' },
            { id: 'risk-assessment', label: 'Risk Assessment', icon: <AlertTriangle className="w-4 h-4" />, color: 'red' },
            { id: 'export-data', label: 'Export All Data', icon: <Download className="w-4 h-4" />, color: 'blue' }
          ],
          widgets: ['kpi-ribbon', 'risk-heatmap', 'branch-comparison', 'compliance-scorecard', 'recent-alerts', 'top-risks']
        };
      case 'branch_ciso':
        return {
          title: 'Branch CISO Dashboard',
          description: 'Branch-level security management and compliance oversight',
          quickActions: [
            { id: 'scan-branch', label: 'Scan Branch Assets', icon: <Target className="w-4 h-4" />, color: 'orange' },
            { id: 'compliance-check', label: 'Compliance Check', icon: <Shield className="w-4 h-4" />, color: 'green' },
            { id: 'vulnerability-review', label: 'Review Vulnerabilities', icon: <AlertTriangle className="w-4 h-4" />, color: 'red' },
            { id: 'team-report', label: 'Team Report', icon: <Users className="w-4 h-4" />, color: 'blue' }
          ],
          widgets: ['kpi-ribbon', 'asset-inventory', 'vulnerability-trends', 'compliance-status', 'recent-alerts', 'patch-recommendations']
        };
      case 'branch_it_manager':
        return {
          title: 'Branch IT Manager Dashboard',
          description: 'Day-to-day asset management and vulnerability remediation',
          quickActions: [
            { id: 'patch-assets', label: 'Apply Patches', icon: <CheckCircle className="w-4 h-4" />, color: 'green' },
            { id: 'scan-assets', label: 'Scan Assets', icon: <Target className="w-4 h-4" />, color: 'orange' },
            { id: 'create-tickets', label: 'Create Tickets', icon: <Clock className="w-4 h-4" />, color: 'blue' },
            { id: 'asset-inventory', label: 'Asset Inventory', icon: <Building2 className="w-4 h-4" />, color: 'purple' }
          ],
          widgets: ['kpi-ribbon', 'asset-inventory', 'vulnerability-list', 'patch-queue', 'scan-status', 'recent-tickets']
        };
      case 'security_analyst':
        return {
          title: 'Security Analyst Dashboard',
          description: 'Threat investigation and vulnerability analysis',
          quickActions: [
            { id: 'investigate-alert', label: 'Investigate Alert', icon: <AlertTriangle className="w-4 h-4" />, color: 'red' },
            { id: 'export-data', label: 'Export Data', icon: <Download className="w-4 h-4" />, color: 'blue' },
            { id: 'create-report', label: 'Create Report', icon: <TrendingUp className="w-4 h-4" />, color: 'green' },
            { id: 'threat-hunting', label: 'Threat Hunting', icon: <Target className="w-4 h-4" />, color: 'orange' }
          ],
          widgets: ['kpi-ribbon', 'threat-intelligence', 'vulnerability-analysis', 'attack-surface', 'recent-alerts', 'investigation-tools']
        };
      case 'patch_engineer':
        return {
          title: 'Patch Engineer Dashboard',
          description: 'Patch management and deployment coordination',
          quickActions: [
            { id: 'deploy-patches', label: 'Deploy Patches', icon: <CheckCircle className="w-4 h-4" />, color: 'green' },
            { id: 'test-patches', label: 'Test Patches', icon: <Target className="w-4 h-4" />, color: 'orange' },
            { id: 'rollback-patches', label: 'Rollback Patches', icon: <X className="w-4 h-4" />, color: 'red' },
            { id: 'patch-schedule', label: 'Schedule Patches', icon: <Clock className="w-4 h-4" />, color: 'blue' }
          ],
          widgets: ['kpi-ribbon', 'patch-queue', 'deployment-status', 'rollback-plans', 'patch-testing', 'deployment-timeline']
        };
      default:
        return {
          title: 'Dashboard',
          description: 'Security management dashboard',
          quickActions: [],
          widgets: ['kpi-ribbon']
        };
    }
  };

  const roleConfig = getRoleConfig(userRole);

  useEffect(() => {
    setLoading(true);
    // Simulate loading
    setTimeout(() => setLoading(false), 1000);
  }, [userRole]);

  if (loading) {
    return (
      <div className={`p-6 ${className}`}>
        <div className="animate-pulse space-y-6">
          <div className="h-8 bg-gray-200 rounded border-3 border-black"></div>
          <div className="h-32 bg-gray-200 rounded border-3 border-black"></div>
          <div className="h-64 bg-gray-200 rounded border-3 border-black"></div>
        </div>
      </div>
    );
  }

  return (
    <div className={`p-6 space-y-6 ${className}`}>
      {/* Role-specific Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold uppercase tracking-wider">{roleConfig.title}</h1>
          <p className="text-gray-600">{roleConfig.description}</p>
        </div>
        <div className="flex items-center gap-2">
          <span className="px-3 py-1 bg-orange-100 text-orange-800 border-2 border-orange-300 rounded text-sm font-bold uppercase tracking-wider">
            {userRole.replace('_', ' ')}
          </span>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="p-6 border-3 border-black bg-white rounded shadow-lg">
        <h2 className="text-xl font-bold uppercase tracking-wider mb-4">Quick Actions</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {roleConfig.quickActions.map((action) => (
            <button
              key={action.id}
              className={`
                p-4 rounded border-3 border-black bg-white shadow-lg hover:shadow-xl
                hover:translate-x-1 hover:translate-y-1 transition-all duration-150 ease-in-out
                hover:shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]
                ${action.color === 'orange' ? 'hover:bg-orange-50' : ''}
                ${action.color === 'green' ? 'hover:bg-green-50' : ''}
                ${action.color === 'red' ? 'hover:bg-red-50' : ''}
                ${action.color === 'blue' ? 'hover:bg-blue-50' : ''}
                ${action.color === 'purple' ? 'hover:bg-purple-50' : ''}
              `}
            >
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded border-2 border-black ${
                  action.color === 'orange' ? 'bg-orange-100 text-orange-800' : ''
                } ${action.color === 'green' ? 'bg-green-100 text-green-800' : ''}
                ${action.color === 'red' ? 'bg-red-100 text-red-800' : ''}
                ${action.color === 'blue' ? 'bg-blue-100 text-blue-800' : ''}
                ${action.color === 'purple' ? 'bg-purple-100 text-purple-800' : ''}`}>
                  {action.icon}
                </div>
                <span className="font-bold text-sm">{action.label}</span>
              </div>
            </button>
          ))}
        </div>
      </div>

      {/* Role-specific KPIs */}
      <div className="p-6 border-3 border-black bg-white rounded shadow-lg">
        <h2 className="text-xl font-bold uppercase tracking-wider mb-4">Key Metrics</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {userRole === 'global_ciso' && (
            <>
              <div className="p-4 border-2 border-black bg-red-50 rounded">
                <div className="flex items-center gap-2 mb-2">
                  <AlertTriangle className="w-5 h-5 text-red-600" />
                  <span className="font-bold">Critical Branches</span>
                </div>
                <div className="text-2xl font-bold text-red-600">3</div>
                <div className="text-sm text-gray-600">Require immediate attention</div>
              </div>
              <div className="p-4 border-2 border-black bg-orange-50 rounded">
                <div className="flex items-center gap-2 mb-2">
                  <Building2 className="w-5 h-5 text-orange-600" />
                  <span className="font-bold">Total Branches</span>
                </div>
                <div className="text-2xl font-bold text-orange-600">12</div>
                <div className="text-sm text-gray-600">Active locations</div>
              </div>
              <div className="p-4 border-2 border-black bg-green-50 rounded">
                <div className="flex items-center gap-2 mb-2">
                  <Shield className="w-5 h-5 text-green-600" />
                  <span className="font-bold">Compliance Score</span>
                </div>
                <div className="text-2xl font-bold text-green-600">87%</div>
                <div className="text-sm text-gray-600">Enterprise average</div>
              </div>
              <div className="p-4 border-2 border-black bg-blue-50 rounded">
                <div className="flex items-center gap-2 mb-2">
                  <TrendingUp className="w-5 h-5 text-blue-600" />
                  <span className="font-bold">Risk Trend</span>
                </div>
                <div className="text-2xl font-bold text-blue-600">-12%</div>
                <div className="text-sm text-gray-600">vs last month</div>
              </div>
            </>
          )}
          
          {userRole === 'branch_ciso' && (
            <>
              <div className="p-4 border-2 border-black bg-red-50 rounded">
                <div className="flex items-center gap-2 mb-2">
                  <AlertTriangle className="w-5 h-5 text-red-600" />
                  <span className="font-bold">Critical Assets</span>
                </div>
                <div className="text-2xl font-bold text-red-600">8</div>
                <div className="text-sm text-gray-600">Need immediate patching</div>
              </div>
              <div className="p-4 border-2 border-black bg-orange-50 rounded">
                <div className="flex items-center gap-2 mb-2">
                  <Target className="w-5 h-5 text-orange-600" />
                  <span className="font-bold">Scan Coverage</span>
                </div>
                <div className="text-2xl font-bold text-orange-600">94%</div>
                <div className="text-sm text-gray-600">Assets scanned</div>
              </div>
              <div className="p-4 border-2 border-black bg-green-50 rounded">
                <div className="flex items-center gap-2 mb-2">
                  <Shield className="w-5 h-5 text-green-600" />
                  <span className="font-bold">Compliance</span>
                </div>
                <div className="text-2xl font-bold text-green-600">89%</div>
                <div className="text-sm text-gray-600">Branch compliance</div>
              </div>
              <div className="p-4 border-2 border-black bg-blue-50 rounded">
                <div className="flex items-center gap-2 mb-2">
                  <Clock className="w-5 h-5 text-blue-600" />
                  <span className="font-bold">MTTR</span>
                </div>
                <div className="text-2xl font-bold text-blue-600">5.2d</div>
                <div className="text-sm text-gray-600">Mean time to remediate</div>
              </div>
            </>
          )}

          {userRole === 'branch_it_manager' && (
            <>
              <div className="p-4 border-2 border-black bg-red-50 rounded">
                <div className="flex items-center gap-2 mb-2">
                  <AlertTriangle className="w-5 h-5 text-red-600" />
                  <span className="font-bold">Pending Patches</span>
                </div>
                <div className="text-2xl font-bold text-red-600">23</div>
                <div className="text-sm text-gray-600">Critical patches</div>
              </div>
              <div className="p-4 border-2 border-black bg-orange-50 rounded">
                <div className="flex items-center gap-2 mb-2">
                  <Building2 className="w-5 h-5 text-orange-600" />
                  <span className="font-bold">Total Assets</span>
                </div>
                <div className="text-2xl font-bold text-orange-600">156</div>
                <div className="text-sm text-gray-600">Managed assets</div>
              </div>
              <div className="p-4 border-2 border-black bg-green-50 rounded">
                <div className="flex items-center gap-2 mb-2">
                  <CheckCircle className="w-5 h-5 text-green-600" />
                  <span className="font-bold">Patched Today</span>
                </div>
                <div className="text-2xl font-bold text-green-600">12</div>
                <div className="text-sm text-gray-600">Successfully patched</div>
              </div>
              <div className="p-4 border-2 border-black bg-blue-50 rounded">
                <div className="flex items-center gap-2 mb-2">
                  <Clock className="w-5 h-5 text-blue-600" />
                  <span className="font-bold">Next Scan</span>
                </div>
                <div className="text-2xl font-bold text-blue-600">2h</div>
                <div className="text-sm text-gray-600">Scheduled scan</div>
              </div>
            </>
          )}
        </div>
      </div>

      {/* Main Dashboard Content */}
      <GlobalDashboard userRole={userRole} />
    </div>
  );
};

export default RoleBasedDashboard;
