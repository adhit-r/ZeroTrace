import React from 'react';
import { useQuery } from '@tanstack/react-query';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Loader2 } from 'lucide-react';

interface Technology {
  name: string;
  category: string;
  instances: number;
  version?: string;
  riskLevel: 'critical' | 'high' | 'medium' | 'low';
  discoveredAssets: string[];
}

interface TechStackData {
  timestamp: string;
  totalTechnologies: number;
  byCategory: Record<string, number>;
  technologies: Technology[];
}

const fetchTechStack = async (): Promise<TechStackData> => {
  const response = await fetch('/api/v1/tech-stack');
  if (!response.ok) throw new Error('Failed to fetch tech stack');
  return response.json();
};

const getRiskColor = (risk: string) => {
  switch (risk) {
    case 'critical':
      return 'bg-red-100 text-red-800';
    case 'high':
      return 'bg-orange-100 text-orange-800';
    case 'medium':
      return 'bg-yellow-100 text-yellow-800';
    case 'low':
      return 'bg-green-100 text-green-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

const TechStack: React.FC = () => {
  const { data, isLoading, error } = useQuery<TechStackData>({
    queryKey: ['techStack'],
    queryFn: fetchTechStack,
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
              <CardTitle className="text-red-900">Error Loading Tech Stack</CardTitle>
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
          <h1 className="text-3xl font-bold text-gray-900">Technology Stack Discovery</h1>
          <p className="text-gray-600 mt-2">
            Complete inventory of software, frameworks, and technologies detected across all assets
          </p>
        </div>

        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Total Technologies</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold text-gray-900">{data?.totalTechnologies || 0}</p>
            </CardContent>
          </Card>

          {data?.byCategory &&
            Object.entries(data.byCategory).slice(0, 3).map(([category, count]) => (
              <Card key={category}>
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm font-medium text-gray-600">{category}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-3xl font-bold text-gray-900">{count}</p>
                </CardContent>
              </Card>
            ))}
        </div>

        {/* Technology List by Category */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {data?.technologies && data.technologies.length > 0 ? (
            data.technologies.map((tech, idx) => (
              <Card key={idx} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div>
                      <CardTitle className="text-lg">{tech.name}</CardTitle>
                      <CardDescription>{tech.category}</CardDescription>
                    </div>
                    <Badge className={`${getRiskColor(tech.riskLevel)} capitalize`}>
                      {tech.riskLevel} Risk
                    </Badge>
                  </div>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <p className="text-sm text-gray-600">Instances</p>
                      <p className="text-2xl font-bold text-gray-900">{tech.instances}</p>
                    </div>
                    {tech.version && (
                      <div>
                        <p className="text-sm text-gray-600">Version</p>
                        <p className="text-lg font-semibold text-gray-900">{tech.version}</p>
                      </div>
                    )}
                  </div>
                  <div>
                    <p className="text-sm text-gray-600 mb-2">Found in Assets</p>
                    <div className="flex flex-wrap gap-2">
                      {tech.discoveredAssets.map((asset, idx) => (
                        <Badge key={idx} variant="outline" className="text-xs">
                          {asset}
                        </Badge>
                      ))}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))
          ) : (
            <Card className="col-span-2">
              <CardContent className="pt-6 text-center text-gray-500">
                No technologies discovered yet
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default TechStack;
