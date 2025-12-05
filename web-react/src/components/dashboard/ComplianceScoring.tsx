import React, { useState, useEffect } from 'react';
import { 
  CheckCircle, 
  Shield, 
  Award, 
  BarChart3, 
  TrendingUp, 
  Eye, 
  Download, 
  RefreshCw, 
  Settings, 
  FileText, 
  Server
} from 'lucide-react';

interface ComplianceFramework {
  id: string;
  name: string;
  version: string;
  description: string;
  category: 'security' | 'privacy' | 'governance' | 'operational';
  score: number;
  maxScore: number;
  status: 'compliant' | 'non_compliant' | 'partial' | 'not_assessed';
  lastAssessed: string;
  nextAssessment: string;
  controls: ComplianceControl[];
  requirements: ComplianceRequirement[];
}

interface ComplianceControl {
  id: string;
  name: string;
  description: string;
  category: string;
  status: 'compliant' | 'non_compliant' | 'not_applicable' | 'not_assessed';
  evidence: string[];
  lastAssessed: string;
  assessor: string;
  notes: string;
  remediation: string[];
  priority: 'high' | 'medium' | 'low';
}

interface ComplianceRequirement {
  id: string;
  title: string;
  description: string;
  category: string;
  status: 'met' | 'not_met' | 'partial' | 'not_applicable';
  evidence: string[];
  lastVerified: string;
  verifiedBy: string;
  notes: string;
  actionItems: string[];
}

interface ComplianceAsset {
  id: string;
  name: string;
  type: string;
  framework: string;
  score: number;
  status: 'compliant' | 'non_compliant' | 'partial';
  lastAssessed: string;
  assessor: string;
  issues: number;
  controls: number;
  compliantControls: number;
}

interface ComplianceData {
  frameworks: ComplianceFramework[];
  assets: ComplianceAsset[];
  summary: {
    overallScore: number;
    totalFrameworks: number;
    compliantFrameworks: number;
    totalAssets: number;
    compliantAssets: number;
    totalControls: number;
    compliantControls: number;
    totalRequirements: number;
    metRequirements: number;
    upcomingAssessments: number;
    overdueAssessments: number;
  };
  trends: Array<{
    date: string;
    score: number;
    frameworks: number;
    assets: number;
  }>;
}

interface ComplianceScoringProps {
  className?: string;
}

const ComplianceScoring: React.FC<ComplianceScoringProps> = ({ className = '' }) => {
  const [data, setData] = useState<ComplianceData | null>(null);
  const [loading, setLoading] = useState(true);
  const [viewMode, setViewMode] = useState<'overview' | 'frameworks' | 'assets' | 'controls' | 'trends'>('overview');
  const [timeRange, setTimeRange] = useState<'30d' | '90d' | '1y'>('90d');

  useEffect(() => {
    const fetchComplianceData = async () => {
      setLoading(true);
      try {
        // Fetch agent data to build compliance information
        const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/agents/`);
        if (!response.ok) {
          throw new Error('Failed to fetch agent data');
        }
        
        const result = await response.json();
        const agents = result.data || [];
        
        // Generate compliance frameworks
        const frameworks: ComplianceFramework[] = [
          {
            id: 'iso27001',
            name: 'ISO 27001',
            version: '2022',
            description: 'Information Security Management System',
            category: 'security',
            score: 85,
            maxScore: 100,
            status: 'compliant',
            lastAssessed: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString(),
            nextAssessment: new Date(Date.now() + 335 * 24 * 60 * 60 * 1000).toISOString(),
            controls: generateControls('ISO 27001'),
            requirements: generateRequirements('ISO 27001')
          },
          {
            id: 'soc2',
            name: 'SOC 2',
            version: '2017',
            description: 'Service Organization Control 2',
            category: 'security',
            score: 78,
            maxScore: 100,
            status: 'partial',
            lastAssessed: new Date(Date.now() - 15 * 24 * 60 * 60 * 1000).toISOString(),
            nextAssessment: new Date(Date.now() + 350 * 24 * 60 * 60 * 1000).toISOString(),
            controls: generateControls('SOC 2'),
            requirements: generateRequirements('SOC 2')
          },
          {
            id: 'gdpr',
            name: 'GDPR',
            version: '2018',
            description: 'General Data Protection Regulation',
            category: 'privacy',
            score: 92,
            maxScore: 100,
            status: 'compliant',
            lastAssessed: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
            nextAssessment: new Date(Date.now() + 358 * 24 * 60 * 60 * 1000).toISOString(),
            controls: generateControls('GDPR'),
            requirements: generateRequirements('GDPR')
          },
          {
            id: 'pci-dss',
            name: 'PCI DSS',
            version: '4.0',
            description: 'Payment Card Industry Data Security Standard',
            category: 'security',
            score: 65,
            maxScore: 100,
            status: 'non_compliant',
            lastAssessed: new Date(Date.now() - 45 * 24 * 60 * 60 * 1000).toISOString(),
            nextAssessment: new Date(Date.now() + 320 * 24 * 60 * 60 * 1000).toISOString(),
            controls: generateControls('PCI DSS'),
            requirements: generateRequirements('PCI DSS')
          }
        ];
        
        // Generate compliance assets
        const assets: ComplianceAsset[] = agents.map((agent: any, index: number) => ({
          id: agent.id,
          name: agent.hostname || `Asset-${index + 1}`,
          type: agent.metadata?.type || 'workstation',
          framework: frameworks[Math.floor(Math.random() * frameworks.length)].name,
          score: Math.floor(Math.random() * 40) + 60,
          status: Math.random() > 0.3 ? 'compliant' : Math.random() > 0.5 ? 'partial' : 'non_compliant',
          lastAssessed: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
          assessor: 'Security Team',
          issues: Math.floor(Math.random() * 10),
          controls: Math.floor(Math.random() * 20) + 10,
          compliantControls: Math.floor(Math.random() * 15) + 5
        }));
        
        // Generate trends
        const trends = generateTrends(timeRange);
        
        const summary = {
          overallScore: Math.floor(frameworks.reduce((sum, f) => sum + f.score, 0) / frameworks.length),
          totalFrameworks: frameworks.length,
          compliantFrameworks: frameworks.filter(f => f.status === 'compliant').length,
          totalAssets: assets.length,
          compliantAssets: assets.filter(a => a.status === 'compliant').length,
          totalControls: frameworks.reduce((sum, f) => sum + f.controls.length, 0),
          compliantControls: frameworks.reduce((sum, f) => sum + f.controls.filter(c => c.status === 'compliant').length, 0),
          totalRequirements: frameworks.reduce((sum, f) => sum + f.requirements.length, 0),
          metRequirements: frameworks.reduce((sum, f) => sum + f.requirements.filter(r => r.status === 'met').length, 0),
          upcomingAssessments: frameworks.filter(f => new Date(f.nextAssessment) <= new Date(Date.now() + 30 * 24 * 60 * 60 * 1000)).length,
          overdueAssessments: frameworks.filter(f => new Date(f.nextAssessment) < new Date()).length
        };
        
        setData({
          frameworks,
          assets,
          summary,
          trends
        });
      } catch (error) {
        console.error('Failed to fetch compliance data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchComplianceData();
  }, [timeRange]);

  const generateControls = (framework: string): ComplianceControl[] => {
    const controlTemplates = [
      { name: 'Access Control', category: 'Security', priority: 'high' as const },
      { name: 'Data Encryption', category: 'Security', priority: 'high' as const },
      { name: 'Incident Response', category: 'Security', priority: 'medium' as const },
      { name: 'Vulnerability Management', category: 'Security', priority: 'high' as const },
      { name: 'Security Awareness', category: 'Governance', priority: 'medium' as const },
      { name: 'Risk Assessment', category: 'Governance', priority: 'high' as const },
      { name: 'Data Backup', category: 'Operational', priority: 'medium' as const },
      { name: 'Change Management', category: 'Operational', priority: 'low' as const }
    ];
    
    return controlTemplates.map((template, index) => ({
      id: `${framework.toLowerCase()}-control-${index}`,
      name: template.name,
      description: `${template.name} control for ${framework}`,
      category: template.category,
      status: Math.random() > 0.2 ? 'compliant' : Math.random() > 0.5 ? 'non_compliant' : 'not_assessed',
      evidence: ['Policy document', 'Technical implementation', 'Audit logs'],
      lastAssessed: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
      assessor: 'Compliance Team',
      notes: `Assessment notes for ${template.name}`,
      remediation: ['Update policy', 'Implement controls', 'Train staff'],
      priority: template.priority
    }));
  };

  const generateRequirements = (framework: string): ComplianceRequirement[] => {
    const requirementTemplates = [
      { title: 'Data Protection', category: 'Privacy' },
      { title: 'Access Management', category: 'Security' },
      { title: 'Incident Handling', category: 'Security' },
      { title: 'Data Retention', category: 'Privacy' },
      { title: 'Audit Logging', category: 'Governance' },
      { title: 'Business Continuity', category: 'Operational' }
    ];
    
    return requirementTemplates.map((template, index) => ({
      id: `${framework.toLowerCase()}-req-${index}`,
      title: template.title,
      description: `${template.title} requirement for ${framework}`,
      category: template.category,
      status: Math.random() > 0.3 ? 'met' : Math.random() > 0.5 ? 'partial' : 'not_met',
      evidence: ['Documentation', 'Implementation', 'Testing'],
      lastVerified: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
      verifiedBy: 'Compliance Team',
      notes: `Verification notes for ${template.title}`,
      actionItems: ['Review documentation', 'Update procedures', 'Conduct training']
    }));
  };

  const generateTrends = (range: string): Array<{ date: string; score: number; frameworks: number; assets: number }> => {
    const days = range === '30d' ? 30 : range === '90d' ? 90 : 365;
    const trends = [];
    
    for (let i = days - 1; i >= 0; i--) {
      const date = new Date(Date.now() - i * 24 * 60 * 60 * 1000);
      trends.push({
        date: date.toISOString().split('T')[0],
        score: Math.floor(Math.random() * 20) + 70,
        frameworks: Math.floor(Math.random() * 5) + 3,
        assets: Math.floor(Math.random() * 20) + 10
      });
    }
    
    return trends;
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'compliant': return 'text-green-600 bg-green-100';
      case 'non_compliant': return 'text-red-600 bg-red-100';
      case 'partial': return 'text-yellow-600 bg-yellow-100';
      case 'not_assessed': return 'text-gray-600 bg-gray-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  const getScoreColor = (score: number) => {
    if (score >= 90) return 'text-green-600';
    if (score >= 80) return 'text-blue-600';
    if (score >= 70) return 'text-yellow-600';
    if (score >= 60) return 'text-orange-600';
    return 'text-red-600';
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  if (loading) {
    return (
      <div className={`p-6 ${className}`}>
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className={`p-6 ${className}`}>
        <div className="text-center text-gray-500">
          <Shield className="h-12 w-12 mx-auto mb-4" />
          <p>No compliance data available</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Compliance Scoring</h2>
          <p className="text-gray-600">Comprehensive compliance framework assessment</p>
        </div>
        <div className="flex items-center gap-2">
          <select
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value as any)}
            className="px-4 py-2 border-2 border-black rounded focus:outline-none"
          >
            <option value="30d">Last 30 Days</option>
            <option value="90d">Last 90 Days</option>
            <option value="1y">Last Year</option>
          </select>
          <button className="p-2 bg-blue-600 text-white rounded border-2 border-blue-700 hover:bg-blue-700 transition-colors">
            <RefreshCw className="h-4 w-4" />
          </button>
        </div>
      </div>

      {/* View Mode Tabs */}
      <div className="flex space-x-1 bg-gray-100 p-1 rounded-lg">
        {[
          { id: 'overview', label: 'Overview', icon: BarChart3 },
          { id: 'frameworks', label: 'Frameworks', icon: Shield },
          { id: 'assets', label: 'Assets', icon: Server },
          { id: 'controls', label: 'Controls', icon: CheckCircle },
          { id: 'trends', label: 'Trends', icon: TrendingUp }
        ].map((tab) => (
          <button
            key={tab.id}
            onClick={() => setViewMode(tab.id as any)}
            className={`flex items-center gap-2 px-4 py-2 rounded-md transition-colors ${
              viewMode === tab.id
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            <tab.icon className="h-4 w-4" />
            {tab.label}
          </button>
        ))}
      </div>

      {/* Overview Mode */}
      {viewMode === 'overview' && (
        <div className="space-y-6">
          {/* Summary Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Overall Score</p>
                  <p className={`text-3xl font-bold ${getScoreColor(data.summary.overallScore)}`}>
                    {data.summary.overallScore}%
                  </p>
                  <p className="text-sm text-gray-600">{data.summary.compliantFrameworks}/{data.summary.totalFrameworks} frameworks</p>
                </div>
                <Award className="h-12 w-12 text-blue-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Compliant Assets</p>
                  <p className="text-3xl font-bold text-green-600">{data.summary.compliantAssets}</p>
                  <p className="text-sm text-gray-600">of {data.summary.totalAssets} total</p>
                </div>
                <CheckCircle className="h-12 w-12 text-green-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Controls</p>
                  <p className="text-3xl font-bold text-blue-600">{data.summary.compliantControls}</p>
                  <p className="text-sm text-gray-600">of {data.summary.totalControls} total</p>
                </div>
                <Shield className="h-12 w-12 text-blue-600" />
              </div>
            </div>
            
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600">Requirements</p>
                  <p className="text-3xl font-bold text-purple-600">{data.summary.metRequirements}</p>
                  <p className="text-sm text-gray-600">of {data.summary.totalRequirements} total</p>
                </div>
                <FileText className="h-12 w-12 text-purple-600" />
              </div>
            </div>
          </div>

          {/* Framework Overview */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <h3 className="text-lg font-bold mb-4">Framework Scores</h3>
              <div className="space-y-4">
                {data.frameworks.map((framework) => (
                  <div key={framework.id} className="flex items-center justify-between p-3 bg-gray-50 rounded border-2 border-gray-200">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-blue-100 rounded">
                        <Shield className="h-4 w-4 text-blue-600" />
                      </div>
                      <div>
                        <p className="font-medium">{framework.name}</p>
                        <p className="text-sm text-gray-600">{framework.description}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className={`text-2xl font-bold ${getScoreColor(framework.score)}`}>
                        {framework.score}%
                      </p>
                      <span className={`px-2 py-1 rounded text-xs font-bold ${getStatusColor(framework.status)}`}>
                        {framework.status.replace('_', ' ')}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
              <h3 className="text-lg font-bold mb-4">Assessment Status</h3>
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Upcoming Assessments</span>
                  <span className="font-bold text-orange-600">{data.summary.upcomingAssessments}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Overdue Assessments</span>
                  <span className="font-bold text-red-600">{data.summary.overdueAssessments}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Total Frameworks</span>
                  <span className="font-bold text-blue-600">{data.summary.totalFrameworks}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Compliant Frameworks</span>
                  <span className="font-bold text-green-600">{data.summary.compliantFrameworks}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Frameworks Mode */}
      {viewMode === 'frameworks' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {data.frameworks.map((framework) => (
              <div key={framework.id} className="p-6 bg-white border-3 border-black rounded shadow-lg">
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <h3 className="text-lg font-bold">{framework.name} {framework.version}</h3>
                    <p className="text-gray-600">{framework.description}</p>
                  </div>
                  <div className="text-right">
                    <p className={`text-3xl font-bold ${getScoreColor(framework.score)}`}>
                      {framework.score}%
                    </p>
                    <span className={`px-2 py-1 rounded text-xs font-bold ${getStatusColor(framework.status)}`}>
                      {framework.status.replace('_', ' ')}
                    </span>
                  </div>
                </div>
                
                <div className="space-y-3">
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Controls:</span>
                    <span className="font-bold">{framework.controls.length}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Requirements:</span>
                    <span className="font-bold">{framework.requirements.length}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Last Assessed:</span>
                    <span className="font-bold">{formatDate(framework.lastAssessed)}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Next Assessment:</span>
                    <span className="font-bold">{formatDate(framework.nextAssessment)}</span>
                  </div>
                </div>
                
                <div className="mt-4 flex items-center gap-2">
                  <button className="flex-1 px-3 py-2 bg-blue-100 text-blue-800 rounded border-2 border-blue-300 hover:bg-blue-200 transition-colors">
                    <Eye className="h-4 w-4 inline mr-1" />
                    View Details
                  </button>
                  <button className="px-3 py-2 bg-gray-100 text-gray-800 rounded border-2 border-gray-300 hover:bg-gray-200 transition-colors">
                    <Download className="h-4 w-4" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Assets Mode */}
      {viewMode === 'assets' && (
        <div className="space-y-6">
          <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
            <h3 className="text-lg font-bold mb-4">Asset Compliance Status</h3>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Asset</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Framework</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Score</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Issues</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Last Assessed</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {data.assets.map((asset) => (
                    <tr key={asset.id} className="hover:bg-gray-50">
                      <td className="px-4 py-3">
                        <div>
                          <div className="font-medium text-gray-900">{asset.name}</div>
                          <div className="text-sm text-gray-500">{asset.type}</div>
                        </div>
                      </td>
                      <td className="px-4 py-3 text-sm text-gray-900">{asset.framework}</td>
                      <td className="px-4 py-3">
                        <span className={`text-lg font-bold ${getScoreColor(asset.score)}`}>
                          {asset.score}%
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        <span className={`px-2 py-1 rounded text-xs font-bold ${getStatusColor(asset.status)}`}>
                          {asset.status.replace('_', ' ')}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-sm text-gray-900">{asset.issues}</td>
                      <td className="px-4 py-3 text-sm text-gray-900">{formatDate(asset.lastAssessed)}</td>
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <button className="p-1 text-blue-600 hover:bg-blue-100 rounded">
                            <Eye className="h-4 w-4" />
                          </button>
                          <button className="p-1 text-gray-600 hover:bg-gray-100 rounded">
                            <Settings className="h-4 w-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {/* Controls Mode */}
      {viewMode === 'controls' && (
        <div className="space-y-6">
          <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
            <h3 className="text-lg font-bold mb-4">Compliance Controls</h3>
            <div className="space-y-4">
              {data.frameworks.map((framework) => (
                <div key={framework.id} className="p-4 bg-gray-50 rounded border-2 border-gray-200">
                  <div className="flex items-center justify-between mb-3">
                    <h4 className="font-bold">{framework.name}</h4>
                    <span className="text-sm text-gray-600">{framework.controls.length} controls</span>
                  </div>
                  
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                    {framework.controls.map((control) => (
                      <div key={control.id} className="p-3 bg-white rounded border border-gray-300">
                        <div className="flex items-start justify-between mb-2">
                          <div>
                            <p className="font-medium">{control.name}</p>
                            <p className="text-sm text-gray-600">{control.description}</p>
                          </div>
                          <div className="flex items-center gap-2">
                            <span className={`px-2 py-1 rounded text-xs font-bold ${getStatusColor(control.status)}`}>
                              {control.status.replace('_', ' ')}
                            </span>
                            <span className={`px-2 py-1 rounded text-xs font-bold ${
                              control.priority === 'high' ? 'bg-red-100 text-red-800' :
                              control.priority === 'medium' ? 'bg-yellow-100 text-yellow-800' :
                              'bg-green-100 text-green-800'
                            }`}>
                              {control.priority}
                            </span>
                          </div>
                        </div>
                        
                        <div className="text-xs text-gray-600">
                          <p>Last Assessed: {formatDate(control.lastAssessed)}</p>
                          <p>Assessor: {control.assessor}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Trends Mode */}
      {viewMode === 'trends' && (
        <div className="space-y-6">
          <div className="p-6 bg-white border-3 border-black rounded shadow-lg">
            <h3 className="text-lg font-bold mb-4">Compliance Trends</h3>
            <div className="space-y-4">
              {data.trends.slice(-7).map((trend, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded border-2 border-gray-200">
                  <div className="flex items-center gap-4">
                    <div className="text-sm font-medium">{formatDate(trend.date)}</div>
                    <div className="flex items-center gap-2">
                      <span className="text-sm text-gray-600">Score:</span>
                      <span className="font-bold">{trend.score}%</span>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <div className="flex items-center gap-1">
                      <Shield className="h-4 w-4 text-blue-600" />
                      <span className="text-sm">{trend.frameworks}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Server className="h-4 w-4 text-green-600" />
                      <span className="text-sm">{trend.assets}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default ComplianceScoring;
