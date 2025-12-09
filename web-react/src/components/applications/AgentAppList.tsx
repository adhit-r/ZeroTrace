import React, { useState } from 'react';
import { ChevronDown, ChevronRight, Server, Package, AlertTriangle } from 'lucide-react';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import type { Application } from '@/services/applicationService';

interface AgentAppListProps {
  applications: Application[];
  onAppClick: (app: Application) => void;
}

const AgentAppList: React.FC<AgentAppListProps> = ({ applications, onAppClick }) => {
  const [expandedAgents, setExpandedAgents] = useState<Set<string>>(new Set());

  // Group applications by agent
  const appsByAgent = applications.reduce((acc, app) => {
    if (!acc[app.agentId]) {
      acc[app.agentId] = {
        agentId: app.agentId,
        agentName: app.agentName,
        apps: [],
      };
    }
    acc[app.agentId].apps.push(app);
    return acc;
  }, {} as Record<string, { agentId: string; agentName: string; apps: Application[] }>);

  const toggleAgent = (agentId: string) => {
    const newExpanded = new Set(expandedAgents);
    if (newExpanded.has(agentId)) {
      newExpanded.delete(agentId);
    } else {
      newExpanded.add(agentId);
    }
    setExpandedAgents(newExpanded);
  };

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
    <div className="space-y-4">
      {Object.values(appsByAgent).map((agentGroup) => {
        const isExpanded = expandedAgents.has(agentGroup.agentId);
        const vulnerableCount = agentGroup.apps.filter(a => a.status === 'vulnerable').length;
        const criticalCount = agentGroup.apps.filter(a => a.riskLevel === 'critical').length;

        return (
          <Card
            key={agentGroup.agentId}
            className="bg-white border-4 border-black rounded-lg overflow-hidden"
          >
            {/* Agent Header */}
            <button
              onClick={() => toggleAgent(agentGroup.agentId)}
              className="w-full p-4 flex items-center justify-between hover:bg-gray-50 transition-colors"
            >
              <div className="flex items-center gap-3 flex-1">
                {isExpanded ? (
                  <ChevronDown className="h-5 w-5 text-black" />
                ) : (
                  <ChevronRight className="h-5 w-5 text-black" />
                )}
                <Server className="h-5 w-5 text-blue-600" />
                <div className="flex-1 text-left">
                  <h3 className="font-black text-black uppercase">{agentGroup.agentName}</h3>
                  <p className="text-sm text-gray-600">
                    {agentGroup.apps.length} applications
                    {vulnerableCount > 0 && (
                      <span className="ml-2 text-red-600 font-bold">
                        • {vulnerableCount} vulnerable
                      </span>
                    )}
                    {criticalCount > 0 && (
                      <span className="ml-2 text-red-600 font-bold">
                        • {criticalCount} critical
                      </span>
                    )}
                  </p>
                </div>
              </div>
            </button>

            {/* Applications List */}
            {isExpanded && (
              <div className="border-t-4 border-black p-4 space-y-3">
                {agentGroup.apps.map((app) => (
                  <div
                    key={app.id}
                    onClick={() => onAppClick(app)}
                    className="p-4 bg-gray-50 border-2 border-gray-300 rounded-lg cursor-pointer hover:bg-gray-100 transition-colors"
                  >
                    <div className="flex items-start justify-between mb-2">
                      <div className="flex items-center gap-2 flex-1">
                        <Package className="h-4 w-4 text-orange-600" />
                        <div className="flex-1">
                          <h4 className="font-bold text-black">{app.name}</h4>
                          <p className="text-sm text-gray-600">
                            {app.vendor} • v{app.version}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        {app.vulnerabilities > 0 && (
                          <div className="flex items-center gap-1">
                            <AlertTriangle className="h-4 w-4 text-red-600" />
                            <span className="text-sm font-bold text-red-600">
                              {app.vulnerabilities}
                            </span>
                          </div>
                        )}
                        <Badge className={getRiskBadgeColor(app.riskLevel)}>
                          {app.riskLevel.toUpperCase()}
                        </Badge>
                      </div>
                    </div>
                    {app.classification && (
                      <p className="text-xs text-gray-500 uppercase">
                        {app.classification}
                      </p>
                    )}
                  </div>
                ))}
              </div>
            )}
          </Card>
        );
      })}
    </div>
  );
};

export default AgentAppList;

