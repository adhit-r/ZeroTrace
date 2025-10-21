import React from 'react';
import { 
  Shield, 
  Network, 
  Server, 
  Database, 
  Lock, 
  Eye, 
  AlertTriangle, 
  CheckCircle,
  Zap
} from 'lucide-react';

const SecurityCategoryDashboard: React.FC = () => {
  const categories = [
    {
      name: 'Network Security',
      icon: Network,
      score: 92,
      status: 'good',
      color: 'green',
      issues: 2,
      description: 'Firewall, VPN, and network monitoring'
    },
    {
      name: 'System Security',
      icon: Server,
      score: 78,
      status: 'warning',
      color: 'orange',
      issues: 8,
      description: 'OS patches, system hardening, and access controls'
    },
    {
      name: 'Data Security',
      icon: Database,
      score: 85,
      status: 'good',
      color: 'green',
      issues: 3,
      description: 'Encryption, backup, and data classification'
    },
    {
      name: 'Application Security',
      icon: Shield,
      score: 65,
      status: 'critical',
      color: 'red',
      issues: 15,
      description: 'Code security, dependencies, and API protection'
    },
    {
      name: 'Identity & Access',
      icon: Lock,
      score: 88,
      status: 'good',
      color: 'green',
      issues: 1,
      description: 'Authentication, authorization, and user management'
    },
    {
      name: 'Monitoring & Detection',
      icon: Eye,
      score: 82,
      status: 'good',
      color: 'green',
      issues: 4,
      description: 'Logging, monitoring, and incident response'
    }
  ];


  const getStatusColor = (status: string) => {
    switch (status) {
      case 'good':
        return 'text-green-600 bg-green-100 border-green-300';
      case 'warning':
        return 'text-orange-600 bg-orange-100 border-orange-300';
      case 'critical':
        return 'text-red-600 bg-red-100 border-red-300';
      default:
        return 'text-gray-600 bg-gray-100 border-gray-300';
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
        <h2 className="text-2xl font-black text-black uppercase mb-2">Security Categories</h2>
        <p className="text-gray-600 font-bold">Comprehensive security assessment across all categories</p>
      </div>

      {/* Category Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {categories.map((category, index) => {
          const IconComponent = category.icon;
          return (
            <div key={index} className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
              <div className="flex items-center justify-between mb-4">
                <div className="p-3 bg-gray-100 rounded border-2 border-black">
                  <IconComponent className="h-6 w-6 text-gray-600" />
                </div>
                <div className="text-right">
                  <div className={`text-3xl font-bold ${category.color === 'green' ? 'text-green-600' : category.color === 'orange' ? 'text-orange-600' : 'text-red-600'}`}>
                    {category.score}%
                  </div>
                  <div className="text-sm text-gray-600 uppercase tracking-wider">Score</div>
                </div>
              </div>
              
              <h3 className="text-lg font-black text-black uppercase mb-2">{category.name}</h3>
              <p className="text-sm text-gray-600 mb-4">{category.description}</p>
              
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-bold">Status</span>
                  <div className={`px-3 py-1 rounded border-2 font-bold uppercase text-xs ${getStatusColor(category.status)}`}>
                    {category.status}
                  </div>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-bold">Issues</span>
                  <span className={`text-sm font-bold ${category.issues > 10 ? 'text-red-600' : category.issues > 5 ? 'text-orange-600' : 'text-green-600'}`}>
                    {category.issues}
                  </span>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Detailed Analysis */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Top Issues by Category */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <h3 className="text-xl font-black text-black uppercase mb-4">Top Issues by Category</h3>
          <div className="space-y-4">
            <div className="flex items-center justify-between p-3 bg-red-50 border-2 border-red-300 rounded">
              <div className="flex items-center gap-3">
                <Shield className="h-5 w-5 text-red-600" />
                <div>
                  <p className="font-bold text-red-800">Application Security</p>
                  <p className="text-sm text-red-600">15 critical vulnerabilities</p>
                </div>
              </div>
              <span className="text-sm text-red-600 font-bold">Critical</span>
            </div>
            <div className="flex items-center justify-between p-3 bg-orange-50 border-2 border-orange-300 rounded">
              <div className="flex items-center gap-3">
                <Server className="h-5 w-5 text-orange-600" />
                <div>
                  <p className="font-bold text-orange-800">System Security</p>
                  <p className="text-sm text-orange-600">8 security issues</p>
                </div>
              </div>
              <span className="text-sm text-orange-600 font-bold">Warning</span>
            </div>
            <div className="flex items-center justify-between p-3 bg-green-50 border-2 border-green-300 rounded">
              <div className="flex items-center gap-3">
                <Network className="h-5 w-5 text-green-600" />
                <div>
                  <p className="font-bold text-green-800">Network Security</p>
                  <p className="text-sm text-green-600">2 minor issues</p>
                </div>
              </div>
              <span className="text-sm text-green-600 font-bold">Good</span>
            </div>
          </div>
        </div>

        {/* Security Recommendations */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <h3 className="text-xl font-black text-black uppercase mb-4">Security Recommendations</h3>
          <div className="space-y-3">
            <div className="p-3 bg-blue-50 border-2 border-blue-300 rounded">
              <div className="flex items-start gap-3">
                <Zap className="h-5 w-5 text-blue-600 mt-0.5" />
                <div>
                  <p className="font-bold text-blue-800">Update Dependencies</p>
                  <p className="text-sm text-blue-600">15 outdated packages need updates</p>
                </div>
              </div>
            </div>
            <div className="p-3 bg-orange-50 border-2 border-orange-300 rounded">
              <div className="flex items-start gap-3">
                <AlertTriangle className="h-5 w-5 text-orange-600 mt-0.5" />
                <div>
                  <p className="font-bold text-orange-800">System Hardening</p>
                  <p className="text-sm text-orange-600">Apply security patches and configurations</p>
                </div>
              </div>
            </div>
            <div className="p-3 bg-green-50 border-2 border-green-300 rounded">
              <div className="flex items-start gap-3">
                <CheckCircle className="h-5 w-5 text-green-600 mt-0.5" />
                <div>
                  <p className="font-bold text-green-800">Network Monitoring</p>
                  <p className="text-sm text-green-600">Continue current monitoring practices</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SecurityCategoryDashboard;
