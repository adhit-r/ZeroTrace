import React from 'react';
import { useQuery } from '@tanstack/react-query';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Loader2, Download, TrendingUp, TrendingDown } from 'lucide-react';

interface ComplianceFramework {
  name: string;
  status: 'compliant' | 'non-compliant' | 'partial';
  score: number;
  lastAudit: string;
  nextAudit: string;
  controlsPassed: number;
  controlsFailed: number;
  controlsNA: number;
}

interface ComplianceReport {
  generatedDate: string;
  frameworks: ComplianceFramework[];
  overallScore: number;
  overallTrend: 'up' | 'down' | 'stable';
  criticalFindings: number;
  highFindings: number;
  mediumFindings: number;
  lowFindings: number;
}

const fetchComplianceReports = async (): Promise<ComplianceReport> => {
  const response = await fetch('/api/v1/compliance/reports');
  if (!response.ok) throw new Error('Failed to fetch compliance reports');
  return response.json();
};

const getStatusColor = (status: string) => {
  switch (status) {
    case 'compliant':
      return 'bg-green-100 text-green-800';
    case 'non-compliant':
      return 'bg-red-100 text-red-800';
    case 'partial':
      return 'bg-yellow-100 text-yellow-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

const ComplianceReports: React.FC = () => {
  const { data, isLoading, error } = useQuery<ComplianceReport>({
    queryKey: ['complianceReports'],
    queryFn: fetchComplianceReports,
    refetchInterval: 60000,
  });

  if (isLoading) {
    return (
      <Layout>
        <div className="flex items-center justify-center min-h-[400px]">
          <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
        </div>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout>
        <div className="p-6">
          <Card className="border-red-200 bg-red-50">
            <CardHeader>
              <CardTitle className="text-red-900">Error Loading Compliance Reports</CardTitle>
              <CardDescription className="text-red-800">
                {error instanceof Error ? error.message : 'Unknown error occurred'}
              </CardDescription>
            </CardHeader>
          </Card>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-6 p-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Compliance Reports</h1>
            <p className="text-gray-600 mt-2">
              Track your compliance posture across frameworks and standards
            </p>
          </div>
          <button className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
            <Download className="w-4 h-4" />
            Export Report
          </button>
        </div>

        {/* Overall Score */}
        {data && (
          <Card className="border-2 border-blue-200 bg-blue-50">
            <CardHeader>
              <CardTitle className="text-blue-900">Overall Compliance Score</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-8">
                <div>
                  <div className="relative w-32 h-32">
                    <svg className="transform -rotate-90 w-32 h-32">
                      <circle
                        cx="64"
                        cy="64"
                        r="56"
                        fill="none"
                        stroke="#e5e7eb"
                        strokeWidth="8"
                      />
                      <circle
                        cx="64"
                        cy="64"
                        r="56"
                        fill="none"
                        stroke="#3b82f6"
                        strokeWidth="8"
                        strokeDasharray={`${(data.overallScore / 100) * 351.86} 351.86`}
                        className="transition-all"
                      />
                    </svg>
                    <div className="absolute inset-0 flex items-center justify-center">
                      <div className="text-center">
                        <p className="text-3xl font-bold text-blue-900">{data.overallScore}</p>
                        <p className="text-xs text-blue-700">/ 100</p>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="flex-1 space-y-3">
                  <div className="flex items-center gap-2">
                    {data.overallTrend === 'up' ? (
                      <TrendingUp className="w-5 h-5 text-green-500" />
                    ) : data.overallTrend === 'down' ? (
                      <TrendingDown className="w-5 h-5 text-red-500" />
                    ) : (
                      <div className="w-5 h-5 text-gray-500">→</div>
                    )}
                    <span className="text-gray-700">
                      Trend: <span className="font-semibold capitalize">{data.overallTrend}</span>
                    </span>
                  </div>
                  <p className="text-sm text-gray-600">
                    Generated: {new Date(data.generatedDate).toLocaleDateString()}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Findings Summary */}
        {data && (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">Critical</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold text-red-600">{data.criticalFindings}</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">High</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold text-orange-600">{data.highFindings}</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">Medium</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold text-yellow-600">{data.mediumFindings}</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">Low</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-3xl font-bold text-green-600">{data.lowFindings}</p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Frameworks */}
        <div>
          <h2 className="text-2xl font-semibold text-gray-900 mb-4">Compliance Frameworks</h2>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {data?.frameworks && data.frameworks.length > 0 ? (
              data.frameworks.map((framework, idx) => (
                <Card key={idx} className="hover:shadow-lg transition-shadow">
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <div>
                        <CardTitle className="text-lg">{framework.name}</CardTitle>
                        <CardDescription>Compliance Status</CardDescription>
                      </div>
                      <Badge className={`${getStatusColor(framework.status)} capitalize`}>
                        {framework.status}
                      </Badge>
                    </div>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {/* Score Bar */}
                    <div>
                      <div className="flex justify-between items-center mb-2">
                        <span className="text-sm text-gray-600">Compliance Score</span>
                        <span className="text-sm font-bold text-gray-900">{framework.score}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div
                          className={`h-2 rounded-full transition-all ${
                            framework.score >= 80
                              ? 'bg-green-500'
                              : framework.score >= 60
                                ? 'bg-yellow-500'
                                : 'bg-red-500'
                          }`}
                          style={{ width: `${framework.score}%` }}
                        />
                      </div>
                    </div>

                    {/* Controls */}
                    <div className="grid grid-cols-3 gap-3 text-sm">
                      <div className="bg-green-50 p-2 rounded">
                        <p className="text-green-700 font-semibold text-lg">{framework.controlsPassed}</p>
                        <p className="text-green-600 text-xs">Passed</p>
                      </div>
                      <div className="bg-red-50 p-2 rounded">
                        <p className="text-red-700 font-semibold text-lg">{framework.controlsFailed}</p>
                        <p className="text-red-600 text-xs">Failed</p>
                      </div>
                      <div className="bg-gray-50 p-2 rounded">
                        <p className="text-gray-700 font-semibold text-lg">{framework.controlsNA}</p>
                        <p className="text-gray-600 text-xs">N/A</p>
                      </div>
                    </div>

                    {/* Dates */}
                    <div className="text-xs text-gray-600 space-y-1 pt-2 border-t">
                      <p>
                        Last Audit: <span className="font-semibold">{new Date(framework.lastAudit).toLocaleDateString()}</span>
                      </p>
                      <p>
                        Next Audit: <span className="font-semibold">{new Date(framework.nextAudit).toLocaleDateString()}</span>
                      </p>
                    </div>
                  </CardContent>
                </Card>
              ))
            ) : (
              <Card className="col-span-2">
                <CardContent className="pt-6 text-center text-gray-500">
                  No compliance frameworks configured
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default ComplianceReports;
