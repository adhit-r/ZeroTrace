import React from 'react';
import { Package, AlertTriangle, CheckCircle, Shield } from 'lucide-react';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import type { Application, ApplicationStats } from '@/services/applicationService';

interface ApplicationCardsProps {
  stats: ApplicationStats;
  topVulnerableApps: Application[];
}

const ApplicationCards: React.FC<ApplicationCardsProps> = ({ stats, topVulnerableApps }) => {
  const getRiskBadgeColor = (risk: string) => {
    switch (risk) {
      case 'critical':
        return 'bg-red-100 text-red-800 border-red-300';
      case 'high':
        return 'bg-orange-100 text-orange-800 border-orange-300';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800 border-yellow-300';
      case 'low':
        return 'bg-blue-100 text-blue-800 border-blue-300';
      default:
        return 'bg-green-100 text-green-800 border-green-300';
    }
  };

  return (
    <div className="space-y-6">
      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card className="p-6 bg-white border-4 border-black rounded-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 uppercase font-bold mb-1">Total Apps</p>
              <p className="text-3xl font-black text-black">{stats.total}</p>
            </div>
            <Package className="h-8 w-8 text-blue-600" />
          </div>
        </Card>

        <Card className="p-6 bg-white border-4 border-black rounded-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 uppercase font-bold mb-1">Vulnerable</p>
              <p className="text-3xl font-black text-red-600">{stats.vulnerable}</p>
            </div>
            <AlertTriangle className="h-8 w-8 text-red-600" />
          </div>
        </Card>

        <Card className="p-6 bg-white border-4 border-black rounded-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 uppercase font-bold mb-1">Critical Risk</p>
              <p className="text-3xl font-black text-red-600">{stats.byRiskLevel.critical}</p>
            </div>
            <Shield className="h-8 w-8 text-red-600" />
          </div>
        </Card>

        <Card className="p-6 bg-white border-4 border-black rounded-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 uppercase font-bold mb-1">Safe</p>
              <p className="text-3xl font-black text-green-600">{stats.safe}</p>
            </div>
            <CheckCircle className="h-8 w-8 text-green-600" />
          </div>
        </Card>
      </div>

      {/* Risk Level Breakdown */}
      <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
        <Card className="p-4 bg-white border-4 border-black rounded-lg text-center">
          <p className="text-xs text-gray-600 uppercase font-bold mb-1">Critical</p>
          <p className="text-2xl font-black text-red-600">{stats.byRiskLevel.critical}</p>
        </Card>
        <Card className="p-4 bg-white border-4 border-black rounded-lg text-center">
          <p className="text-xs text-gray-600 uppercase font-bold mb-1">High</p>
          <p className="text-2xl font-black text-orange-600">{stats.byRiskLevel.high}</p>
        </Card>
        <Card className="p-4 bg-white border-4 border-black rounded-lg text-center">
          <p className="text-xs text-gray-600 uppercase font-bold mb-1">Medium</p>
          <p className="text-2xl font-black text-yellow-600">{stats.byRiskLevel.medium}</p>
        </Card>
        <Card className="p-4 bg-white border-4 border-black rounded-lg text-center">
          <p className="text-xs text-gray-600 uppercase font-bold mb-1">Low</p>
          <p className="text-2xl font-black text-blue-600">{stats.byRiskLevel.low}</p>
        </Card>
        <Card className="p-4 bg-white border-4 border-black rounded-lg text-center">
          <p className="text-xs text-gray-600 uppercase font-bold mb-1">Safe</p>
          <p className="text-2xl font-black text-green-600">{stats.byRiskLevel.safe}</p>
        </Card>
      </div>

      {/* Top Vulnerable Applications */}
      {topVulnerableApps.length > 0 && (
        <Card className="p-6 bg-white border-4 border-black rounded-lg">
          <h3 className="text-xl font-black text-black uppercase mb-4">Top Vulnerable Applications</h3>
          <div className="space-y-3">
            {topVulnerableApps.slice(0, 5).map((app) => (
              <div
                key={app.id}
                className="flex items-center justify-between p-4 bg-gray-50 border-2 border-gray-300 rounded-lg"
              >
                <div className="flex items-center gap-3 flex-1">
                  <Package className="h-5 w-5 text-orange-600" />
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <h4 className="font-bold text-black">{app.name}</h4>
                      <Badge className={getRiskBadgeColor(app.riskLevel)}>
                        {app.riskLevel.toUpperCase()}
                      </Badge>
                    </div>
                    <p className="text-sm text-gray-600">
                      {app.vendor} • v{app.version} • {app.agentName}
                    </p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-lg font-black text-red-600">{app.vulnerabilities}</p>
                  <p className="text-xs text-gray-600">Vulnerabilities</p>
                </div>
              </div>
            ))}
          </div>
        </Card>
      )}
    </div>
  );
};

export default ApplicationCards;

