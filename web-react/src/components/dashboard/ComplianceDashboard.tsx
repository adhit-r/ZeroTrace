import React, { useState, useEffect } from 'react';
import {
  Shield,
  AlertTriangle,
  CheckCircle,
  ChevronDown,
  ChevronUp,
  Lock,
  Users,
  Globe,
  Database,
  Server,
  AlertCircle,
  XCircle,
  BarChart3,
  RefreshCw,
  Filter,
  Download
} from 'lucide-react';
import { Bar, Doughnut, Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend
);

interface ComplianceFramework {
  id: string;
  name: string;
  description: string;
  icon: React.ReactNode;
  color: string;
  score: number;
  total: number;
  passed: number;
  failed: number;
  status: 'compliant' | 'non_compliant' | 'partial';
  lastAssessment: string;
  nextAssessment: string;
  requirements: ComplianceRequirement[];
  trends: {
    score: number[];
    passed: number[];
    failed: number[];
  };
}

interface ComplianceRequirement {
  id: string;
  title: string;
  description: string;
  category: string;
  priority: 'critical' | 'high' | 'medium' | 'low';
  status: 'compliant' | 'non_compliant' | 'partial' | 'not_assessed';
  evidence: string[];
  remediation: string;
  references: string[];
  lastChecked: string;
}

interface ComplianceGap {
  framework: string;
  requirement: string;
  category: string;
  priority: string;
  status: string;
  gap: string;
  remediation: string;
  effort: 'low' | 'medium' | 'high';
  timeline: string;
}

interface ComplianceDashboardProps {
  className?: string;
}

const ComplianceDashboard: React.FC<ComplianceDashboardProps> = ({
  className = ''
}) => {
  const [frameworks, setFrameworks] = useState<ComplianceFramework[]>([]);
  const [selectedFramework, setSelectedFramework] = useState<string>('');
  const [gaps, setGaps] = useState<ComplianceGap[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [expandedSections, setExpandedSections] = useState<Set<string>>(new Set(['overview', 'frameworks']));

  const complianceFrameworks = [
    { id: 'cis', name: 'CIS Benchmarks', icon: <Shield className="h-5 w-5" />, color: 'blue' },
    { id: 'pci', name: 'PCI-DSS', icon: <Lock className="h-5 w-5" />, color: 'green' },
    { id: 'hipaa', name: 'HIPAA', icon: <Users className="h-5 w-5" />, color: 'purple' },
    { id: 'gdpr', name: 'GDPR', icon: <Globe className="h-5 w-5" />, color: 'orange' },
    { id: 'soc2', name: 'SOC 2', icon: <Database className="h-5 w-5" />, color: 'cyan' },
    { id: 'iso27001', name: 'ISO 27001', icon: <Server className="h-5 w-5" />, color: 'indigo' }
  ];

  useEffect(() => {
    const loadComplianceData = async () => {
      try {
        setIsLoading(true);
        setError(null);

        // Fetch compliance status
        const response = await fetch('/api/v2/compliance/status');
        if (!response.ok) {
          throw new Error('Failed to fetch compliance data');
        }

        const complianceData = await response.json();

        // Transform API data
        const transformedFrameworks: ComplianceFramework[] = Object.entries(complianceData.frameworks || {}).map(([frameworkId, data]: [string, any]) => {
          const frameworkConfig = complianceFrameworks.find(f => f.id === frameworkId);
          return {
            id: frameworkId,
            name: frameworkConfig?.name || frameworkId.toUpperCase(),
            description: data.description || `Compliance framework for ${frameworkId.toUpperCase()}`,
            icon: frameworkConfig?.icon || <Shield className="h-5 w-5" />,
            color: frameworkConfig?.color || 'gray',
            score: data.score || 0,
            total: data.total || 0,
            passed: data.passed || 0,
            failed: data.failed || 0,
            status: data.score >= 80 ? 'compliant' : data.score >= 60 ? 'partial' : 'non_compliant',
            lastAssessment: data.last_assessment || new Date().toISOString(),
            nextAssessment: data.next_assessment || new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(),
            requirements: data.requirements?.map((req: any) => ({
              id: req.id,
              title: req.title,
              description: req.description,
              category: req.category,
              priority: req.priority,
              status: req.status,
              evidence: req.evidence || [],
              remediation: req.remediation,
              references: req.references || [],
              lastChecked: req.last_checked || new Date().toISOString()
            })) || [],
            trends: {
              score: data.trends?.score || [0, 0, 0, 0, 0, 0, 0, data.score || 0],
              passed: data.trends?.passed || [0, 0, 0, 0, 0, 0, 0, data.passed || 0],
              failed: data.trends?.failed || [0, 0, 0, 0, 0, 0, 0, data.failed || 0]
            }
          };
        });

        setFrameworks(transformedFrameworks);

        // Fetch gap analysis
        const gapsResponse = await fetch('/api/v2/compliance/gaps');
        if (gapsResponse.ok) {
          const gapsData = await gapsResponse.json();
          setGaps(gapsData.gaps || []);
        }

      } catch (err) {
        console.error('Failed to load compliance data:', err);
        setError('Failed to load compliance dashboard data');
      } finally {
        setIsLoading(false);
      }
    };

    loadComplianceData();
  }, []);

  const toggleSection = (section: string) => {
    const newExpanded = new Set(expandedSections);
    if (newExpanded.has(section)) {
      newExpanded.delete(section);
    } else {
      newExpanded.add(section);
    }
    setExpandedSections(newExpanded);
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'compliant':
        return <CheckCircle className="h-5 w-5 text-green-600" />;
      case 'partial':
        return <AlertCircle className="h-5 w-5 text-yellow-600" />;
      case 'non_compliant':
        return <XCircle className="h-5 w-5 text-red-600" />;
      default:
        return <AlertCircle className="h-5 w-5 text-gray-600" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'compliant':
        return 'bg-green-100 text-green-800 border-green-300';
      case 'partial':
        return 'bg-yellow-100 text-yellow-800 border-yellow-300';
      case 'non_compliant':
        return 'bg-red-100 text-red-800 border-red-300';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-300';
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
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

  const getEffortColor = (effort: string) => {
    switch (effort) {
      case 'low':
        return 'bg-green-100 text-green-800 border-green-300';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800 border-yellow-300';
      case 'high':
        return 'bg-red-100 text-red-800 border-red-300';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-300';
    }
  };

  if (isLoading) {
    return (
      <div className={`p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal ${className}`}>
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded border-2 border-black"></div>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="h-24 bg-gray-200 rounded border-2 border-black"></div>
            ))}
          </div>
          <div className="h-96 bg-gray-200 rounded border-2 border-black"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal ${className}`}>
        <div className="text-center py-8">
          <AlertTriangle className="h-12 w-12 text-red-500 mx-auto mb-4" />
          <h3 className="text-lg font-bold text-red-800 mb-2">Error Loading Compliance Data</h3>
          <p className="text-red-600 mb-4">{error}</p>
          <button
            onClick={() => window.location.reload()}
            className="px-4 py-2 bg-red-100 text-red-800 border-2 border-red-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-red-200 transition-colors"
          >
            <RefreshCw className="h-4 w-4 mr-2 inline-block" />
            Retry
          </button>
        </div>
      </div>
    );
  }

  const totalFrameworks = frameworks.length;
  const compliantFrameworks = frameworks.filter(f => f.status === 'compliant').length;
  const partialFrameworks = frameworks.filter(f => f.status === 'partial').length;
  const nonCompliantFrameworks = frameworks.filter(f => f.status === 'non_compliant').length;
  const avgScore = frameworks.length > 0 ? frameworks.reduce((acc, f) => acc + f.score, 0) / frameworks.length : 0;

  const frameworkData = {
    labels: frameworks.map(f => f.name),
    datasets: [
      {
        label: 'Compliance Score',
        data: frameworks.map(f => f.score),
        backgroundColor: frameworks.map(f => {
          switch (f.status) {
            case 'compliant': return 'rgba(34, 197, 94, 0.8)';
            case 'partial': return 'rgba(245, 158, 11, 0.8)';
            case 'non_compliant': return 'rgba(239, 68, 68, 0.8)';
            default: return 'rgba(107, 114, 128, 0.8)';
          }
        }),
        borderColor: frameworks.map(_f => 'rgba(0, 0, 0, 1)'),
        borderWidth: 3
      }
    ]
  };

  const statusData = {
    labels: ['Compliant', 'Partial', 'Non-Compliant'],
    datasets: [
      {
        label: 'Frameworks',
        data: [compliantFrameworks, partialFrameworks, nonCompliantFrameworks],
        backgroundColor: [
          'rgba(34, 197, 94, 0.8)',
          'rgba(245, 158, 11, 0.8)',
          'rgba(239, 68, 68, 0.8)'
        ],
        borderColor: [
          'rgba(0, 0, 0, 1)',
          'rgba(0, 0, 0, 1)',
          'rgba(0, 0, 0, 1)'
        ],
        borderWidth: 3
      }
    ]
  };

  const selectedFrameworkData = frameworks.find(f => f.id === selectedFramework);

  return (
    <div className={`bg-white border-3 border-black rounded-lg shadow-neo-brutal ${className}`}>
      {/* Header */}
      <div className="p-6 border-b-3 border-black">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="p-3 bg-green-100 rounded border-2 border-black">
              <Shield className="h-8 w-8 text-green-600" />
            </div>
            <div>
              <h1 className="text-3xl font-black uppercase tracking-wider text-black">
                Compliance Dashboard
              </h1>
              <p className="text-gray-600">Compliance framework assessment and gap analysis</p>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <button className="px-4 py-2 bg-gray-100 text-gray-800 border-2 border-gray-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-gray-200 transition-colors">
              <Filter className="h-4 w-4 mr-2 inline-block" />
              Filter
            </button>
            <button className="px-4 py-2 bg-orange-100 text-orange-800 border-2 border-orange-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-orange-200 transition-colors">
              <Download className="h-4 w-4 mr-2 inline-block" />
              Export
            </button>
          </div>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="p-6 border-b-3 border-black">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          {/* Total Frameworks */}
          <div className="p-4 bg-gray-50 border-2 border-black rounded">
            <div className="flex items-center justify-between mb-2">
              <div className="text-3xl font-bold text-black">{totalFrameworks}</div>
              <div className="p-2 bg-blue-100 rounded border border-black">
                <Shield className="h-6 w-6 text-gray-800" />
              </div>
            </div>
            <div className="text-sm text-gray-600 uppercase tracking-wider">Total Frameworks</div>
          </div>

          {/* Average Score */}
          <div className="p-4 bg-gray-50 border-2 border-black rounded">
            <div className="flex items-center justify-between mb-2">
              <div className="text-3xl font-bold text-black">{avgScore.toFixed(1)}%</div>
              <div className="p-2 bg-green-100 rounded border border-black">
                <BarChart3 className="h-6 w-6 text-gray-800" />
              </div>
            </div>
            <div className="text-sm text-gray-600 uppercase tracking-wider">Average Score</div>
          </div>

          {/* Compliant Frameworks */}
          <div className="p-4 bg-gray-50 border-2 border-black rounded">
            <div className="flex items-center justify-between mb-2">
              <div className="text-3xl font-bold text-green-600">{compliantFrameworks}</div>
              <div className="p-2 bg-green-100 rounded border border-black">
                <CheckCircle className="h-6 w-6 text-gray-800" />
              </div>
            </div>
            <div className="text-sm text-gray-600 uppercase tracking-wider">Compliant</div>
          </div>

          {/* Non-Compliant Frameworks */}
          <div className="p-4 bg-gray-50 border-2 border-black rounded">
            <div className="flex items-center justify-between mb-2">
              <div className="text-3xl font-bold text-red-600">{nonCompliantFrameworks}</div>
              <div className="p-2 bg-red-100 rounded border border-black">
                <XCircle className="h-6 w-6 text-gray-800" />
              </div>
            </div>
            <div className="text-sm text-gray-600 uppercase tracking-wider">Non-Compliant</div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="p-6 space-y-6">
        {/* Framework Selection */}
        <div className="border-2 border-black rounded">
          <div className="p-4 bg-gray-50 border-b-2 border-black">
            <h3 className="text-lg font-bold uppercase tracking-wider text-black">Compliance Frameworks</h3>
          </div>
          <div className="p-6">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {frameworks.map((framework) => (
                <button
                  key={framework.id}
                  onClick={() => setSelectedFramework(selectedFramework === framework.id ? '' : framework.id)}
                  className={`p-4 border-2 border-black rounded hover:bg-gray-50 transition-colors ${selectedFramework === framework.id ? 'bg-blue-50 border-blue-300' : 'bg-white'
                    }`}
                >
                  <div className="flex items-center gap-3 mb-3">
                    <div className={`p-2 bg-${framework.color}-100 rounded border border-black`}>
                      {framework.icon}
                    </div>
                    <div className="flex-1 text-left">
                      <div className="font-bold text-black">{framework.name}</div>
                      <div className="text-sm text-gray-600">{framework.description}</div>
                    </div>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      {getStatusIcon(framework.status)}
                      <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider border-2 ${getStatusColor(framework.status)}`}>
                        {framework.status}
                      </span>
                    </div>
                    <div className="text-right">
                      <div className="text-lg font-bold text-black">{framework.score.toFixed(1)}%</div>
                      <div className="text-xs text-gray-600">{framework.passed}/{framework.total} passed</div>
                    </div>
                  </div>
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Selected Framework Details */}
        {selectedFrameworkData && (
          <div className="border-2 border-black rounded">
            <div className="p-4 bg-gray-50 border-b-2 border-black">
              <h3 className="text-lg font-bold uppercase tracking-wider text-black">
                {selectedFrameworkData.name} Details
              </h3>
            </div>
            <div className="p-6">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Requirements */}
                <div>
                  <h4 className="text-md font-bold uppercase tracking-wider text-black mb-4">Requirements</h4>
                  <div className="space-y-3">
                    {selectedFrameworkData.requirements.map((req) => (
                      <div key={req.id} className="p-3 bg-gray-50 border-2 border-black rounded">
                        <div className="flex items-center justify-between mb-2">
                          <div className="font-bold text-black">{req.title}</div>
                          <div className="flex items-center gap-2">
                            <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider border-2 ${getStatusColor(req.status)}`}>
                              {req.status}
                            </span>
                            <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider border-2 ${getPriorityColor(req.priority)}`}>
                              {req.priority}
                            </span>
                          </div>
                        </div>
                        <div className="text-sm text-gray-600 mb-2">{req.description}</div>
                        <div className="text-xs text-gray-500">
                          Category: {req.category} • Last checked: {new Date(req.lastChecked).toLocaleDateString()}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Trends */}
                <div>
                  <h4 className="text-md font-bold uppercase tracking-wider text-black mb-4">Trends</h4>
                  <div className="h-64">
                    <Line
                      data={{
                        labels: ['7d ago', '6d ago', '5d ago', '4d ago', '3d ago', '2d ago', '1d ago', 'Today'],
                        datasets: [
                          {
                            label: 'Score',
                            data: selectedFrameworkData.trends.score,
                            fill: false,
                            backgroundColor: 'rgb(34, 197, 94)',
                            borderColor: 'rgba(34, 197, 94, 0.8)',
                            borderWidth: 3,
                            tension: 0.4
                          },
                          {
                            label: 'Passed',
                            data: selectedFrameworkData.trends.passed,
                            fill: false,
                            backgroundColor: 'rgb(59, 130, 246)',
                            borderColor: 'rgba(59, 130, 246, 0.8)',
                            borderWidth: 3,
                            tension: 0.4
                          },
                          {
                            label: 'Failed',
                            data: selectedFrameworkData.trends.failed,
                            fill: false,
                            backgroundColor: 'rgb(239, 68, 68)',
                            borderColor: 'rgba(239, 68, 68, 0.8)',
                            borderWidth: 3,
                            tension: 0.4
                          }
                        ]
                      }}
                      options={{
                        responsive: true,
                        maintainAspectRatio: false,
                        scales: {
                          y: {
                            beginAtZero: true
                          }
                        }
                      }}
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Gap Analysis */}
        <div className="border-2 border-black rounded">
          <button
            onClick={() => toggleSection('gaps')}
            className="w-full p-4 bg-gray-50 border-b-2 border-black rounded-t flex items-center justify-between hover:bg-gray-100 transition-colors"
          >
            <h3 className="text-lg font-bold uppercase tracking-wider text-black">Gap Analysis</h3>
            {expandedSections.has('gaps') ?
              <ChevronUp className="h-5 w-5 text-gray-600" /> :
              <ChevronDown className="h-5 w-5 text-gray-600" />
            }
          </button>
          {expandedSections.has('gaps') && (
            <div className="p-6">
              <div className="space-y-4">
                {gaps.map((gap, index) => (
                  <div key={index} className="p-4 bg-gray-50 border-2 border-black rounded">
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <div className="font-bold text-black">{gap.framework}</div>
                        <div className="text-sm text-gray-600">{gap.requirement}</div>
                      </div>
                      <div className="flex items-center gap-2">
                        <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider border-2 ${getEffortColor(gap.effort)}`}>
                          {gap.effort} effort
                        </span>
                        <span className={`px-2 py-1 rounded text-xs font-bold uppercase tracking-wider border-2 ${getPriorityColor(gap.priority)}`}>
                          {gap.priority}
                        </span>
                      </div>
                    </div>
                    <div className="text-sm text-gray-600 mb-2">{gap.gap}</div>
                    <div className="text-sm text-gray-800 mb-2">
                      <strong>Remediation:</strong> {gap.remediation}
                    </div>
                    <div className="text-xs text-gray-500">
                      Timeline: {gap.timeline} • Category: {gap.category}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Charts */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Framework Scores */}
          <div className="border-2 border-black rounded">
            <div className="p-4 bg-gray-50 border-b-2 border-black">
              <h3 className="text-lg font-bold uppercase tracking-wider text-black">Framework Scores</h3>
            </div>
            <div className="p-6">
              <div className="h-64">
                <Bar
                  data={frameworkData}
                  options={{
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                      y: {
                        beginAtZero: true,
                        max: 100
                      }
                    }
                  }}
                />
              </div>
            </div>
          </div>

          {/* Compliance Status */}
          <div className="border-2 border-black rounded">
            <div className="p-4 bg-gray-50 border-b-2 border-black">
              <h3 className="text-lg font-bold uppercase tracking-wider text-black">Compliance Status</h3>
            </div>
            <div className="p-6">
              <div className="h-64">
                <Doughnut
                  data={statusData}
                  options={{
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                      legend: {
                        position: 'bottom'
                      }
                    }
                  }}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ComplianceDashboard;
