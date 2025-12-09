import React, { useState, useEffect, useMemo } from 'react';
import { Package, AlertTriangle } from 'lucide-react';
import { Card } from '@/components/ui/card';
import { applicationService } from '@/services/applicationService';
import type { Application, ApplicationStats } from '@/services/applicationService';
import ApplicationCharts from '@/components/applications/ApplicationCharts';
import ApplicationCards from '@/components/applications/ApplicationCards';
import ApplicationFilters from '@/components/applications/ApplicationFilters';
import type { GroupingType, FilterType } from '@/components/applications/ApplicationFilters';
import AgentAppList from '@/components/applications/AgentAppList';
import AppDetailModal from '@/components/applications/AppDetailModal';

const ApplicationAnalysis: React.FC = () => {
  const [applications, setApplications] = useState<Application[]>([]);
  const [stats, setStats] = useState<ApplicationStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [grouping, setGrouping] = useState<GroupingType>('agent');
  const [filterType, setFilterType] = useState<FilterType>('all');
  const [selectedApp, setSelectedApp] = useState<Application | null>(null);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setIsLoading(true);
      const [apps, appStats] = await Promise.all([
        applicationService.getAllApplications(),
        applicationService.getApplicationStats(),
      ]);
      setApplications(apps);
      setStats(appStats);
    } catch (error) {
      console.error('Failed to fetch application data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // Filter and group applications
  const filteredAndGroupedApps = useMemo(() => {
    let filtered = [...applications];

    // Apply search filter
    if (searchTerm) {
      filtered = filtered.filter(
        app =>
          app.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          app.vendor.toLowerCase().includes(searchTerm.toLowerCase()) ||
          app.version.toLowerCase().includes(searchTerm.toLowerCase()) ||
          app.agentName.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Apply status filter
    if (filterType !== 'all') {
      filtered = filtered.filter(app => app.status === filterType);
    }

    // Group applications
    if (grouping === 'risk') {
      const grouped: Record<string, Application[]> = {
        critical: [],
        high: [],
        medium: [],
        low: [],
        safe: [],
      };
      filtered.forEach(app => {
        grouped[app.riskLevel].push(app);
      });
      return grouped;
    } else if (grouping === 'classification') {
      const grouped: Record<string, Application[]> = {};
      filtered.forEach(app => {
        const classification = app.classification || 'unknown';
        if (!grouped[classification]) {
          grouped[classification] = [];
        }
        grouped[classification].push(app);
      });
      return grouped;
    } else {
      // Group by agent (default)
      const grouped: Record<string, Application[]> = {};
      filtered.forEach(app => {
        if (!grouped[app.agentId]) {
          grouped[app.agentId] = [];
        }
        grouped[app.agentId].push(app);
      });
      return grouped;
    }
  }, [applications, searchTerm, filterType, grouping]);

  // Get top vulnerable applications
  const topVulnerableApps = useMemo(() => {
    return [...applications]
      .filter(app => app.status === 'vulnerable')
      .sort((a, b) => {
        // Sort by risk level first, then by vulnerability count
        const riskOrder = { critical: 5, high: 4, medium: 3, low: 2, safe: 1 };
        const riskDiff = riskOrder[b.riskLevel] - riskOrder[a.riskLevel];
        if (riskDiff !== 0) return riskDiff;
        return b.vulnerabilities - a.vulnerabilities;
      });
  }, [applications]);

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading application analysis...</p>
          </div>
        </div>
      </div>
    );
  }

  if (!stats) {
    return (
      <div className="p-6">
        <Card className="p-12 text-center bg-white border-4 border-black rounded-lg">
          <AlertTriangle className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-xl font-bold text-gray-900 mb-2">Failed to Load Data</h3>
          <p className="text-gray-600">Unable to fetch application data. Please try again later.</p>
        </Card>
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-black text-black uppercase mb-2">Application Analysis</h1>
        <p className="text-gray-600">
          Comprehensive analysis of installed applications across all agents with risk assessment and vulnerability tracking
        </p>
      </div>

      {/* Summary Cards */}
      <ApplicationCards stats={stats} topVulnerableApps={topVulnerableApps} />

      {/* Charts */}
      <ApplicationCharts applications={applications} stats={stats} />

      {/* Filters */}
      <ApplicationFilters
        searchTerm={searchTerm}
        onSearchChange={setSearchTerm}
        grouping={grouping}
        onGroupingChange={setGrouping}
        filterType={filterType}
        onFilterChange={setFilterType}
      />

      {/* Applications List - Grouped by Agent */}
      {grouping === 'agent' && (
        <div>
          <h2 className="text-2xl font-black text-black uppercase mb-4">Applications by Agent</h2>
          <AgentAppList
            applications={Object.values(filteredAndGroupedApps).flat()}
            onAppClick={setSelectedApp}
          />
        </div>
      )}

      {/* Applications List - Grouped by Risk */}
      {grouping === 'risk' && (
        <div className="space-y-6">
          <h2 className="text-2xl font-black text-black uppercase">Applications by Risk Level</h2>
          {(['critical', 'high', 'medium', 'low', 'safe'] as const).map(riskLevel => {
            const apps = filteredAndGroupedApps[riskLevel] || [];
            if (apps.length === 0) return null;

            return (
              <Card key={riskLevel} className="p-6 bg-white border-4 border-black rounded-lg">
                <h3 className="text-xl font-black text-black uppercase mb-4">
                  {riskLevel.toUpperCase()} Risk ({apps.length})
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {apps.map(app => (
                    <div
                      key={app.id}
                      onClick={() => setSelectedApp(app)}
                      className="p-4 bg-gray-50 border-2 border-gray-300 rounded-lg cursor-pointer hover:bg-gray-100 transition-colors"
                    >
                      <div className="flex items-start justify-between mb-2">
                        <div className="flex-1">
                          <h4 className="font-bold text-black">{app.name}</h4>
                          <p className="text-sm text-gray-600">
                            {app.vendor} • v{app.version}
                          </p>
                        </div>
                        {app.vulnerabilities > 0 && (
                          <div className="flex items-center gap-1">
                            <AlertTriangle className="h-4 w-4 text-red-600" />
                            <span className="text-sm font-bold text-red-600">
                              {app.vulnerabilities}
                            </span>
                          </div>
                        )}
                      </div>
                      <p className="text-xs text-gray-500">{app.agentName}</p>
                    </div>
                  ))}
                </div>
              </Card>
            );
          })}
        </div>
      )}

      {/* Applications List - Grouped by Classification */}
      {grouping === 'classification' && (
        <div className="space-y-6">
          <h2 className="text-2xl font-black text-black uppercase">Applications by Classification</h2>
          {Object.entries(filteredAndGroupedApps).map(([classification, apps]) => (
            <Card key={classification} className="p-6 bg-white border-4 border-black rounded-lg">
              <h3 className="text-xl font-black text-black uppercase mb-4">
                {classification.toUpperCase()} ({apps.length})
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {apps.map(app => (
                  <div
                    key={app.id}
                    onClick={() => setSelectedApp(app)}
                    className="p-4 bg-gray-50 border-2 border-gray-300 rounded-lg cursor-pointer hover:bg-gray-100 transition-colors"
                  >
                    <div className="flex items-start justify-between mb-2">
                      <div className="flex-1">
                        <h4 className="font-bold text-black">{app.name}</h4>
                        <p className="text-sm text-gray-600">
                          {app.vendor} • v{app.version}
                        </p>
                      </div>
                      {app.vulnerabilities > 0 && (
                        <div className="flex items-center gap-1">
                          <AlertTriangle className="h-4 w-4 text-red-600" />
                          <span className="text-sm font-bold text-red-600">
                            {app.vulnerabilities}
                          </span>
                        </div>
                      )}
                    </div>
                    <p className="text-xs text-gray-500">{app.agentName}</p>
                  </div>
                ))}
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* Empty State */}
      {Object.values(filteredAndGroupedApps).flat().length === 0 && (
        <Card className="p-12 text-center bg-white border-4 border-black rounded-lg">
          <Package className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-xl font-bold text-gray-900 mb-2">No Applications Found</h3>
          <p className="text-gray-600">
            {searchTerm || filterType !== 'all'
              ? 'Try adjusting your search or filters'
              : 'No applications detected yet. Agents will report installed applications after scanning.'}
          </p>
        </Card>
      )}

      {/* App Detail Modal */}
      <AppDetailModal application={selectedApp} onClose={() => setSelectedApp(null)} />
    </div>
  );
};

export default ApplicationAnalysis;

