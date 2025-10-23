import React from 'react';
import { useQuery } from '@tanstack/react-query';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Loader2, Activity, AlertCircle, CheckCircle } from 'lucide-react';

interface ScannerModule {
  id: string;
  name: string;
  category: string;
  status: 'active' | 'inactive' | 'error';
  description: string;
  itemsScanned: number;
  vulnerabilitiesFound: number;
  lastRun: string;
  nextRun: string;
  successRate: number;
  errorCount: number;
  avgScanTime: number;
}

interface ScannerData {
  totalScanners: number;
  activeScanners: number;
  failedScanners: number;
  modules: ScannerModule[];
  overallHealth: number;
}

const fetchScannerDetails = async (): Promise<ScannerData> => {
  const response = await fetch('/api/v1/scanners/details');
  if (!response.ok) throw new Error('Failed to fetch scanner details');
  return response.json();
};

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'active':
      return <Activity className="w-5 h-5 text-green-500" />;
    case 'error':
      return <AlertCircle className="w-5 h-5 text-red-500" />;
    case 'inactive':
      return <CheckCircle className="w-5 h-5 text-gray-500" />;
    default:
      return null;
  }
};

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'active':
      return 'bg-green-100 text-green-800';
    case 'inactive':
      return 'bg-gray-100 text-gray-800';
    case 'error':
      return 'bg-red-100 text-red-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

const ScannerDetails: React.FC = () => {
  const { data, isLoading, error } = useQuery<ScannerData>({
    queryKey: ['scannerDetails'],
    queryFn: fetchScannerDetails,
    refetchInterval: 30000,
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
              <CardTitle className="text-red-900">Error Loading Scanner Details</CardTitle>
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
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Scanner Modules</h1>
          <p className="text-gray-600 mt-2">Monitor individual scanner performance and health</p>
        </div>

        {/* Summary Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Total Scanners</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold text-gray-900">{data?.totalScanners || 0}</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Active</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold text-green-600">{data?.activeScanners || 0}</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Failed</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold text-red-600">{data?.failedScanners || 0}</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Overall Health</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold text-blue-600">{data?.overallHealth || 0}%</p>
            </CardContent>
          </Card>
        </div>

        {/* Scanner Modules Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {data?.modules && data.modules.length > 0 ? (
            data.modules.map((module) => (
              <Card key={module.id} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-3">
                      {getStatusIcon(module.status)}
                      <div>
                        <CardTitle className="text-lg">{module.name}</CardTitle>
                        <CardDescription className="text-xs">{module.category}</CardDescription>
                      </div>
                    </div>
                    <Badge className={`${getStatusBadgeClass(module.status)} capitalize`}>
                      {module.status}
                    </Badge>
                  </div>
                </CardHeader>
                <CardContent className="space-y-4">
                  {/* Description */}
                  <p className="text-sm text-gray-600">{module.description}</p>

                  {/* Success Rate */}
                  <div>
                    <div className="flex justify-between items-center mb-2">
                      <span className="text-sm text-gray-600">Success Rate</span>
                      <span className="text-sm font-bold text-gray-900">{module.successRate}%</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div
                        className={`h-2 rounded-full transition-all ${
                          module.successRate >= 95
                            ? 'bg-green-500'
                            : module.successRate >= 80
                              ? 'bg-yellow-500'
                              : 'bg-red-500'
                        }`}
                        style={{ width: `${module.successRate}%` }}
                      />
                    </div>
                  </div>

                  {/* Stats Grid */}
                  <div className="grid grid-cols-3 gap-3 text-sm">
                    <div className="bg-blue-50 p-2 rounded">
                      <p className="text-blue-700 font-semibold text-lg">{module.itemsScanned}</p>
                      <p className="text-blue-600 text-xs">Items Scanned</p>
                    </div>
                    <div className="bg-red-50 p-2 rounded">
                      <p className="text-red-700 font-semibold text-lg">{module.vulnerabilitiesFound}</p>
                      <p className="text-red-600 text-xs">Vulnerabilities</p>
                    </div>
                    <div className="bg-orange-50 p-2 rounded">
                      <p className="text-orange-700 font-semibold text-lg">{module.errorCount}</p>
                      <p className="text-orange-600 text-xs">Errors</p>
                    </div>
                  </div>

                  {/* Timing Info */}
                  <div className="grid grid-cols-2 gap-2 text-xs text-gray-600 pt-2 border-t">
                    <div>
                      <p className="text-gray-500">Last Run</p>
                      <p className="font-semibold text-gray-900">
                        {new Date(module.lastRun).toLocaleTimeString()}
                      </p>
                    </div>
                    <div>
                      <p className="text-gray-500">Avg Scan Time</p>
                      <p className="font-semibold text-gray-900">{module.avgScanTime}s</p>
                    </div>
                  </div>

                  {/* Next Run */}
                  <div className="text-xs text-gray-600">
                    <p className="text-gray-500">Next Run</p>
                    <p className="font-semibold text-gray-900">
                      {new Date(module.nextRun).toLocaleDateString()} at{' '}
                      {new Date(module.nextRun).toLocaleTimeString()}
                    </p>
                  </div>
                </CardContent>
              </Card>
            ))
          ) : (
            <Card className="col-span-2">
              <CardContent className="pt-6 text-center text-gray-500">
                No scanner modules available
              </CardContent>
            </Card>
          )}
        </div>

        {/* Scanner Categories Summary */}
        {data?.modules && data.modules.length > 0 && (
          <Card>
            <CardHeader>
              <CardTitle>Scanners by Category</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                {Array.from(new Set(data.modules.map((m) => m.category))).map((category) => {
                  const count = data.modules.filter((m) => m.category === category).length;
                  const active = data.modules.filter(
                    (m) => m.category === category && m.status === 'active'
                  ).length;
                  return (
                    <div key={category} className="text-center p-4 bg-gray-50 rounded-lg">
                      <p className="font-semibold text-gray-900">{count}</p>
                      <p className="text-xs text-gray-600 mt-1">{category}</p>
                      <p className="text-xs text-green-600 mt-1">{active} active</p>
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </Layout>
  );
};

export default ScannerDetails;
