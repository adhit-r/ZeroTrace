import React from 'react';
import { useQuery } from '@tanstack/react-query';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Loader2, CheckCircle, Clock, AlertCircle } from 'lucide-react';

interface ScanJob {
  id: string;
  name: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  progress: number;
  startTime: string;
  endTime?: string;
  assetsScanned: number;
  vulnerabilitiesFound: number;
  errorCount: number;
}

interface ProcessingData {
  activeJobs: ScanJob[];
  completedJobs: ScanJob[];
  totalQueued: number;
  avgProcessingTime: number;
}

const fetchScanProcessing = async (): Promise<ProcessingData> => {
  const response = await fetch('/api/v1/scan/processing');
  if (!response.ok) throw new Error('Failed to fetch scan processing');
  return response.json();
};

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'completed':
      return <CheckCircle className="w-5 h-5 text-green-500" />;
    case 'running':
      return <Loader2 className="w-5 h-5 text-blue-500 animate-spin" />;
    case 'failed':
      return <AlertCircle className="w-5 h-5 text-red-500" />;
    case 'pending':
      return <Clock className="w-5 h-5 text-yellow-500" />;
    default:
      return null;
  }
};

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'completed':
      return 'bg-green-100 text-green-800';
    case 'running':
      return 'bg-blue-100 text-blue-800';
    case 'failed':
      return 'bg-red-100 text-red-800';
    case 'pending':
      return 'bg-yellow-100 text-yellow-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

const ScanProcessing: React.FC = () => {
  const { data, isLoading, error } = useQuery<ProcessingData>({
    queryKey: ['scanProcessing'],
    queryFn: fetchScanProcessing,
    refetchInterval: 10000,
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
              <CardTitle className="text-red-900">Error Loading Scan Processing</CardTitle>
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
          <h1 className="text-3xl font-bold text-gray-900">Scan Processing</h1>
          <p className="text-gray-600 mt-2">Monitor active and completed vulnerability scans</p>
        </div>

        {/* Summary Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Active Scans</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold text-blue-600">{data?.activeJobs.length || 0}</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Queued</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold text-yellow-600">{data?.totalQueued || 0}</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Avg Processing Time</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold text-gray-900">
                {data?.avgProcessingTime ? `${data.avgProcessingTime}s` : 'N/A'}
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Active Scans */}
        <div>
          <h2 className="text-2xl font-semibold text-gray-900 mb-4">Active Scans</h2>
          <div className="space-y-3">
            {data?.activeJobs && data.activeJobs.length > 0 ? (
              data.activeJobs.map((job) => (
                <Card key={job.id} className="hover:shadow-lg transition-shadow">
                  <CardHeader className="pb-3">
                    <div className="flex items-start justify-between">
                      <div className="flex items-center gap-3">
                        {getStatusIcon(job.status)}
                        <div>
                          <CardTitle className="text-lg">{job.name}</CardTitle>
                          <CardDescription className="text-xs">ID: {job.id}</CardDescription>
                        </div>
                      </div>
                      <Badge className={`${getStatusBadgeClass(job.status)} capitalize`}>
                        {job.status}
                      </Badge>
                    </div>
                  </CardHeader>
                  <CardContent className="space-y-3">
                    {/* Progress Bar */}
                    <div>
                      <div className="flex justify-between items-center mb-2">
                        <span className="text-sm text-gray-600">Progress</span>
                        <span className="text-sm font-semibold text-gray-900">{job.progress}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div
                          className="bg-blue-500 h-2 rounded-full transition-all"
                          style={{ width: `${job.progress}%` }}
                        />
                      </div>
                    </div>

                    {/* Stats Grid */}
                    <div className="grid grid-cols-3 gap-4">
                      <div>
                        <p className="text-xs text-gray-600">Assets Scanned</p>
                        <p className="text-lg font-bold text-gray-900">{job.assetsScanned}</p>
                      </div>
                      <div>
                        <p className="text-xs text-gray-600">Vulnerabilities</p>
                        <p className="text-lg font-bold text-red-600">{job.vulnerabilitiesFound}</p>
                      </div>
                      <div>
                        <p className="text-xs text-gray-600">Errors</p>
                        <p className="text-lg font-bold text-orange-600">{job.errorCount}</p>
                      </div>
                    </div>

                    {/* Timeline */}
                    <div className="text-xs text-gray-500">
                      Started: {new Date(job.startTime).toLocaleString()}
                    </div>
                  </CardContent>
                </Card>
              ))
            ) : (
              <Card>
                <CardContent className="pt-6 text-center text-gray-500">
                  No active scans
                </CardContent>
              </Card>
            )}
          </div>
        </div>

        {/* Completed Scans */}
        <div>
          <h2 className="text-2xl font-semibold text-gray-900 mb-4">Recent Completed Scans</h2>
          <div className="space-y-3">
            {data?.completedJobs && data.completedJobs.length > 0 ? (
              data.completedJobs.slice(0, 5).map((job) => (
                <Card key={job.id} className="opacity-75">
                  <CardHeader className="pb-3">
                    <div className="flex items-start justify-between">
                      <div className="flex items-center gap-3">
                        {getStatusIcon(job.status)}
                        <div>
                          <CardTitle className="text-lg">{job.name}</CardTitle>
                          <CardDescription className="text-xs">
                            Completed: {job.endTime ? new Date(job.endTime).toLocaleString() : 'N/A'}
                          </CardDescription>
                        </div>
                      </div>
                      <Badge className={`${getStatusBadgeClass(job.status)} capitalize`}>
                        {job.status}
                      </Badge>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-3 gap-4 text-sm">
                      <div>
                        <p className="text-gray-600">Assets Scanned</p>
                        <p className="font-bold text-gray-900">{job.assetsScanned}</p>
                      </div>
                      <div>
                        <p className="text-gray-600">Vulnerabilities</p>
                        <p className="font-bold text-red-600">{job.vulnerabilitiesFound}</p>
                      </div>
                      <div>
                        <p className="text-gray-600">Errors</p>
                        <p className="font-bold text-orange-600">{job.errorCount}</p>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))
            ) : (
              <Card>
                <CardContent className="pt-6 text-center text-gray-500">
                  No completed scans
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default ScanProcessing;
