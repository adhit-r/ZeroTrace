import React, { useState, useEffect } from 'react';
import {
  Shield,
  Network,
  Database,
  Lock,
  Container,
  Brain,
  Wifi,
  Eye,
  Globe,
  Zap,
  AlertTriangle,
  BarChart3,
  RefreshCw,
  Download,
  Settings
} from 'lucide-react';
import ComprehensiveSecurityDashboard from '../components/dashboard/ComprehensiveSecurityDashboard';
import SecurityCategoryDashboard from '../components/dashboard/SecurityCategoryDashboard';
import NetworkSecurityDashboard from '../components/dashboard/NetworkSecurityDashboard';
import ComplianceDashboard from '../components/dashboard/ComplianceDashboard';
import { api } from '../services/api';

interface SecurityDashboardProps { }

const SecurityDashboard: React.FC<SecurityDashboardProps> = () => {
  const [selectedView, setSelectedView] = useState<string>('comprehensive');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);


  const dashboardViews = [
    { id: 'comprehensive', name: 'Comprehensive Overview', icon: <BarChart3 className="h-5 w-5" /> },
    { id: 'compliance', name: 'Compliance', icon: <Shield className="h-5 w-5" /> },
    { id: 'network', name: 'Network Security', icon: <Network className="h-5 w-5" /> },
    { id: 'system', name: 'System Vulnerabilities', icon: <AlertTriangle className="h-5 w-5" /> },
    { id: 'auth', name: 'Authentication', icon: <Lock className="h-5 w-5" /> },
    { id: 'database', name: 'Database Security', icon: <Database className="h-5 w-5" /> },
    { id: 'api', name: 'API Security', icon: <Globe className="h-5 w-5" /> },
    { id: 'container', name: 'Container Security', icon: <Container className="h-5 w-5" /> },
    { id: 'ai', name: 'AI/ML Security', icon: <Brain className="h-5 w-5" /> },
    { id: 'iot', name: 'IoT/OT Security', icon: <Wifi className="h-5 w-5" /> },
    { id: 'privacy', name: 'Privacy & Compliance', icon: <Eye className="h-5 w-5" /> },
    { id: 'web3', name: 'Web3 Security', icon: <Zap className="h-5 w-5" /> }
  ];

  const [agents, setAgents] = useState<any[]>([]);

  useEffect(() => {
    const loadDashboardData = async () => {
      try {
        setIsLoading(true);
        setError(null);

        // Fetch agents for dashboards
        const agentsResponse = await api.get('/api/agents/');
        setAgents(agentsResponse.data.data || []);

        setIsLoading(false);
      } catch (err) {
        console.error('Failed to load dashboard data:', err);
        setError('Failed to load security dashboard data');
        setIsLoading(false);
      }
    };

    loadDashboardData();
  }, []);

  const renderDashboard = () => {
    switch (selectedView) {
      case 'comprehensive':
        return <ComprehensiveSecurityDashboard />;
      case 'compliance':
        return <ComplianceDashboard />;
      case 'network':
        return <NetworkSecurityDashboard agents={agents} />;
      default:
        return <SecurityCategoryDashboard {...({ category: selectedView } as any)} />;
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="animate-pulse space-y-6">
            <div className="h-8 bg-gray-200 rounded border-3 border-black"></div>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
              {Array.from({ length: 4 }).map((_, i) => (
                <div key={i} className="h-32 bg-gray-200 rounded border-3 border-black"></div>
              ))}
            </div>
            <div className="h-96 bg-gray-200 rounded border-3 border-black"></div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="p-6 bg-red-50 border-3 border-red-300 rounded shadow-neubrutalist-lg">
            <div className="flex items-center gap-3 mb-4">
              <AlertTriangle className="h-6 w-6 text-red-600" />
              <h2 className="text-xl font-bold text-red-800">Error Loading Dashboard</h2>
            </div>
            <p className="text-red-700 mb-4">{error}</p>
            <button
              onClick={() => window.location.reload()}
              className="px-4 py-2 bg-red-100 text-red-800 border-2 border-red-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-red-200 transition-colors"
            >
              <RefreshCw className="h-4 w-4 mr-2 inline-block" />
              Retry
            </button>
          </div>
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
            <div className="flex items-center gap-4">
              <div className="p-3 bg-blue-100 rounded border-2 border-black">
                <Shield className="h-8 w-8 text-blue-600" />
              </div>
              <div>
                <h1 className="text-3xl font-black uppercase tracking-wider text-black">
                  Security Dashboard
                </h1>
                <p className="text-gray-600">Comprehensive security analysis across all categories</p>
              </div>
            </div>
            <div className="flex items-center gap-4">
              <button className="px-4 py-2 bg-gray-100 text-gray-800 border-2 border-gray-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-gray-200 transition-colors">
                <Settings className="h-4 w-4 mr-2 inline-block" />
                Settings
              </button>
              <button className="px-4 py-2 bg-orange-100 text-orange-800 border-2 border-orange-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-orange-200 transition-colors">
                <Download className="h-4 w-4 mr-2 inline-block" />
                Export
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Navigation Tabs */}
      <div className="bg-white border-b-3 border-black">
        <div className="max-w-7xl mx-auto px-6">
          <div className="flex space-x-1 overflow-x-auto">
            {dashboardViews.map((view) => (
              <button
                key={view.id}
                onClick={() => setSelectedView(view.id)}
                className={`flex items-center gap-2 px-4 py-3 text-sm font-bold uppercase tracking-wider border-b-3 transition-colors whitespace-nowrap ${selectedView === view.id
                  ? 'text-blue-600 border-blue-600 bg-blue-50'
                  : 'text-gray-600 border-transparent hover:text-gray-800 hover:border-gray-300'
                  }`}
              >
                {view.icon}
                {view.name}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto p-6">
        {renderDashboard()}
      </div>
    </div >
  );
};

export default SecurityDashboard;
